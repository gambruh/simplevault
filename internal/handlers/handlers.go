// Package handlers provides handlers logic for a http/https server
// It works with two storages - authentication and data storage
// It also provides NewService construction
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
	"github.com/gambruh/gophkeeper/internal/storage"
	"github.com/gambruh/gophkeeper/internal/storage/database"
)

// WebService is a class to
type WebService struct {
	Storage     Storage
	AuthStorage AuthStorage
	Mu          *sync.Mutex
}

// AuthStorage stores login and passwords of app users
// Different databases for authentication implementation can be used
type AuthStorage interface {
	Register(login string, password string) error
	VerifyCredentials(login string, password string) error
}

// Storage interface is a data storage. Implementation may vary
type Storage interface {
	SetLoginCred(username string, logincreds storage.EncryptedData) error
	SetNote(username string, note storage.EncryptedData) error
	SetBinary(username string, binary storage.Binary) error
	SetCard(username string, card storage.EncryptedData) error
	GetLoginCred(username string, name string) (storage.EncryptedData, error)
	GetNote(username string, name string) (storage.EncryptedData, error)
	GetBinary(username string, name string) (storage.Binary, error)
	GetCard(username string, name string) (storage.EncryptedData, error)
	ListLoginCreds(username string) ([]string, error)
	ListNotes(username string) ([]string, error)
	ListBinaries(username string) ([]string, error)
	ListCards(username string) ([]string, error)
}

var (
	ErrWrongCredentials = errors.New("wrong login/password")
	ErrUsernameIsTaken  = errors.New("username is taken")
)

// handlers for a server
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
		r.Post("/api/notes/add", h.AddNote)
		r.Post("/api/notes/get", h.GetNote)
		r.Get("/api/notes/list", h.ListNotes)
		r.Post("/api/binaries/add", h.AddBinary)
		r.Post("/api/binaries/get", h.GetBinary)
		r.Get("/api/binaries/list", h.ListBinaries)
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

// Register is a new user registration handler
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

// Login allows a user to login
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
	case auth.ErrUserNotFound:
		fmt.Println("Invalid login credentials:", data.Login)
		w.WriteHeader(http.StatusUnauthorized)
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

// AddLoginCreds sets new login credentials in a database
// returns error if there is already login with the same for a current user
func (h *WebService) AddLoginCreds(w http.ResponseWriter, r *http.Request) {
	var logincred storage.EncryptedData

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
		if database.IsUniqueConstraintViolation(err) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		log.Println("Unexpected case in AddLoginCreds Handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// AddCard sets new card data in a database
// returns http.StatusConflict error if there is already card with the same name for a current user
func (h *WebService) AddCard(w http.ResponseWriter, r *http.Request) {
	var carddata storage.EncryptedData

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
	default:
		if database.IsUniqueConstraintViolation(err) {
			w.WriteHeader(http.StatusConflict)
			return
		}
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
	case storage.ErrDataNotFound:
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
	case storage.ErrDataNotFound:
		w.WriteHeader(http.StatusNoContent)
	default:
		log.Println("error in ListCards handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *WebService) GetCard(w http.ResponseWriter, r *http.Request) {
	var carddata storage.EncryptedData
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

	card, err := h.Storage.GetCard(username.(string), carddata.Name)

	switch err {
	case nil:
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(card)
	case storage.ErrDataNotFound:
		w.WriteHeader(http.StatusNoContent)
	default:
		log.Println("error in GetCard handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *WebService) GetLoginCreds(w http.ResponseWriter, r *http.Request) {
	var input storage.EncryptedData

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
	case storage.ErrDataNotFound:
		w.WriteHeader(http.StatusNoContent)
	default:
		log.Println("error in GetLoginCreds handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *WebService) AddNote(w http.ResponseWriter, r *http.Request) {
	var note storage.EncryptedData

	contentType := r.Header.Get("Content-type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := r.Context().Value(config.UserID("userID"))

	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		log.Println("error in AddNote handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.Storage.SetNote(username.(string), note)
	switch err {
	case nil:
		w.WriteHeader(http.StatusAccepted)
	default:
		if database.IsUniqueConstraintViolation(err) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		log.Println("Unexpected case in AddNote Handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *WebService) GetNote(w http.ResponseWriter, r *http.Request) {
	var input storage.EncryptedData

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

	note, err := h.Storage.GetNote(username.(string), input.Name)

	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(note)
	case storage.ErrDataNotFound:
		w.WriteHeader(http.StatusNoContent)
	default:
		log.Println("error in GetNote handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *WebService) ListNotes(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(config.UserID("userID"))

	notes, err := h.Storage.ListNotes(username.(string))
	switch err {
	case nil:
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(notes)
	case storage.ErrDataNotFound:
		w.WriteHeader(http.StatusNoContent)
	default:
		log.Println("error in ListNotes handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *WebService) AddBinary(w http.ResponseWriter, r *http.Request) {
	var binary storage.Binary

	contentType := r.Header.Get("Content-type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := r.Context().Value(config.UserID("userID"))

	err := json.NewDecoder(r.Body).Decode(&binary)
	if err != nil {
		log.Println("error in AddBinary handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.Storage.SetBinary(username.(string), binary)
	switch err {
	case nil:
		w.WriteHeader(http.StatusAccepted)
	default:
		if database.IsUniqueConstraintViolation(err) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		log.Println("Unexpected case in AddBinary Handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *WebService) GetBinary(w http.ResponseWriter, r *http.Request) {
	var input storage.Binary

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

	binary, err := h.Storage.GetBinary(username.(string), input.Name)

	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(binary)
	case storage.ErrDataNotFound:
		w.WriteHeader(http.StatusNoContent)
	default:
		log.Println("error in GetNote handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *WebService) ListBinaries(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(config.UserID("userID"))

	notes, err := h.Storage.ListBinaries(username.(string))
	switch err {
	case nil:
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(notes)
	case storage.ErrDataNotFound:
		w.WriteHeader(http.StatusNoContent)
	default:
		log.Println("error in ListNotes handler:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
