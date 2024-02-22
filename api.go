package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Man-Crest/Go-Bank-Api/models"
	"github.com/Man-Crest/Go-Bank-Api/storage"
	"github.com/gorilla/mux"
)

type ApiServer struct {
	listenAddress string
	DB            *storage.PostgresStore
}

func NewServer(listenAddress string, db *storage.PostgresStore) *ApiServer {
	return &ApiServer{
		listenAddress: listenAddress,
		DB:            db,
	}
}

func (s *ApiServer) run() {
	router := mux.NewRouter()

	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccounts))

	fmt.Printf("server is running on port number: %s", s.listenAddress)
	http.ListenAndServe(s.listenAddress, router)
}

func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	// if r.Method == "POST" {
	// 	return s.handleCreateAccount(w, r)
	// }
	if r.Method == "PUT" {
		return s.handleTransfer(w, r)
	}
	// if r.Method == "DELETE" {
	// 	return s.handleDeleteAccount(w, r)
	// }
	return nil
}
func (s *ApiServer) handleAccounts(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccounts(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	// if r.Method == "PUT" {
	// 	return s.handleTransfer(w, r)
	// }
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return nil
}

func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	data := models.TempAcc()
	WriteJSON(w, http.StatusOK, data)
	return nil
}
func (s *ApiServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {

	accounts, err := s.DB.GetAccounts()
	if err != nil {
		fmt.Errorf(err.Error(), "error getting all accounts")
	}
	fmt.Println(accounts)
	WriteJSON(w, http.StatusOK, accounts)
	return nil
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	NewAccReq := struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&NewAccReq)
	if err != nil {
		fmt.Errorf("error decodeing values from body")
	}

	account := storage.NewAccountFunc(NewAccReq.FirstName, NewAccReq.LastName)

	acc, _ := s.DB.CreateAccount(account)

	if err != nil {
		return err
		log.Fatal(err)
	}
	WriteJSON(w, http.StatusOK, acc)
	return nil
}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	DeleteAccReqID := mux.Vars(r)

	fmt.Printf("id  %v", DeleteAccReqID["id"])

	DeleteAccReqID2, err := strconv.Atoi(DeleteAccReqID["id"])

	if err != nil {
		log.Fatal(err)
	}

	err = s.DB.DeleteAccount(DeleteAccReqID2)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func makeHTTPHandleFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := f(w, r)

		if err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)

	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(v)
}

type ApiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}
