package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"github.com/goriiin/kotyari-bots_backend/pkg/evals"
	"github.com/goriiin/kotyari-bots_backend/pkg/grok"
	proxyPkg "github.com/goriiin/kotyari-bots_backend/pkg/proxy"
)

func main() {
	ctx := context.Background()

	// Конфигурация для Judge
	cfg := evals.Config{
		Timeout: 15 * time.Second,
		Model:   "grok-2-mini",
	}

	// Загрузка конфигурации для Grok клиента и прокси
	llmConfig, err := config.New[grok.GrokClientConfig]()
	if err != nil {
		log.Fatalf("failed to load llm config: %v", err)
	}

	proxyConfig, err := config.New[proxyPkg.ProxyConfig]()
	if err != nil {
		log.Fatalf("failed to load proxy config: %v", err)
	}

	// Создание клиента GroK с использованием прокси
	llm, err := grok.NewGrokClient(llmConfig, proxyConfig)
	if err != nil {
		log.Fatalf("failed to create grok client: %v", err)
	}

	// Создание нового Judge
	judge := evals.NewJudge(cfg, llm)

	// Промпты для оценки
	userPrompt := "Расскажи о породе кошек мейн-кун."
	profilePrompt := ""
	botPrompt := "Ты — дружелюбный ассистент, который любит кошек."

	// Кандидаты для оценки
	candidates := []model.Candidate{
		{
			Title: "Мейн-кун: нежный гигант",
			Text:  "Мейн-куны — одна из самых крупных пород домашних кошек. Они известны своим дружелюбным и игривым характером, что делает их прекрасными компаньонами. Несмотря на свой размер, они очень нежны и хорошо ладят с детьми и другими животными.",
		},
		{
			Title: "Все о мейн-кунах",
			Text:  "Мейн-кун — порода кошек, которая произошла из Северной Америки. Их отличительные черты — большой размер, пушистый хвост и кисточки на ушах, как у рыси. Они требуют регулярного ухода за шерстью, чтобы избежать колтунов.",
		},
		{
			Title: "Мейн-куны",
			Text:  "Большие кошки с кисточками на ушах. Дружелюбные.",
		},
	}

	// Выбор лучшего кандидата
	bestCandidate, err := judge.SelectBest(ctx, userPrompt, profilePrompt, botPrompt, candidates)
	if err != nil {
		log.Fatalf("failed to select best candidate: %v", err)
	}

	fmt.Printf("Заголовок: %s\n", bestCandidate.Title)
	fmt.Printf("Текст: %s\n", bestCandidate.Text)
}
