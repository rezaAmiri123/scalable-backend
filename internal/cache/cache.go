package cache

import (
	"context"

	"github.com/rezaAmiri123/scalable-backend/internal/entity"
)

type Cache interface{
	TagArticles(ctx context.Context, tagSlug string)([]entity.Article,error)
}
