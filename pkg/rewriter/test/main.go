package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/goriiin/kotyari-bots_backend/pkg/config"
	"github.com/goriiin/kotyari-bots_backend/pkg/grok"
	proxyPkg "github.com/goriiin/kotyari-bots_backend/pkg/proxy"
	"github.com/goriiin/kotyari-bots_backend/pkg/rewriter"
)

func main() {
	cfg := rewriter.Config{
		NumRewrites: 5,
		Timeout:     10 * time.Second,
	}

	ctx := context.Background()

	llmConfig, err := config.New[grok.GrokClientConfig]()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", llmConfig.ApiKey)

	proxyConfig, err := config.New[proxyPkg.ProxyConfig]()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", proxyConfig.ProxyAPI)

	llm, err := grok.NewGrokClient(llmConfig, proxyConfig)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", llm)

	rw := rewriter.NewGrokRewriter(cfg, llm, "grok-3-mini")

	// Тестовые входные данные
	user := "мейнкун"
	profile := "Пиши информативно и дружелюбно для начинающих владельцев."
	bot := "Структурируй подзаголовками, добавляй практические советы."

	// 1) Переписывание
	variants, err := rw.Rewrite(ctx, user, profile, bot)
	if err != nil {
		log.Fatalf("rewrite error: %v", err)
	}
	if len(variants) == 0 {
		log.Fatalf("rewrite produced no variants")
	}
	fmt.Println("Переписанные варианты:")
	for i, v := range variants {
		fmt.Printf("  %d) %s\n", i+1, v)
	}
}
