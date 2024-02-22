package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Man-Crest/Go-Bank-Api/storage"
	"github.com/golang-jwt/jwt/v5"
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

	router.HandleFunc("/account/{id}", WithJWTAuth(makeHTTPHandleFunc(s.handleAccountByID)))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccounts))

	fmt.Printf("server is running on port number: %s", s.listenAddress)
	http.ListenAndServe(s.listenAddress, router)
}

func WithJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("jwt request is been made")

		TokenString := r.Header.Get("x-jwt-token")

		_, err := ValidateJWT(TokenString)

		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "eorror in valiadting the jwt token"})
			return
		}

		handlerFunc(w, r)
	}
}
func ValidateJWT(TokenString string) (*jwt.Token, error) {
	secret := "fahubfhbliBdrv"
	return jwt.Parse(TokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func (s *ApiServer) handleAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccountbyID(w, r)
	}
	// if r.Method == "POST" {
	// 	return s.handleCreateAccount(w, r)
	// }
	if r.Method == "PUT" {
		return s.handleTransfer(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
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
	// if r.Method == "DELETE" {
	// 	return s.handleDeleteAccount(w, r)
	// }
	return nil
}

func (s *ApiServer) handleGetAccountbyID(w http.ResponseWriter, r *http.Request) error {

	id, err := getID(r)
	if err != nil {
		return err
	}

	data, err := s.DB.GetAccountByID(id)

	if err != nil {
		return err
	}
	WriteJSON(w, http.StatusOK, data)
	return nil
}
func (s *ApiServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {

	fmt.Println("inside all account get")
	accounts, err := s.DB.GetAccounts()
	if err != nil {
		return err
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
		return err
	}

	account := storage.NewAccountFunc(NewAccReq.FirstName, NewAccReq.LastName)

	acc, _ := s.DB.CreateAccount(account)

	if err != nil {
		return err
	}
	WriteJSON(w, http.StatusOK, acc)
	return nil
}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	id, err := getID(r)

	if err != nil {
		return err
	}

	err = s.DB.DeleteAccount(id)
	if err != nil {
		return err
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

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}
