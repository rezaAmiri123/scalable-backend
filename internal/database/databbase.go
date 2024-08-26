package database

import (
	"context"
	"errors"

	"github.com/rezaAmiri123/scalable-backend/internal/entity"
)
var(
	ErrEntityNotfound = errors.New("entity not found")
)

type Database interface{
	Migrate()error
	GetAuthor(ctx context.Context,id uint)(entity.Author,error)
	CreateAuthor(ctx context.Context, author *entity.Author)error

	ListTags(ctx context.Context)([]entity.Tag,error)
	ListTagArticles(ctx context.Context, slug string)([]entity.Article,error)
	CreateTag(ctx context.Context, tag *entity.Tag)error

	ListArticles(ctx context.Context) ([]entity.Article, error)
	CreateArticle(ctx context.Context, article *entity.Article, tags []string) error
}
