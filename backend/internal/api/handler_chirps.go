package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	httputil "github.com/Joshua-Lucas/go-chirpy/internal"
	"github.com/Joshua-Lucas/go-chirpy/internal/auth"
	"github.com/Joshua-Lucas/go-chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *APIConfig) CreateChirpHandler(w http.ResponseWriter, r *http.Request) {

	type reqBody struct {
		Body string `json:"body"`
	}

	// Valid user
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		httputil.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.TokenSecret)

	if err != nil {
		httputil.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	b := reqBody{}
	err = decoder.Decode(&b)

	if err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	cleanedChirp, err := validateChirp(b.Body)

	if err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	params := database.CreateChirpParams{
		Body:   cleanedChirp,
		UserID: userId,
	}

	newChirp, err := cfg.DBQueries.CreateChirp(r.Context(), params)

	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "Something went wrong creating chirp")
		return
	}

	wBody := Chirp{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserId:    newChirp.UserID,
	}

	err = httputil.RespondWithJSON(w, http.StatusCreated, wBody)

	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "Something went wrong responding ")
		return
	}
}

// validateChirp checks that a chirp body is within the allowed length and
// replaces any banned words with "****". It returns the cleaned body or an
// error if the chirp is too long.
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

func (cfg *APIConfig) GetAllChripsHandler(w http.ResponseWriter, r *http.Request) {

	filter := r.URL.Query().Get("author_id")
	sortQuery := r.URL.Query().Get("sort")

	var dbChirps []database.Chirp

	if filter != "" {
		filter, err := uuid.Parse(filter)

		if err != nil {

			httputil.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		dbChirps, err = cfg.DBQueries.ListChirpsByAuthor(r.Context(), filter)

		if err != nil {
			httputil.RespondWithError(w, http.StatusInternalServerError, "Something went wrong creating chirp")
			return
		}

	} else {
		var err error

		dbChirps, err = cfg.DBQueries.ListChirps(r.Context())

		if err != nil {
			httputil.RespondWithError(w, http.StatusInternalServerError, "Something went wrong creating chirp")
			return
		}
	}

	if sortQuery == "desc" {
		sort.Slice(dbChirps, func(i, j int) bool { return dbChirps[i].CreatedAt.After(dbChirps[j].CreatedAt) })
	}

	chirpList := make([]Chirp, 0, len(dbChirps))

	for _, v := range dbChirps {
		chirp := Chirp{

			ID:        v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Body:      v.Body,
			UserId:    v.UserID,
		}
		chirpList = append(chirpList, chirp)
	}

	err := httputil.RespondWithJSON(w, http.StatusOK, chirpList)

	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "Something went wrong responding ")
		return
	}
}

func (cfg *APIConfig) GetChripHandler(w http.ResponseWriter, r *http.Request) {

	chirpId, err := uuid.Parse(r.PathValue("chirpId"))

	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "Something went wrong responding ")
		return
	}

	dbChirp, err := cfg.DBQueries.GetChirp(r.Context(), chirpId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	}

	err = httputil.RespondWithJSON(w, http.StatusOK, chirp)

	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "Something went wrong responding ")
		return
	}
}

func (cfg *APIConfig) DeleteChripHandler(w http.ResponseWriter, r *http.Request) {

	// Valid user
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		httputil.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.TokenSecret)

	if err != nil {
		httputil.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	chirpId, err := uuid.Parse(r.PathValue("chirpId"))

	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "Something went wrong responding ")
		return
	}

	dbChirp, err := cfg.DBQueries.GetChirp(r.Context(), chirpId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if dbChirp.UserID != userId {
		httputil.RespondWithError(w, http.StatusForbidden, "User is not authorized for this action")
		return
	}

	err = cfg.DBQueries.DeleteChirp(r.Context(), chirpId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
