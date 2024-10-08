package database

import (
	"context"
	"errors"

	"github.com/rezaAmiri123/scalable-backend/internal/entity"
	"github.com/rezaAmiri123/scalable-backend/internal/promhelper"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (g *GormDatabase) CreateAuthor(ctx context.Context, author *entity.Author) error {
	return g.queryMetric.Do("create_author", func() error {
		err := g.db.WithContext(ctx).Create(author).Error
		if err != nil {
			logrus.WithError(err).Error("error while creating an author")
			return err
		}
		return nil
	})
}

func (g *GormDatabase) GetAuthor(ctx context.Context, id uint) (entity.Author, error) {
	var author entity.Author
	return author, g.queryMetric.Do("get_author", func() error {
		err := g.db.WithContext(ctx).Where("id=?", id).First(&author).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return promhelper.NewPromError(promhelper.StatusNotFound, ErrEntityNotfound)
			}
			return err
		}
		return nil
	})
}
