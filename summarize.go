package main

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"github.com/uptrace/bun"
)

const openAIKey = OPENAI_API_KEY // Replace with your actual API key

func summarizeArticles(ctx context.Context, db *bun.DB) error {
	articles := make([]Article, 0)
	err := db.NewSelect().Model(&articles).Where("summarized = false").Limit(5).OrderExpr("published_at DESC").Scan(ctx)
	if err != nil {
		return err
	}

	client := openai.NewClient(openAIKey)

	for _, article := range articles {
		req := openai.CompletionRequest{
			Model:     openai.GPT3Ada, // Choose the appropriate model
			MaxTokens: 150,
			Prompt:    fmt.Sprintf("Please summarize this article: %s", article.Content),
		}
		resp, err := client.CreateCompletion(ctx, req)
		if err != nil {
			return err
		}

		summary := resp.Choices[0].Text
		article.Content = summary
		article.Summarized = true
		_, err = db.NewUpdate().Model(&article).WherePK().Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
