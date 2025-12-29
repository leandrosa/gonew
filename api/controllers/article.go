package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"x.com/api/models"
	"x.com/api/repository"
)

func init() {
	articleRepo = repository.NewArticleRepositoryInMemory()

	var articles = []*models.Article{
		{Id: "8617bf49-39a9-4268-b113-7b6bcd189ba2", Title: "Article 1", Desc: "Article Description 1", Content: "Article Content 1"},
		{Id: "38da7ce2-02b5-471a-90b8-c299f2ef132e", Title: "Article 2", Desc: "Article Description 2", Content: "Article Content 2"},
	}

	for _, article := range articles {
		articleRepo.AddArticle(context.Background(), article)
	}
}

type articleRepository interface {
	AddArticle(_ context.Context, article *models.Article) error
	GetByID(_ context.Context, id string) (*models.Article, error)
	GeAll(_ context.Context) ([]*models.Article, error)
	UpdateArticle(_ context.Context, article *models.Article) error
	DeleteArticle(_ context.Context, id string) error
}

var articleRepo articleRepository

func GetAllArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := articleRepo.GeAll(r.Context())
	if err != nil {
		responseWithJSON(w, http.StatusInternalServerError, models.NewErrorResponse(err))
		return
	}

	responseWithJSON(w, http.StatusOK, articles)
}

func GetArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	article, err := articleRepo.GetByID(r.Context(), id)
	if err != nil {
		responseWithJSON(w, http.StatusInternalServerError, models.NewErrorResponse(err))
		return
	}

	if article == nil {
		models.NewErrorResponse(fmt.Errorf("article id not found: %s", id))
		responseWithJSON(w, http.StatusNotFound, nil)
		return
	}

	responseWithJSON(w, http.StatusOK, article)
}

func AddArticle(w http.ResponseWriter, r *http.Request) {
	// Read request body
	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		responseWithJSON(w, http.StatusBadRequest, models.NewErrorResponse(err))
		return
	}

	article.Id = (uuid.New()).String()
	articleRepo.AddArticle(r.Context(), &article)

	//TODO: Add the location header
	responseWithJSON(w, http.StatusCreated, article)
}

func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Read to request body
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responseWithJSON(w, http.StatusBadRequest, models.NewErrorResponse(err))
		return
	}

	var updatedArticle models.Article
	json.Unmarshal(body, &updatedArticle)

	article, err := articleRepo.GetByID(r.Context(), id)
	if err != nil {
		responseWithJSON(w, http.StatusInternalServerError, models.NewErrorResponse(err))
		return
	}

	if article == nil {
		responseWithJSON(w, http.StatusNotFound, models.NewErrorResponse(fmt.Errorf("article id not found: %s", id)))
		return
	}

	article.Title = updatedArticle.Title
	article.Desc = updatedArticle.Desc
	article.Content = updatedArticle.Content

	if err = articleRepo.UpdateArticle(r.Context(), article); err != nil {
		responseWithJSON(w, http.StatusNotFound, models.NewErrorResponse(fmt.Errorf("article id not found: %s", id)))
		return
	}

	responseWithJSON(w, http.StatusOK, article)
}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := articleRepo.DeleteArticle(r.Context(), id); err != nil {
		responseWithJSON(w, http.StatusBadRequest, models.NewErrorResponse(err))
		return
	}

	responseWithJSON(w, http.StatusNoContent, nil)
}

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
