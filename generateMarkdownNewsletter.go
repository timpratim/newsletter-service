package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/uptrace/bun"
)

func generateMarkdownNewsletter(ctx context.Context, db *bun.DB) (string, error) {
	// Fetch summarized articles from the database
	var articles []Article
	err := db.NewSelect().Model(&articles).Where("summarized = true").OrderExpr("published_at DESC").Limit(5).Scan(ctx)
	if err != nil {
		return "", err
	}

	// Start building the Markdown content
	var markdownContent strings.Builder
	markdownContent.WriteString("# Weekly Newsletter\n\n")

	// Loop through the articles and append them to the Markdown content
	for _, article := range articles {
		markdownContent.WriteString(fmt.Sprintf("## %s\n", article.Title))
		markdownContent.WriteString(fmt.Sprintf("%s\n\n", article.Content))
		markdownContent.WriteString(fmt.Sprintf("[Read more](%s)\n\n", article.Link))
	}

	return markdownContent.String(), nil
}
