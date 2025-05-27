package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	// Define a struct to hold the parameters from the request body
	// This struct will be used to decode the JSON request body
	type parameters struct {
		Body string `json:"body"`
	}

	// Define a struct to hold the return values
	// This struct will be used to encode the response body
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	// Decode the JSON request body into the ChirpRequest struct
	decoder := json.NewDecoder(r.Body)
	params := parameters{}         // Initialize the parameters struct
	err := decoder.Decode(&params) // Decode the request body into params
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	// Define a map of bad words to check against
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	// Check for bad words in the chirp body
	cleaned := cleanChirp(params.Body, badWords)
	// Here you would typically validate the chirp content.
	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleaned,
	})
}

// cleanChirp replaces bad words in the chirp body with "****".
// It takes the body of the chirp and a map of bad words to check against.
// It returns the cleaned chirp body.
func cleanChirp(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
