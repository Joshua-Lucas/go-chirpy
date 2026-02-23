package api

import (
	"encoding/json"
	"net/http"
	"time"

	httputil "github.com/Joshua-Lucas/go-chirpy/internal"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *APIConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	b := body{}
	err := decoder.Decode(&b)

	if err != nil {
		httputil.RespondWithError(w, http.StatusBadGateway, "Something went wrong")
	}

	newUser, err := cfg.DBQueries.CreateUser(r.Context(), b.Email)

	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "Something went wrong creating user")
	}

	wBody := User{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	}

	// Respond with JSON
	err = httputil.RespondWithJSON(w, http.StatusCreated, wBody)

	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "Something went wrong responding ")
	}

}
