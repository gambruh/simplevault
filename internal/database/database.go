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
	SetLoginCred(username string, logincreds EncryptedData) error
	//	SetNote(username string, note Note) error
	//	SetBinary(username string, binary Binary) error
	SetCard(username string, card EncryptedCard) error
	GetLoginCred(username string, name string) (EncryptedData, error)
	//	GetNote(username string, name string) (EncryptedData, error)
	//	GetBinary(username string, name string) (EncryptedData, error)
	GetCard(username string, name string) (EncryptedCard, error)
	ListLoginCreds(username string) ([]string, error)
	//	ListNotes(username string) ([]string, error)
	//	ListBinaries(username string) ([]string, error)
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
		if _, err := s.DB.Exec(createLoginCredsTableQuery); err != nil {
			return err
		}
		if _, err = s.DB.Exec(alterTableLC); err != nil {
			return err
		}
		if _, err = s.DB.Exec(alterTableLC2); err != nil {
			return err
		}
		if _, err = s.DB.Exec(createUniqueLoginCredsConstraint); err != nil {
			return err
		}

	}
	return nil
}

func (s *SQLdb) createCardsTable() error {
	err := s.checkTableExists("gk_cards")
	if err == ErrTableDoesntExist {
		_, err := s.DB.Exec(createCardsTableQuery)
		if err != nil {
			return err
		}
		_, err = s.DB.Exec(createUniqueCardConstraint)
		if err != nil {
			return err
		}
		_, err = s.DB.Exec(encryptedCardsTable)
		if err != nil {
			return err
		}
		_, err = s.DB.Exec(addColumnData)
		if err != nil {
			return err
		}
	}
	return nil

}

func (s *SQLdb) createNotesTable() error {
	err := s.checkTableExists("gk_notes")
	if err == ErrTableDoesntExist {
		_, err := s.DB.Exec(createNotesTableQuery)
		if err != nil {
			return err
		}
		_, err = s.DB.Exec(createUniqueNoteConstraint)
		if err != nil {
			return err
		}
		_, err = s.DB.Exec(alterTableNotes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SQLdb) createBinariesTable() error {
	err := s.checkTableExists("gk_binaries")
	if err == ErrTableDoesntExist {
		if _, err := s.DB.Exec(createBinariesTableQuery); err != nil {
			return err
		}
		if _, err := s.DB.Exec(createUniqueBinaryConstraint); err != nil {
			return err
		}
	}

	return nil
}

func (s *SQLdb) SetCard(username string, cardData EncryptedCard) error {

	var cardname string
	err := s.DB.QueryRow(checkCardNameQuery, cardData.Cardname, username).Scan(&cardname)
	if err != sql.ErrNoRows {
		return ErrMetanameIsTaken
	}

	_, err = s.DB.Exec(setCardQuery, cardData.Cardname, cardData.Data, username)
	if err != nil {
		return fmt.Errorf("error setting card data in SetCard:%w", err)
	}
	return nil
}

func (s *SQLdb) SetLoginCred(username string, loginData EncryptedData) error {
	_, err := s.DB.Exec(setLoginCredsQuery, loginData.Name, loginData.Data, username)
	if err != nil {
		return fmt.Errorf("error setting data in SetLoginCred:%w", err)
	}
	return nil
}

func (s *SQLdb) GetCard(username string, cardname string) (EncryptedCard, error) {
	var cardData EncryptedCard
	err := s.DB.QueryRow(getCardQuery, cardname, username).Scan(&cardData.Cardname, &cardData.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return EncryptedCard{}, ErrDataNotFound
		} else {
			return EncryptedCard{}, fmt.Errorf("error in GetCard:%w", err)
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

func (s *SQLdb) GetLoginCred(username string, loginname string) (logincred EncryptedData, err error) {
	var encrData EncryptedData
	err = s.DB.QueryRow(getLoginCredsQuery, loginname, username).Scan(&encrData.Name, &encrData.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return EncryptedData{}, ErrDataNotFound
		} else {
			return EncryptedData{}, fmt.Errorf("error in GetLoginCred:%w", err)
		}
	}
	return encrData, nil
}

func (s *SQLdb) ListLoginCreds(username string) (logincrednames []string, err error) {

	rows, err := s.DB.Query(listLoginCredsQuery, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrDataNotFound
		}
		return nil, fmt.Errorf("couldn't ask database in ListLoginCreds:%w", err)
	}

	for rows.Next() {
		var itemname string

		err := rows.Scan(&itemname)
		if err != nil {
			return nil, fmt.Errorf("error scanning in ListLoginCreds:%w", err)
		}
		logincrednames = append(logincrednames, itemname)

	}

	err = rows.Err()

	if err != nil {
		return nil, fmt.Errorf("error scanning with rows.Next() in ListLoginCreds:%w", err)
	}

	return logincrednames, nil
}
