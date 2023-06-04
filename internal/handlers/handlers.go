package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/gambruh/gophkeeper/internal/auth"
	"github.com/gambruh/gophkeeper/internal/config"
	"github.com/gambruh/gophkeeper/internal/database"
)

type WebService struct {
	Storage     Storage
	AuthStorage AuthStorage
	Mu          *sync.Mutex
}

type AuthStorage interface {
	Register(login string, password string) error
	VerifyCredentials(login string, password string) error
}

type Storage interface {
	SetLoginCred(username string, logincreds database.LoginCreds) error
	//	SetNote(username string, note Note) error
	//	SetBinary(username string, binary Binary) error
	SetCard(username string, card database.Card) error
	GetLoginCred(username string, name string) (database.LoginCreds, error)
	//	GetNote(username string, name string) (Note, error)
	//	GetBinary(username string, name string) (Binary, error)
	GetCard(username string, name string) (database.Card, error)
	ListLoginCreds(username string) ([]database.LoginCreds, error)
	//	ListNotes(username string) ([]Note, error)
	//	ListBinaries(username string) ([]Binary, error)
	ListCards(username string) ([]string, error)
}

var (
	ErrWrongCredentials = errors.New("wrong login/password")
	ErrUsernameIsTaken  = errors.New("username is taken")
)

func (h *WebService) Service() http.Handler {

	r := chi.NewRouter()
	r.Use(middleware.Compress(5, "text/plain", "text/html", "application/json"))

	r.Post("/api/user/register", h.Register)
	r.Post("/api/user/login", h.Login)

	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Post("/api/logincreds/add", h.AddLoginCreds)
		r.Post("/api/logincreds/get", h.GetLoginCreds)
		r.Get("/api/logincreds/list", h.ListLoginCreds)
		r.Post("/api/cards/add", h.AddCard)
		r.Post("/api/cards/get", h.GetCard)
		r.Get("/api/cards/list", h.ListCards)
	})

	return r
}

func NewService(storage Storage, authstorage AuthStorage) *WebService {
	return &WebService{
		Storage:     storage,
		AuthStorage: authstorage,
		Mu:          &sync.Mutex{},
	}
}

func (h *WebService) Register(w http.ResponseWriter, r *http.Request) {
	var data auth.LoginData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if data.Login == "" {
		w.Write([]byte("Empty login field"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if data.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty password field"))
		return
	}

	err = h.AuthStorage.Register(data.Login, data.Password)
	switch err {
	case auth.ErrUsernameIsTaken:
		fmt.Println("Username is taken")
		w.WriteHeader(http.StatusConflict)
		return
	case nil:
		// Generate token
		token, err := auth.GenerateToken(data.Login)
		if err != nil {
			fmt.Println("error when generating token", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Set the token in "Cookies"
		http.SetCookie(w, &http.Cookie{
			Name:  "gophkeeper-auth",
			Value: token,
		})
		// Return a success response
		w.WriteHeader(http.StatusOK)
	default:
		log.Println("Unexpected case in new user registration:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *WebService) Login(w http.ResponseWriter, r *http.Request) {
	var data auth.LoginData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("Wrong login credentials format:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Verify the user's credentials
	err = h.AuthStorage.VerifyCredentials(data.Login, data.Password)
	switch err {
	case nil:
		//login and password are verified
	case auth.ErrWrongPassword:
		fmt.Println("Invalid login credentials:", data.Login)
		w.WriteHeader(http.StatusUnauthorized)
		return
	default:
		fmt.Println("error when verifying login credentials:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Generate a token
	token, err := auth.GenerateToken(data.Login)
	if err != nil {
		fmt.Println("error when generating token", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the token as a cookie in the response
	http.SetCookie(w, &http.Cookie{
		Name:  "gophkeeper-auth",
		Value: token,
	})

	// Return a success response
	w.WriteHeader(http.StatusOK)
}

func (h *WebService) AddLoginCreds(w http.ResponseWriter, r *http.Request) {
	var logincred database.LoginCreds

	contentType := r.Header.Get("Content-type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := r.Context().Value(config.UserID("userID"))

	err := json.NewDecoder(r.Body).Decode(&logincred)
	if err != nil {
		log.Println("error in AddLoginCreds handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.Storage.SetLoginCred(username.(string), logincred)
	switch err {
	case nil:
		w.WriteHeader(http.StatusAccepted)
	default:
		log.Println("Unexpected case in AddLoginCreds Handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *WebService) AddCard(w http.ResponseWriter, r *http.Request) {
	var carddata database.Card

	contentType := r.Header.Get("Content-type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := r.Context().Value(config.UserID("userID"))

	err := json.NewDecoder(r.Body).Decode(&carddata)
	if err != nil {
		fmt.Println("error in AddCard handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.Storage.SetCard(username.(string), carddata)
	switch err {
	case nil:
		w.WriteHeader(http.StatusAccepted)
	case database.ErrMetanameIsTaken:
		w.WriteHeader(http.StatusConflict)
	default:
		fmt.Println("Unexpected case in AddCard Handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *WebService) ListCards(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(config.UserID("userID"))

	cards, err := h.Storage.ListCards(username.(string))
	switch err {
	case nil:
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cards)
	case database.ErrDataNotFound:
		w.WriteHeader(http.StatusNoContent)
	default:
		log.Println("error in ListCards handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *WebService) ListLoginCreds(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(config.UserID("userID"))

	logins, err := h.Storage.ListLoginCreds(username.(string))
	switch err {
	case nil:
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(logins)
	case database.ErrDataNotFound:
		w.WriteHeader(http.StatusNoContent)
	default:
		log.Println("error in ListCards handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *WebService) GetCard(w http.ResponseWriter, r *http.Request) {
	var carddata database.Card
	username := r.Context().Value(config.UserID("userID"))

	contentType := r.Header.Get("Content-type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&carddata)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	card, err := h.Storage.GetCard(username.(string), carddata.Cardname)

	switch err {
	case nil:
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(card)
	case database.ErrDataNotFound:
		w.WriteHeader(http.StatusNoContent)
	default:
		log.Println("error in GetCard handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *WebService) GetLoginCreds(w http.ResponseWriter, r *http.Request) {
	var input database.LoginCreds

	username := r.Context().Value(config.UserID("userID"))

	contentType := r.Header.Get("Content-type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logincreds, err := h.Storage.GetLoginCred(username.(string), input.Name)

	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(logincreds)
	case database.ErrDataNotFound:
		w.WriteHeader(http.StatusNoContent)
	default:
		log.Println("error in GetCard handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
