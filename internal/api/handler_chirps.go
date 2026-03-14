package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	httputil "github.com/Joshua-Lucas/go-chirpy/internal"
)

func (cfg *APIConfig) CreateChirp(w http.ResponseWriter, r http.Request) {

	type reqBody struct {
		Body   string `json:"body"`
		UserId string `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	b := reqBody{}
	err := decoder.Decode(&b)

	if err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}
	
	// TODO: Validate chirp string
	// TODO: Run Create query
	// TODO: Hanlde the return 
}

// validateChirp checks that a chirp body is within the allowed length
// and replaces any banned words with "****".
// It returns the cleaned body or an error if the chirp is too long.
func validateChirp(body string) (string, error) {

	max_characters := 140

	// Validate length of string
	if len(body) > max_characters {
		return "", errors.New("Chirp is too long")
	}

	// Validate profane words
	var profaneWords = map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Fields(body)
	for i := range words {
		if _, ok := profaneWords[strings.ToLower(words[i])]; ok {
			words[i] = "****"
		}
	}

	parsedWords := strings.Join(words, " ")

	return parsedWords, nil
}
