package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeJson(w http.ResponseWriter, status int, v any) error {
	enc := json.NewEncoder(w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if v != nil {
		if err := enc.Encode(v); err != nil {
			fmt.Println(err.Error())
			fmt.Println("error when writing json in writeJson()")
			return err
		}
	}

	return nil
}

func printError(err error, message string) {
	fmt.Println(err.Error())
	fmt.Println(message)
}
