package repository

import (
	"context"
	"fmt"

	"x.com/api/models"
)

type ArticleRepositoryInMemory struct {
	store map[string]*models.Article
}

func NewArticleRepositoryInMemory() *ArticleRepositoryInMemory {
	return &ArticleRepositoryInMemory{
		store: map[string]*models.Article{},
	}
}

func (r ArticleRepositoryInMemory) AddArticle(_ context.Context, article *models.Article) error {
	r.store[article.Id] = article

	return nil
}

func (r ArticleRepositoryInMemory) GetByID(_ context.Context, id string) (*models.Article, error) {

	article, ok := r.store[id]
	if !ok {
		return article, fmt.Errorf("id not found; id: %s", id)
	}

	return article, nil
}

func (r ArticleRepositoryInMemory) GeAll(_ context.Context) ([]*models.Article, error) {
	article := []*models.Article{}
	for _, a := range r.store {
		article = append(article, a)
	}

	return article, nil
}

func (r ArticleRepositoryInMemory) UpdateArticle(_ context.Context, article *models.Article) error {
	r.store[article.Id] = article
	return nil
}

func (r ArticleRepositoryInMemory) DeleteArticle(_ context.Context, id string) error {
	delete(r.store, id)
	return nil
}
