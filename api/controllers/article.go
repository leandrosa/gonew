package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"x.com/api/models"
)

var articles = []models.Article{
	{Id: "8617bf49-39a9-4268-b113-7b6bcd189ba2", Title: "Article 1", Desc: "Article Description 1", Content: "Article Content 1"},
	{Id: "38da7ce2-02b5-471a-90b8-c299f2ef132e", Title: "Article 2", Desc: "Article Description 2", Content: "Article Content 2"},
}

func GetAllArticles(w http.ResponseWriter, r *http.Request) {
	responseWithJSON(w, http.StatusOK, articles)
}

func GetArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for _, article := range articles {
		if article.Id == id {
			responseWithJSON(w, http.StatusOK, article)
			return
		}
	}

	responseWithJSON(w, http.StatusNotFound, nil)
}

func AddArticle(w http.ResponseWriter, r *http.Request) {
	// Read to request body
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}
	var article models.Article
	json.Unmarshal(body, &article)

	article.Id = (uuid.New()).String()
	articles = append(articles, article)

	//TODO: Add the location header
	responseWithJSON(w, http.StatusCreated, article)
}

func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Read request body
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var updatedArticle models.Article
	json.Unmarshal(body, &updatedArticle)

	for index, article := range articles {
		if article.Id == id {
			article.Title = updatedArticle.Title
			article.Desc = updatedArticle.Desc
			article.Content = updatedArticle.Content

			articles[index] = article

			responseWithJSON(w, http.StatusOK, article)
			break
		}
	}
}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for index, article := range articles {
		if article.Id == id {
			articles = append(articles[:index], articles[index+1:]...)
			break
		}
	}

	responseWithJSON(w, http.StatusNoContent, nil)
}

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
