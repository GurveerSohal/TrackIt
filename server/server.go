// has all the endpoints
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type Server struct {
	database *Database
	router   *chi.Mux
}

func writeJson(w http.ResponseWriter, status int, v any) error {
	enc := json.NewEncoder(w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := enc.Encode(v); err != nil {
		fmt.Println(err.Error())
		fmt.Println("error when writing json in writeJson()")
		return err
	}

	return nil
}

func (s *Server) init() {
	s.router.Get("/api/health", handleHealth)

	http.ListenAndServe(":8080", s.router)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Message string `json:"message"`
	}{
		Message: "health seems fine",
	}

	writeJson(w, http.StatusOK, body)
}
