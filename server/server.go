// has all the endpoints
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

type TokenVerifyRequest struct { 
	Token string
}

type MyCustomClaims struct {
	Username string `json:"username"`
	Uid uuid.UUID `json:"uid"`
	jwt.RegisteredClaims
}

func (s *Server) init() {
	s.router.Get("/api/health", s.handleHealth)
	s.router.Post("/api/login", s.handleLogin)
	s.router.Post("/api/signup", s.handleSignup)
	s.router.Post("/api/token/verify", s.handleTokenVerify)
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

func (s *Server) handleTokenVerify(w http.ResponseWriter, r *http.Request) {
	body := new(TokenVerifyRequest)
	json.NewDecoder(r.Body).Decode(body)

	tokenString := body.Token

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		fmt.Printf("%v %v\n", claims.Username, claims.RegisteredClaims.Issuer)
		res := struct {
			Message string `json:"message"`
		}{
			Message: "token valid",
		}
		writeJson(w, http.StatusOK, res)
	} else {
		fmt.Println(err)
		res := struct {
			Message string `json:"message"`
		}{
			Message: "token invalid",
		}
		writeJson(w, http.StatusUnauthorized, res)
	}
	
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

	mySigningKey := []byte(os.Getenv("JWT_SECRET"))
	fmt.Println(mySigningKey)
	
	// Create claims with multiple fields populated
	claims := MyCustomClaims{
		user.Username,
		user.Id,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer: "trackitserver",
		},
	}
	
	fmt.Printf("username: %v\n", claims.Username)
	fmt.Printf("uid: %v\n", claims.Uid)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println("error when creating jwt")
		res := struct {
			Message string `json:"message"`
		}{
			Message: "internal_server_error",
		}
	
		writeJson(w, http.StatusInternalServerError, res)
	}

	res := struct {
		Token string `json:"token"`
	}{
		Token: ss,
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

	mySigningKey := []byte(os.Getenv("JWT_SECRET"))
	fmt.Println(mySigningKey)

	type MyCustomClaims struct {
		Username string `json:"username"`
		Uid uuid.UUID `json:"uid"`
		jwt.RegisteredClaims
	}
	
	// Create claims with multiple fields populated
	claims := MyCustomClaims{
		user.Username,
		user.Id,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer: "trackitserver",
		},
	}
	
	fmt.Printf("username: %v\n", claims.Username)
	fmt.Printf("uid: %v\n", claims.Uid)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println("error when creating jwt")
		res := struct {
			Message string `json:"message"`
		}{
			Message: "internal_server_error",
		}
	
		writeJson(w, http.StatusInternalServerError, res)
	}

	res := struct {
		Token string `json:"token"`
	}{
		Token: ss,
	}

	writeJson(w, http.StatusOK, res)
}