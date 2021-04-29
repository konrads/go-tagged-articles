package db

import (
	"github.com/konrads/go-tagged-articles/pkg/model"
)

// Db interface allows substitution of implementation (and mocking)
type DB interface {
	GetArticle(id string) (*model.Article, error)
	SaveArticle(article *model.Article) (int, error)
	GetArticles(tag string, date model.Date) ([]*model.Article, error)
	Close() error
}
