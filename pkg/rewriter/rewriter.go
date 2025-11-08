package rewriter

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	posts "github.com/goriiin/kotyari-bots_backend/api/protos/posts/gen"
	jsoniter "github.com/json-iterator/go"
)

type Config struct {
	NumRewrites int           // K: 3–5 по умолчанию
	Timeout     time.Duration // общий таймаут на переписывание
	MaxLen      int           // необяз: обрезка слишком длинных вариантов
}

type grokClient interface {
	Generate(ctx context.Context, systemPrompt, userPrompt string) (string, error)
	GeneratePost(ctx context.Context, botPrompt, taskText, profilePrompt string) (string, error)
}

// Реализация на Grok
type Rewriter struct {
	cfg    Config
	client grokClient
	model  string // напр. "grok-2" / "grok-2-mini"
}

func NewGrokRewriter(cfg Config, client grokClient, model string) *Rewriter {
	// Значения по умолчанию
	if cfg.NumRewrites <= 0 {
		cfg.NumRewrites = 3
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 15 * time.Second
	}
	return &Rewriter{cfg: cfg, client: client, model: model}
}

func (r *Rewriter) Rewrite(ctx context.Context, user, profile, bot string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.Timeout)
	defer cancel()

	system := RewriterSystemPrompt
	userMsg := buildUserPrompt(user, profile, bot, r.cfg.NumRewrites)

	raw, err := r.client.Generate(ctx, system, userMsg)
	if err != nil {
		return nil, err
	}
	rewrites, err := parseJSONList(raw)
	if err != nil {
		return []string{user}, nil
	}

	rewrites = dedupAndClean(rewrites, r.cfg.MaxLen)
	if len(rewrites) == 0 {
		rewrites = []string{user}
	}
	return rewrites, nil
}

func buildUserPrompt(user, profile, bot string, k int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "K: %d\n", k)
	if profile != "" {
		fmt.Fprintf(&b, "profile_prompt: %q\n", profile)
	}
	if bot != "" {
		fmt.Fprintf(&b, "bot_prompt: %q\n", bot)
	}
	fmt.Fprintf(&b, "input: %q\n", user)
	return b.String()
}

func parseJSONList(s string) ([]string, error) {
	var arr []string
	dec := jsoniter.NewDecoder(strings.NewReader(s))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&arr); err != nil {
		// часто модели добавляют префиксы; попробуем найти первую '['
		if idx := strings.Index(s, "["); idx >= 0 {
			if j := strings.LastIndex(s, "]"); j > idx {
				sub := s[idx : j+1]
				if err2 := jsoniter.Unmarshal([]byte(sub), &arr); err2 == nil {
					return arr, nil
				}
			}
		}
		return nil, errors.New("model output is not a JSON string array")
	}
	return arr, nil
}

func dedupAndClean(in []string, maxLen int) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if maxLen > 0 && len([]rune(s)) > maxLen {
			words := strings.Fields(s)
			for len([]rune(strings.Join(words, " "))) > maxLen && len(words) > 1 {
				words = words[:len(words)-1]
			}
			s = strings.Join(words, " ")
		}
		key := strings.ToLower(s)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

func BuildBatchRequests(user, profile, bot string, rewrites []string) *posts.GetPostsRequest {
	req := &posts.GetPostsRequest{}
	for _, r := range rewrites {
		req.PostsRequest = append(req.PostsRequest, &posts.GetPostRequest{
			UserPrompt:    r,
			ProfilePrompt: profile,
			BotPrompt:     bot,
		})
	}
	return req
}

// Заглушка выбора лучшего варианта — потом замените на LLM-as-a-Judge
func SelectBest(candidates []*posts.GetPostResponse) *posts.GetPostResponse {
	if len(candidates) == 0 {
		return nil
	}
	return candidates[0]
}
