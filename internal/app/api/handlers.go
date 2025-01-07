package api

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_srror"`
}

func initHeasers(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

func (api *API) GetAllArticles(writer http.ResponseWriter, req *http.Request) {
	initHeasers(writer)
	api.logger.Info("Get ALl Articles Get /api/v1/articles")
	articles, err := api.storage.Article().SelectAll()
	if err != nil {
		api.logger.Info("Error While Articles.SelectAll : ", err)
		msg := Message{
			StatusCode: 501,
			Message:    "We have some troubles to accessing database. Try again later",
			IsError:    true,
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(articles)
}

func (api *API) GetArticlesById(writer http.ResponseWriter, req *http.Request) {}

func (api *API) DeleteArticleById(writer http.ResponseWriter, req *http.Request) {}

func (api *API) PostArticle(writer http.ResponseWriter, req *http.Request) {}

func (api *API) PostUserRegister(writer http.ResponseWriter, req *http.Request) {}
