package storage

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Man-Crest/Go-Bank-Api/models"
	_ "github.com/lib/pq"
)

type Storaage interface {
	CreateAccountTable() error
	CreateAccount(a *models.Account) (*models.Account, error)
	DeleteAccount(int) error
	UpdateAccount(*models.Account) (*models.Account, error)
	GetAccountByID(int) error
	GetAccounts() ([]*models.Account, error)
}

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresConnection() (*PostgresStore, error) {
	connectionString := "user=Man_user dbname=go-bank password=password sslmode=disable"

	db, err := sql.Open("postgres", connectionString)
	// defer os.Exit(1)

	if err != nil {
		fmt.Errorf("Error ouccers in db conection :%w", err)
	}

	if err = db.Ping(); err != nil {
		fmt.Errorf("DB ping failed :%w", err)
	}

	return &PostgresStore{
		DB: db,
	}, nil
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
		id INTEGER PRIMARY KEY,
		firstname TEXT NOT NULL,
		lastname TEXT NOT NULL,
		account_number INTEGER NOT NULL UNIQUE, 
		balance REAL NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	`
	_, err := s.DB.Exec(query)
	if err != nil {
		log.Fatal(err, "error at create account")
		return err
	}
	return nil
}

func NewAccountFunc(firstname string, lastname string) *models.Account {
	return &models.Account{
		ID:        rand.Intn(1000),
		FirstName: firstname,
		LastName:  lastname,
		AccNumber: rand.Intn(1000000000),
		Balance:   0,
		CreatedAt: time.Now(),
	}
}

func (s *PostgresStore) CreateAccount(acc *models.Account) (*models.Account, error) {
	query := "INSERT INTO accounts (id,firstname,lastname,account_number,balance) VALUES ($1,$2,$3,$4,$5)"

	resp, err := s.DB.Exec(query, acc.ID, acc.FirstName, acc.LastName, acc.AccNumber, acc.Balance)
	if err != nil {
		fmt.Errorf("error ouccers in create newaccount query :%v", err)
		return nil, err
	}
	fmt.Println(resp)
	return acc, nil
}
func (s *PostgresStore) GetAccountByID(id int) (*models.Account, error) {
	query := "SELECT FROM * accounts WHERE id=$1"

	resp := s.DB.QueryRow(query, id)
	// if err != nil {
	// 	fmt.Errorf("error ouccers in create newaccount query :%v", err)
	// 	return nil, err
	// }

	respAccount := &models.Account{}
	resp.Scan(
		&respAccount.ID,
		&respAccount.FirstName,
		&respAccount.LastName,
		&respAccount.AccNumber,
		&respAccount.Balance,
		&respAccount.CreatedAt,
	)

	fmt.Println(respAccount)
	return respAccount, nil
}
func (s *PostgresStore) DeleteAccount(id int) error {
	query := "DELETE FROM accounts WHERE id=$1"

	_ = s.DB.QueryRow(query, id)

	return nil
}
func (s *PostgresStore) UpdateAccount(acc *models.Account) (*models.Account, error) {
	query := `UPDATE accounts
	SET firstname = $1, lastname = $2,
	WHERE id = $3;`

	resp := s.DB.QueryRow(query, acc.FirstName, acc.LastName, acc.ID)

	respAccount := &models.Account{}
	resp.Scan(
		&respAccount.ID,
		&respAccount.FirstName,
		&respAccount.LastName,
		&respAccount.AccNumber,
		&respAccount.Balance,
		&respAccount.CreatedAt,
	)
	return respAccount, nil
}

func (s *PostgresStore) GetAccounts() ([]*models.Account, error) {
	query := "SELECT * FROM accounts"

	rows, err := s.DB.Query(query)
	if err != nil {
		fmt.Errorf("error ouccers in create newaccount query :%v", err)
		return nil, err
	}

	accounts := []*models.Account{}

	for rows.Next() {

		Account := &models.Account{}
		rows.Scan(
			&Account.ID,
			&Account.FirstName,
			&Account.LastName,
			&Account.AccNumber,
			&Account.Balance,
			&Account.CreatedAt,
		)
		accounts = append(accounts, Account)
	}
	return accounts, nil
}
