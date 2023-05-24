package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// writes the value v as json to the stream w
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
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

type APIError struct {
	Error string
}

type APIServer struct {
	listenAddr string
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	// TO DO reroute it to account
	// router.HandleFunc("/", makeHTTPHandlerFunc(s.handleAccount))

	router.HandleFunc("/account", makeHTTPHandlerFunc(s.handleAccount))

	router.HandleFunc("/account/{uuid}", makeHTTPHandlerFunc(s.handleGetAccount))

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{listenAddr: listenAddr}
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch method := r.Method; method {
	case "GET":
		return s.handleGetAccount(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	default:
		return fmt.Errorf("method not allowed: %s", method)
	}
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(mux.Vars(r)["uuid"])
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// TO DO check for a user in a database

	return WriteJSON(w, http.StatusOK, &Account{
		ID:       id,
		Username: "dummy",
		Email:    "user",
	})
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
