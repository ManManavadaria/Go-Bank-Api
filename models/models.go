package models

import "time"

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	AccNumber int       `json:"account_number"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"create_at"`
}

func TempAcc() *Account {
	return &Account{
		ID:        1,
		FirstName: "asfsf",
		LastName:  "dfgsdf",
		AccNumber: 121211212,
		Balance:   123123123,
		CreatedAt: time.Now().UTC(),
	}
}
