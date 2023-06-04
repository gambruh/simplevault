package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/gambruh/gophkeeper/internal/config"
)

type Storage interface {
	SetLoginCred(username string, logincreds LoginCreds) error
	//	SetNote(username string, note Note) error
	//	SetBinary(username string, binary Binary) error
	SetCard(username string, card Card) error
	GetLoginCred(username string, name string) (LoginCreds, error)
	//	GetNote(username string, name string) (Note, error)
	//	GetBinary(username string, name string) (Binary, error)
	GetCard(username string, name string) (Card, error)
	ListLoginCreds(username string) ([]LoginCreds, error)
	//	ListNotes(username string) ([]Note, error)
	//	ListBinaries(username string) ([]Binary, error)
	ListCards(username string) ([]string, error)
}

type SQLdb struct {
	DB *sql.DB
}

// типы ошибок
var (
	ErrTableDoesntExist = errors.New("table doesn't exist")
	ErrWrongPassword    = errors.New("wrong password")
	ErrDataNotFound     = errors.New("requested data not found in storage")
	ErrMetanameIsTaken  = errors.New("metaname is already in use")
)

func NewSQLdb(postgresStr string) *SQLdb {
	DB, _ := sql.Open("postgres", postgresStr)
	return &SQLdb{
		DB: DB,
	}
}

func GetDB() (defstorage Storage) {

	db := NewSQLdb(config.Cfg.Database)
	err := db.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defstorage = db

	return defstorage
}

func (s *SQLdb) CheckConn(dbAddress string) error {
	db, err := sql.Open("postgres", config.Cfg.Database)
	if err != nil {
		fmt.Printf("error while opening DB:%v\n", err)
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		fmt.Printf("error while pinging: %v\n", err)
		return err
	}
	return nil
}

// InitDatabase creates tables in an empty database
func (s *SQLdb) InitDatabase() error {
	err := s.createLoginCredsTable()
	if err != nil {
		return fmt.Errorf("error creating logincreds table:%w", err)
	}
	err = s.createCardsTable()
	if err != nil {
		return fmt.Errorf("error creating cards table:%w", err)
	}
	err = s.createNotesTable()
	if err != nil {
		return fmt.Errorf("error creating notes table:%w", err)
	}
	err = s.createBinariesTable()
	if err != nil {
		return fmt.Errorf("error creating binaries table:%w", err)
	}

	return nil
}

func (s *SQLdb) checkTableExists(tablename string) error {
	var check bool
	err := s.DB.QueryRow(checkTableExistsQuery, tablename).Scan(&check)
	if err != nil {
		fmt.Printf("error checking if table exists: %v", err)
		return err
	}
	if !check {
		return ErrTableDoesntExist
	}
	return nil
}

func (s *SQLdb) createLoginCredsTable() error {
	err := s.checkTableExists("gk_logincreds")
	if err == ErrTableDoesntExist {
		_, err := s.DB.Exec(createLoginCredsTableQuery)
		return err
	}
	return nil
}

func (s *SQLdb) createCardsTable() error {
	err := s.checkTableExists("gk_cards")
	if err == ErrTableDoesntExist {
		_, err := s.DB.Exec(createCardsTableQuery)
		return err
	}
	return nil

}

func (s *SQLdb) createNotesTable() error {
	err := s.checkTableExists("gk_notes")
	if err == ErrTableDoesntExist {
		_, err := s.DB.Exec(createNotesTableQuery)
		return err
	}
	return nil

}

func (s *SQLdb) createBinariesTable() error {
	err := s.checkTableExists("gk_binaries")
	if err == ErrTableDoesntExist {
		_, err := s.DB.Exec(createBinariesTableQuery)
		return err
	}
	return nil
}

func (s *SQLdb) SetCard(username string, cardData Card) error {

	var cardname string
	err := s.DB.QueryRow(checkCardNameQuery, cardData.Cardname).Scan(&cardname)
	if err != sql.ErrNoRows {
		return ErrMetanameIsTaken
	}

	_, err = s.DB.Exec(setCardQuery, cardData.Cardname, cardData.Number, cardData.Name, cardData.Surname, cardData.ValidTill, cardData.Code, username)
	if err != nil {
		return fmt.Errorf("error setting card data in SetCard:%w", err)
	}
	return nil
}

func (s *SQLdb) SetLoginCred(username string, loginData LoginCreds) error {
	_, err := s.DB.Exec(setLoginCredsQuery, loginData.Name, loginData.Login, loginData.Password, loginData.Site, username)
	if err != nil {
		return fmt.Errorf("error setting card data in SetLoginCred:%w", err)
	}
	return nil
}

func (s *SQLdb) GetCard(username string, cardname string) (Card, error) {
	var cardData Card
	err := s.DB.QueryRow(getCardQuery, cardname, username).Scan(&cardData.Cardname, &cardData.Number, &cardData.Name, &cardData.Surname, &cardData.ValidTill, &cardData.Code)
	if err != nil {
		if err == sql.ErrNoRows {
			return Card{}, ErrDataNotFound
		} else {
			return Card{}, fmt.Errorf("error in GetCard:%w", err)
		}
	}
	return cardData, nil
}

func (s *SQLdb) ListCards(username string) (cardnames []string, err error) {

	rows, err := s.DB.Query(listCardsQuery, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrDataNotFound
		}
		return nil, fmt.Errorf("couldn't ask database in ListCards:%w", err)
	}

	for rows.Next() {
		var cardname string

		err := rows.Scan(&cardname)
		if err != nil {
			return nil, fmt.Errorf("error scanning in ListCards:%w", err)
		}
		cardnames = append(cardnames, cardname)

	}

	err = rows.Err()

	if err != nil {
		return nil, fmt.Errorf("error scanning with rows.Next() in ListCards:%w", err)
	}

	return cardnames, nil
}

func (s *SQLdb) GetLoginCred(username string, loginname string) (LoginCreds, error) {

	return LoginCreds{}, nil
}

func (s *SQLdb) ListLoginCreds(username string) (logincreds []LoginCreds, err error) {

	return nil, nil
}
