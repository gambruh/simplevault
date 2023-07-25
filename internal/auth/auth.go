// Package auth provides authentication and authorization functions
// it is separated in case if other auth will be used
package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gambruh/gophkeeper/internal/argon2id"
	"github.com/gambruh/gophkeeper/internal/config"
)

const (
	userstablename     = "gk_users"
	passwordstablename = "gk_passwords"
)

type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthStorage interface {
	Register(login string, password string) error
	VerifyCredentials(login string, password string) error
}

type AuthMemStorage struct {
	Data map[string]string
}

type AuthDB struct {
	db *sql.DB
}

// Authentication errors
var (
	ErrUserNotFound     = errors.New("user not found in database")
	ErrTableDoesntExist = errors.New("table doesn't exist")
	ErrUsernameIsTaken  = errors.New("username is taken")
	ErrWrongCredentials = errors.New("wrong login credentials")
	ErrWrongPassword    = errors.New("wrong password")
)

func GenerateToken(login string) (string, error) {
	// Create a new token object, specifying the signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": login,
		"exp":    time.Now().Add(time.Hour * 8).Unix(),
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(config.Cfg.Key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// AuthMiddleware checks cookies attached to http request
// If not valid then http.StatusUnauthorized will return
// If valid it will pass the request to the handler
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		type MyCustomClaims struct {
			UserID string `json:"userID"`
			jwt.StandardClaims
		}

		cookie, err := r.Cookie("gophkeeper-auth")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Cfg.Key), nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*MyCustomClaims)

		if !ok || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), config.UserID("userID"), claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// NewAuthDB returns connection to authDB
func NewAuthDB(postgresStr string) *AuthDB {
	db, _ := sql.Open("postgres", postgresStr)
	return &AuthDB{
		db: db,
	}
}

// GetAuthDB returns new auth storage
func GetAuthDB() (authstorage AuthStorage) {

	db := NewAuthDB(config.Cfg.Database)
	db.InitAuthDB()
	authstorage = db

	return authstorage
}

// CheckTableExists is a helper function to check if there is already a table with specifed name in a database
func (s *AuthDB) CheckTableExists(tablename string) error {
	var check bool

	err := s.db.QueryRow(checkTableExistsQuery, tablename).Scan(&check)
	if err != nil {
		log.Printf("error checking if table exists: %v", err)
		return err
	}
	if !check {
		return ErrTableDoesntExist
	}
	return nil
}

// InitAuthDB creates tables in a fresh database.
func (s *AuthDB) InitAuthDB() error {
	err := s.CreateUsersTable()
	if err != nil {
		return err
	}
	err = s.CreatePasswordsTable()
	if err != nil {
		return err
	}
	return nil
}

// CreateUsersTable creates table "users"
func (s *AuthDB) CreateUsersTable() error {
	err := s.CheckTableExists(userstablename)
	if err == ErrTableDoesntExist {
		_, err := s.db.Exec(createUsersTableQuery)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreatePasswordsTable creates table "passwords"
func (s *AuthDB) CreatePasswordsTable() error {
	err := s.CheckTableExists(passwordstablename)
	if err == ErrTableDoesntExist {
		_, err = s.db.Exec(createPasswordsTableQuery)
		if err != nil {
			log.Println("error when creating passwords table:", err)
			return err
		}
	}
	return nil
}

// Register attempts to save login and password hash in a database
// returns error in case if the login is already exists in the database
func (s *AuthDB) Register(login string, password string) error {
	var username string
	hashedpassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Println("error when trying to hash password:", err)
		return err
	}
	e := s.db.QueryRow(CheckUsernameQuery, login).Scan(&username)
	switch e {
	case sql.ErrNoRows:
		_, err := s.db.Exec(AddNewUserQuery, login, hashedpassword)
		return err
	case nil:
		return ErrUsernameIsTaken
	default:
		fmt.Printf("unexpected error in Register user: %v\n", e)
		return e
	}
}

// VerifyCredentials compares login and password with existing in the database login and password hash
// returns error if no coincidence found, or if password hashes didn't match
func (s *AuthDB) VerifyCredentials(login string, password string) error {
	var (
		id   int
		pass string
	)

	err := s.db.QueryRow(CheckUsernameQuery, login).Scan(&id)
	if err != nil {
		return ErrUserNotFound
	}

	err = s.db.QueryRow(CheckPasswordQuery, id).Scan(&pass)
	if err != nil {
		return ErrWrongPassword
	}

	check, err := argon2id.ComparePasswordAndHash(password, pass)
	if err != nil {
		log.Println("error when trying to compare password and hash:", err)
		return err
	}
	if !check {
		return ErrWrongPassword
	}

	return nil
}

// Register is a method for inmemory implementation of AuthStorage interface
func (s *AuthMemStorage) Register(login string, password string) error {
	_, contains := s.Data[login]
	if contains {
		return ErrUsernameIsTaken
	}
	s.Data[login] = password
	return nil
}

// VerifyCredentials is a method for inmemory implementation of AuthStorage interface
func (s AuthMemStorage) VerifyCredentials(login string, password string) error {
	check, contains := s.Data[login]
	if contains && check == password {
		return nil
	}
	return ErrWrongPassword
}

// NewMemStorage returns inmemory implementation of AuthStorage interface
func NewMemStorage() *AuthMemStorage {
	return &AuthMemStorage{
		Data: make(map[string]string),
	}
}
