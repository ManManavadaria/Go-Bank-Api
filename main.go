package main

import (
	"log"

	"github.com/Man-Crest/Go-Bank-Api/storage"
)

func main() {
	PostgresStore, err := storage.NewPostgresConnection()
	if err != nil {
		log.Fatal(err)
	}

	err = PostgresStore.CreateAccountTable()
	if err != nil {
		log.Fatal(err)
	}

	Server := NewServer(":7000", PostgresStore)
	Server.run()
}
