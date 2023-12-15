package main

import (
	"github.com/uptrace/bun"
	"time"
)

type Article struct {
	bun.BaseModel `bun:"table:articles,alias:a"`
	ID            int64     `bun:"id,pk,autoincrement"`
	Title         string    `bun:"title,notnull"`
	Link          string    `bun:"link,notnull"`
	PublishedAt   time.Time `bun:"published_at"`
	Content       string    `bun:"content,nullzero"`
	Summarized    bool      `bun:"summarized,default:false"`
}

// Newsletter is a Bun model.This is a struct that represents a database table. Newsletter contains the fields that are mapped to the table columns.
// bun.BaseModel is a struct that contains the common fields for all models. It is not required, but it is convenient to embed it into your models.
type Newsletter struct {
	bun.BaseModel `bun:"table:newsletters,alias:n"`
	ID            int64     `bun:"id,pk,autoincrement"`
	IssueDate     time.Time `bun:"issue_date,notnull"`
	Content       string    `bun:"content,nullzero"`
}
