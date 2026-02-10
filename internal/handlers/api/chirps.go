package api

import (
	"encoding/json"
	"net/http"
	"strings"

	httputil "github.com/Joshua-Lucas/go-chirpy/internal"
)

func ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	b := body{}
	err := decoder.Decode(&b)

	if err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	max_characters := 140

	// Validate length of string
	if len(b.Body) > max_characters {
		httputil.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// Validate profane words
	var profaneWords = map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Fields(b.Body)
	for i := range words {
		if _, ok := profaneWords[strings.ToLower(words[i])]; ok {
			words[i] = "****"
		}
	}

	parsedWords := strings.Join(words, " ")

	respBody := returnVals{
		CleanedBody: parsedWords,
	}

	httputil.RespondWithJSON(w, http.StatusOK, respBody)

}
