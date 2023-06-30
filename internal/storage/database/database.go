// package database is responsible for interaction with postgres DB
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"

	"github.com/gambruh/gophkeeper/internal/config"
	"github.com/gambruh/gophkeeper/internal/storage"
)

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

type SQLdb struct {
	DB *sql.DB
}

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
		return storage.ErrTableDoesntExist
	}
	return nil
}

func (s *SQLdb) createLoginCredsTable() error {
	err := s.checkTableExists("gk_logincreds")
	if err == storage.ErrTableDoesntExist {
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
	if err == storage.ErrTableDoesntExist {
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
	if err == storage.ErrTableDoesntExist {
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
	if err == storage.ErrTableDoesntExist {
		if _, err := s.DB.Exec(createBinariesTableQuery); err != nil {
			return err
		}
		if _, err := s.DB.Exec(createUniqueBinaryConstraint); err != nil {
			return err
		}
	}

	return nil
}

func (s *SQLdb) SetCard(username string, cardData storage.EncryptedData) error {

	var cardname string
	err := s.DB.QueryRow(checkCardNameQuery, cardData.Name, username).Scan(&cardname)
	if err != sql.ErrNoRows {
		return storage.ErrMetanameIsTaken
	}

	_, err = s.DB.Exec(setCardQuery, cardData.Name, cardData.Data, username)
	if err != nil {
		return fmt.Errorf("error setting card data in SetCard:%w", err)
	}
	return nil
}

func (s *SQLdb) GetCard(username string, cardname string) (storage.EncryptedData, error) {
	var cardData storage.EncryptedData
	err := s.DB.QueryRow(getCardQuery, cardname, username).Scan(&cardData.Name, &cardData.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return storage.EncryptedData{}, storage.ErrDataNotFound
		} else {
			return storage.EncryptedData{}, fmt.Errorf("error in GetCard:%w", err)
		}
	}
	return cardData, nil
}

func (s *SQLdb) ListCards(username string) (cardnames []string, err error) {

	rows, err := s.DB.Query(listCardsQuery, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrDataNotFound
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

func (s *SQLdb) SetLoginCred(username string, loginData storage.EncryptedData) error {
	_, err := s.DB.Exec(setLoginCredsQuery, loginData.Name, loginData.Data, username)
	if err != nil {
		return fmt.Errorf("error setting data in SetLoginCred:%w", err)
	}
	return nil
}

func (s *SQLdb) GetLoginCred(username string, loginname string) (logincred storage.EncryptedData, err error) {
	var encrData storage.EncryptedData
	err = s.DB.QueryRow(getLoginCredsQuery, loginname, username).Scan(&encrData.Name, &encrData.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return storage.EncryptedData{}, storage.ErrDataNotFound
		} else {
			return storage.EncryptedData{}, fmt.Errorf("error in GetLoginCred:%w", err)
		}
	}
	return encrData, nil
}

func (s *SQLdb) ListLoginCreds(username string) (logincrednames []string, err error) {

	rows, err := s.DB.Query(listLoginCredsQuery, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrDataNotFound
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

func (s *SQLdb) SetNote(username string, data storage.EncryptedData) error {
	_, err := s.DB.Exec(setNoteQuery, data.Name, data.Data, username)
	if err != nil {
		return fmt.Errorf("error setting data in SetLoginCred:%w", err)
	}
	return nil
}

func (s *SQLdb) GetNote(username string, notename string) (encrData storage.EncryptedData, err error) {

	err = s.DB.QueryRow(getNoteQuery, notename, username).Scan(&encrData.Name, &encrData.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return storage.EncryptedData{}, storage.ErrDataNotFound
		} else {
			return storage.EncryptedData{}, fmt.Errorf("error in GetLoginCred:%w", err)
		}
	}
	return encrData, nil
}

func (s *SQLdb) ListNotes(username string) (notenames []string, err error) {

	rows, err := s.DB.Query(listNotesQuery, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrDataNotFound
		}
		return nil, fmt.Errorf("couldn't ask database in ListLoginCreds:%w", err)
	}

	for rows.Next() {
		var itemname string

		err := rows.Scan(&itemname)
		if err != nil {
			return nil, fmt.Errorf("error scanning in ListLoginCreds:%w", err)
		}
		notenames = append(notenames, itemname)

	}

	err = rows.Err()

	if err != nil {
		return nil, fmt.Errorf("error scanning with rows.Next() in ListLoginCreds:%w", err)
	}

	return notenames, nil
}

func (s *SQLdb) SetBinary(username string, binary storage.Binary) error {
	_, err := s.DB.Exec(setBinaryQuery, binary.Name, binary.Data, username)
	if err != nil {
		return fmt.Errorf("error setting data in SetLoginCred:%w", err)
	}
	return nil
}

func (s *SQLdb) GetBinary(username string, binaryname string) (binary storage.Binary, err error) {

	err = s.DB.QueryRow(getBinaryQuery, binaryname, username).Scan(&binary.Name, &binary.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return storage.Binary{}, storage.ErrDataNotFound
		} else {
			return storage.Binary{}, fmt.Errorf("error in GetLoginCred:%w", err)
		}
	}
	return binary, nil
}

func (s *SQLdb) ListBinaries(username string) (binarynames []string, err error) {

	rows, err := s.DB.Query(listBinariesQuery, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrDataNotFound
		}
		return nil, fmt.Errorf("couldn't ask database in ListBinaries:%w", err)
	}

	for rows.Next() {
		var itemname string

		err := rows.Scan(&itemname)
		if err != nil {
			return nil, fmt.Errorf("error scanning in ListBinaries:%w", err)
		}
		binarynames = append(binarynames, itemname)

	}

	err = rows.Err()

	if err != nil {
		return nil, fmt.Errorf("error scanning with rows.Next() in ListBinaries:%w", err)
	}

	return binarynames, nil
}

func IsUniqueConstraintViolation(err error) bool {
	if pgErr, ok := err.(*pq.Error); ok {
		return pgErr.Code == "23505"
	}
	return false
}
