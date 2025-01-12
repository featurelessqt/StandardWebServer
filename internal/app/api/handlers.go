package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/standardWebServer/internal/app/middleware"
	"github.com/standardWebServer/internal/app/models"
)

type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_srror"`
}

func initHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

func (api *API) GetAllArticles(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
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

func (api *API) PostArticle(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Post Article POST /api/v1/articles")
	var article models.Article
	err := json.NewDecoder(req.Body).Decode(&article)
	if err != nil {
		api.logger.Info("Invalid json recieved from client")
		msg := Message{
			StatusCode: 400,
			Message:    "Provided json is invalid",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	a, err := api.storage.Article().Create(&article)
	if err != nil {
		api.logger.Info("Troubles while creating new article: ", err)
		msg := Message{
			StatusCode: 501,
			Message:    "We have some troubles to accessing database. Try again.",
			IsError:    true,
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
	}
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(a)
}

func (api *API) GetArticlesById(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Get Article by ID /api.v1/articlew/{id}")
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		api.logger.Info("troubles while parsing {id} parse:", err)
		msg := Message{
			StatusCode: 400,
			Message:    "Unapropriate id value. dont use ID as uncasting to int value.",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	article, ok, err := api.storage.Article().FindArticleById(id)
	if err != nil {
		api.logger.Info("Troubles while accessing databse table (articles) tith id. err:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles to accessing database. Try again",
			IsError:    true,
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	if !ok {
		api.logger.Info("Can not find article with that id in database")
		msg := Message{
			StatusCode: 404,
			Message:    "Article with that ID does not esists in database.",
			IsError:    true,
		}

		writer.WriteHeader(404)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(article)
}

func (api *API) DeleteArticleById(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Delete Article by Id DELETE /api/v1/articles/{id}")
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		api.logger.Info("Troubles while parsing {id} param:", err)
		msg := Message{
			StatusCode: 400,
			Message:    "Unapropriate id value. don't use ID as uncasting to int value.",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	_, ok, err := api.storage.Article().FindArticleById(id)
	if err != nil {
		api.logger.Info("Troubles while accessing database table (articles) with id. err:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles to accessing database. Try again",
			IsError:    true,
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	if !ok {
		api.logger.Info("Can not find article with that ID in database")
		msg := Message{
			StatusCode: 404,
			Message:    "Article with that ID does not exists in database.",
			IsError:    true,
		}

		writer.WriteHeader(404)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	_, err = api.storage.Article().DeleteById(id)
	if err != nil {
		api.logger.Info("Troubles while deleting database elemnt from table (articles) with id. err:", err)
		msg := Message{
			StatusCode: 501,
			Message:    "We have some troubles to accessing database. Try again",
			IsError:    true,
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	writer.WriteHeader(202)
	msg := Message{
		StatusCode: 202,
		Message:    fmt.Sprintf("Article with ID %d successfully deleted.", id),
		IsError:    false,
	}
	json.NewEncoder(writer).Encode(msg)
}

func (api *API) PostUserRegister(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Post User Register POST /api/v1/user/register")
	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		api.logger.Info("Invalid json recieved from client")
		msg := Message{
			StatusCode: 400,
			Message:    "Provided json is invalid",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	_, ok, err := api.storage.User().FindByLogin(user.Login)
	if err != nil {
		api.logger.Info("Troubles while accessing database table (users) with id. err:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles to accessing database. Try again",
			IsError:    true,
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	if ok {
		api.logger.Info("User with that ID already exists")
		msg := Message{
			StatusCode: 400,
			Message:    "User with that login already exists in database",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	//Теперь пытаемся добавить в бд
	userAdded, err := api.storage.User().Create(&user)
	if err != nil {
		api.logger.Info("Troubles while accessing database table (users) with id. err:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles to accessing database. Try again",
			IsError:    true,
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	msg := Message{
		StatusCode: 201,
		Message:    fmt.Sprintf("User {login:%s} successfully registered!", userAdded.Login),
		IsError:    false,
	}
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(msg)
}

func (api *API) PostToAuth(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Post to Auth POST /api/v1/user/auth")
	var userFromJson models.User
	err := json.NewDecoder(req.Body).Decode(&userFromJson)
	if err != nil {
		api.logger.Info("Invalid json recieved from client")
		msg := Message{
			StatusCode: 400,
			Message:    "Provided json is invalid",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	userInDB, ok, err := api.storage.User().FindByLogin(userFromJson.Login)
	if err != nil {
		api.logger.Info("Can not make user searcj in database:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles while accessing database",
			IsError:    true,
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	if !ok {
		api.logger.Info("User with that login does not exists")
		msg := Message{
			StatusCode: 400,
			Message:    "User with that login already exists in database. Try register first",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	if userInDB.Password != userFromJson.Password {
		api.logger.Info("Invalid creditals to auth")
		msg := Message{
			StatusCode: 404,
			Message:    "Your password is invalid",
			IsError:    true,
		}
		writer.WriteHeader(404)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	claims["admin"] = true
	claims["name"] = userInDB.Login
	tokenString, err := token.SignedString(middleware.SecretKey)
	if err != nil {
		api.logger.Info("Can not claim jwt-token")
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles. Try again",
			IsError:    true,
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	msg := Message{
		StatusCode: 201,
		Message:    tokenString,
		IsError:    false,
	}
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(msg)
}
