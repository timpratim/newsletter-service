package main

import (
	"context"
	"github.com/mmcdole/gofeed"
	"github.com/uptrace/bun"
)

func fetchAndStoreArticles(ctx context.Context, db *bun.DB, feedURL string) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return err
	}

	for _, item := range feed.Items {
		count, err := db.NewSelect().Model((*Article)(nil)).
			Where("link = ?", item.Link).
			Count(ctx)
		if err != nil {
			return err
		}

		if count == 0 {
			article := Article{
				Title:       item.Title,
				Link:        item.Link,
				PublishedAt: *item.PublishedParsed,
				Content:     item.Content,
			}

			_, err := db.NewInsert().Model(&article).Exec(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
