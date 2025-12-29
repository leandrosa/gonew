package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"x.com/api/controllers"
	"x.com/api/middlewares"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	// middlewares
	r.Use(middlewares.PanicRecovery)
	r.Use(middlewares.WrapHandlerWithLogging)

	//index
	r.HandleFunc("/", controllers.Index).Methods(http.MethodGet)
	r.HandleFunc("/index", controllers.Index).Methods(http.MethodGet)

	//articles
	r.HandleFunc("/articles", controllers.GetAllArticles).Methods(http.MethodGet)
	r.HandleFunc("/articles/{id}", controllers.GetArticle).Methods(http.MethodGet)
	r.HandleFunc("/articles", controllers.AddArticle).Methods(http.MethodPost)
	r.HandleFunc("/articles/{id}", controllers.UpdateArticle).Methods(http.MethodPut)
	r.HandleFunc("/articles/{id}", controllers.DeleteArticle).Methods(http.MethodDelete)

	return r
}
