package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type APIError struct {
	Error string `json:"error"`
}

type APIServer struct {
	listenAddr string
	store      Storage
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	// TO DO reroute it to account
	// router.HandleFunc("/", makeHTTPHandlerFunc(s.handleAccount))

	router.HandleFunc("/account/", makeHTTPHandlerFunc(s.handleAccount))

	router.HandleFunc("/account/{uuid}/", withJWTAuth(makeHTTPHandlerFunc(s.handleGetAccountByID)))

	router.HandleFunc("/health/", makeHTTPHandlerFunc(s.handleHealth))

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return fmt.Errorf("method not allowed: %s", r.Method)
	}

	return WriteJSON(w, http.StatusOK, "Everything seems fine!")
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch method := r.Method; method {
	case "GET":
		return s.handleGetAccount(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	default:
		return fmt.Errorf("method not allowed: %s", method)
	}
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	switch method := r.Method; method {
	case "GET":
		id, err := getID(r)
		if err != nil {
			return err
		}
		account, err := s.store.GetAccountByID(id)

		if err != nil {
			str := fmt.Sprintf("Account with uuid %s not found", id)
			return WriteJSON(w, http.StatusNotFound, str)
		}

		return WriteJSON(w, http.StatusOK, account)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	default:
		return fmt.Errorf("method not allowed: %s", method)
	}
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	account := NewAccount(req.Username, req.Email)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	tokenString, err := createJWT(account)

	if err != nil {
		return err
	}

	fmt.Println(tokenString)

	body := struct {
		Account *Account
		Token   string
	}{
		Account: account,
		Token:   tokenString,
	}

	return WriteJSON(w, http.StatusCreated, body)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		log.Println(err.Error())
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]any{"deleted": id})
}

// writes the value v as json to the stream w
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET_KEY")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

}

func withJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling jwt auth middleware")

		tokenString := r.Header.Get("x-jwt-token")
		_, err := validateJWT(tokenString)

		if err != nil {
			WriteJSON(w, http.StatusForbidden, APIError{Error: "Invalid token"})
			return
		}

		handlerFunc(w, r)
	}
}

func createJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiersAt": 15000,
		"accountID": account.ID,
	}

	secret := os.Getenv("JWT_SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

// a type (which is the function signature for our handler functions)
// we need this because we return an error, but router.HandleFunc
// takes a function which doesn't return an error
type apiFunc func(http.ResponseWriter, *http.Request) error

// this takes our apiFunc, and returns a function that will call it
// and handle the error, but the signature will match http.HandlerFunc
func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// TO DO
			// later we can change our error and switch and write an appropriate error code
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (uuid.UUID, error) {
	id, err := uuid.Parse(mux.Vars(r)["uuid"])

	if err != nil {
		log.Println("error when parsing uuid", err.Error())
		return id, fmt.Errorf("invalid uuid : %s, %s", mux.Vars(r)["uuid"], err.Error())
	}

	return id, nil
}
