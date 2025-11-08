// path: pkg/evals/judge.go
package evals

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type grokClient interface {
	Generate(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}

type Config struct {
	Timeout time.Duration
	MaxLen  int    // опционально: ограничение длины текста кандидата перед судейством
	Model   string // если потребуется маршрутизация моделей
}

type Judge struct {
	cfg    Config
	client grokClient
}

func NewJudge(cfg Config, client grokClient) *Judge {
	if cfg.Timeout == 0 {
		cfg.Timeout = 15 * time.Second
	}
	return &Judge{cfg: cfg, client: client}
}

type Candidate struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
	// Можно добавить любые тех. поля (id, score из retriever и т.д.)
}

type judgeInput struct {
	UserPrompt    string      `json:"user_prompt"`
	ProfilePrompt string      `json:"profile_prompt,omitempty"`
	BotPrompt     string      `json:"bot_prompt,omitempty"`
	Candidates    []Candidate `json:"candidates"`
}

type judgeOutput struct {
	BestIndex int               `json:"best_index"`
	Reasons   map[string]string `json:"reasons,omitempty"` // "k0": "...", "k1": "..."
}

// SelectBest выбирает лучший ответ из candidates.
// Возвращает индекс лучшего и нормализованного кандидата.
func (j *Judge) SelectBest(ctx context.Context, userPrompt, profilePrompt, botPrompt string, candidates []Candidate) (int, Candidate, error) {
	if len(candidates) == 0 {
		return -1, Candidate{}, errors.New("no candidates")
	}
	if len(candidates) == 1 {
		return 0, candidates[0], nil
	}

	// Подрежем тексты, если задан MaxLen.
	if j.cfg.MaxLen > 0 {
		for i := range candidates {
			candidates[i].Text = trimRunes(candidates[i].Text, j.cfg.MaxLen)
			candidates[i].Title = trimRunes(candidates[i].Title, j.cfg.MaxLen/5)
		}
	}

	in := judgeInput{
		UserPrompt:    userPrompt,
		ProfilePrompt: profilePrompt,
		BotPrompt:     botPrompt,
		Candidates:    candidates,
	}
	usr, _ := jsoniter.Marshal(in)

	ctx, cancel := context.WithTimeout(ctx, j.cfg.Timeout)
	defer cancel()

	sys := judgeSystemPromptRU
	raw, err := j.client.Generate(ctx, sys, string(usr))
	if err != nil {
		return -1, Candidate{}, fmt.Errorf("judge generate failed: %w", err)
	}

	out, err := parseJudgeJSON(raw)
	if err != nil {
		return -1, Candidate{}, fmt.Errorf("judge parse failed: %w; raw=%s", err, raw)
	}
	if out.BestIndex < 0 || out.BestIndex >= len(candidates) {
		return -1, Candidate{}, fmt.Errorf("judge returned invalid index: %d", out.BestIndex)
	}
	return out.BestIndex, candidates[out.BestIndex], nil
}

const judgeSystemPromptRU = `
Ты — строгий судья качества генераций из RAG-конвейера.
Требуется выбрать лучший вариант ответа для публикации.

Критерии (в порядке убывания важности):
1) Релевантность и полнота ответа относительно задачи пользователя.
2) Фактуальность и отсутствие выдумок.
3) Соблюдение ограничений и стиля (бот/профиль/платформа).
4) Ясность, структура (заголовок/текст), читабельность.

Формат ввода: JSON с полями user_prompt, profile_prompt, bot_prompt, candidates[].
Формат вывода: строго JSON: {"best_index": <int>, "reasons": {"k0":"...","k1":"..."}}.
Никакого текста вне JSON.
`

func parseJudgeJSON(s string) (judgeOutput, error) {
	// Часто модели могут добавить текст вокруг; попытаемся выделить первый JSON-объект.
	start := strings.IndexByte(s, '{')
	end := strings.LastIndexByte(s, '}')
	var out judgeOutput
	if start >= 0 && end > start {
		if err := jsoniter.Unmarshal([]byte(s[start:end+1]), &out); err == nil {
			return out, nil
		}
	}
	// Прямая попытка
	if err := jsoniter.Unmarshal([]byte(s), &out); err != nil {
		return out, err
	}
	return out, nil
}

func trimRunes(s string, max int) string {
	if max <= 0 {
		return s
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max]) + "…"
}
