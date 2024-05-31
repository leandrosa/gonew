package controllers

import (
	"encoding/json"
	"net/http"
)

type IndexApi struct {
	Message string `json:"message"`
}

func Index(w http.ResponseWriter, r *http.Request) {

	indexApi := IndexApi{
		Message: "Welcome",
	}

	respondWithJSON(w, http.StatusOK, indexApi)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
