// has all the endpoints
package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	database *Database
	router   *chi.Mux
}

type UserPwdRequest struct {
	Username string
	Password string
}

func (s *Server) init() {
	s.router.Get("/api/health", s.handleHealth)
	s.router.Post("/api/login", s.handleLogin)
	s.router.Post("/api/signup", s.handleSignup)
	http.ListenAndServe(":8080", s.router)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Message string `json:"message"`
	}{
		Message: "health seems fine",
	}

	writeJson(w, http.StatusOK, body)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	body := new(UserPwdRequest)
	json.NewDecoder(r.Body).Decode(body)

	user, err := s.database.getUser(body.Username)
	if err != nil {
		res := struct {
			Message string `json:"message"`
		}{
			Message: "couldn't get user from database",
		}
		writeJson(w, http.StatusUnauthorized, res)
		return
	}

	// compare password with user hash
	if err := bcrypt.CompareHashAndPassword(user.Hash, []byte(body.Password)); err != nil {
		res := struct {
			Message string `json:"message"`
		}{
			Message: "wrong password",
		}
		writeJson(w, http.StatusUnauthorized, res)
		return
	}

	res := struct {
		Message string `json:"message"`
	}{
		Message: "logged in",
	}

	writeJson(w, http.StatusOK, res)
}

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
	body := new(UserPwdRequest)
	json.NewDecoder(r.Body).Decode(body)

	_, err := s.database.getUser(body.Username)
	if err == nil {
		res := struct {
			Message string `json:"message"`
		}{
			Message: "user already exists",
		}
		writeJson(w, http.StatusUnauthorized, res)
		return
	}

	if err := s.database.createUser(body.Username, body.Password); err != nil {
		res := struct {
			Message string `json:"message"`
		}{
			Message: "failed to create user",
		}
		writeJson(w, http.StatusInternalServerError, res)
		return
	}

	writeJson(w, http.StatusCreated, nil)
}