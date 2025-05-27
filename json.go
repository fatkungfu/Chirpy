package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	// Define a struct for the error response
	type errorResponse struct {
		Error string `json:"error"`
	}

	// Set the response header and status code
	response := errorResponse{
		Error: msg,
	}
	respondWithJSON(w, code, response)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response, _ := json.Marshal(payload)
	w.Write(response)
}
