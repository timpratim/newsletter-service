package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

const dsn = "postgres://test:test123@localhost:5432/newsletter_service?sslmode=disable"

func main() {
	db := setupDatabase()
	defer db.Close()

	router := gin.Default()
	setupRoutes(router, db)

	router.Run(":8080")
}

func setupDatabase() *bun.DB {
	// Open a PostgreSQL database.
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Create a Bun db on top of it.
	db := bun.NewDB(pgdb, pgdialect.New())

	// Database migration
	ctx := context.Background()
	_, err := db.NewCreateTable().Model((*Article)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		log.Fatalf("Failed to create articles table: %v", err)
	}

	_, err = db.NewCreateTable().Model((*Newsletter)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		log.Fatalf("Failed to create newsletters table: %v", err)
	}

	return db
}

func setupRoutes(router *gin.Engine, db *bun.DB) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the Newsletter Service"})
	})

	router.GET("/fetch_articles", func(c *gin.Context) {
		ctx := context.Background()
		feedURL := c.Query("feed_url")
		if feedURL == "" {
			c.JSON(400, gin.H{"error": "Feed URL is required"})
			return
		}

		err := fetchAndStoreArticles(ctx, db, feedURL)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch articles"})
			return
		}

		c.JSON(200, gin.H{"message": "Articles fetched successfully"})
	})

	router.GET("/summarize_articles", func(c *gin.Context) {
		ctx := context.Background()
		err := summarizeArticles(ctx, db)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to summarize articles"})
			return
		}

		c.JSON(200, gin.H{"message": "Articles summarized successfully"})
	})

	router.GET("/generate_newsletter", func(c *gin.Context) {
		ctx := context.Background()
		newsletterContent, err := generateMarkdownNewsletter(ctx, db)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate newsletter"})
			return
		}

		c.JSON(200, gin.H{"newsletter": newsletterContent})
	})

	// Add more routes as needed
}
