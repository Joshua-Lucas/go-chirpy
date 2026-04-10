package api

import (
	"encoding/json"
	"net/http"
	"time"

	httputil "github.com/Joshua-Lucas/go-chirpy/internal"
	"github.com/Joshua-Lucas/go-chirpy/internal/auth"
	"github.com/Joshua-Lucas/go-chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *APIConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	b := body{}
	err := decoder.Decode(&b)

	if err != nil {
		httputil.RespondWithError(w, http.StatusBadGateway, "Something went wrong decoding body")
	}

	hashedPassword, err := auth.HashPassword(b.Password)

	if err != nil {
		httputil.RespondWithError(w, http.StatusBadGateway, "Something went wrong hashing password")
	}

	userParams := database.CreateUserParams{
		Email:          b.Email,
		HashedPassword: hashedPassword,
	}

	newUser, err := cfg.DBQueries.CreateUser(r.Context(), userParams)

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

func (cfg *APIConfig) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	b := body{}
	err := decoder.Decode(&b)

	if err != nil {
		httputil.RespondWithError(w, http.StatusBadGateway, "Something went wrong decoding body")
	}

	user, err := cfg.DBQueries.GetUserByEmail(r.Context(), b.Email)

	if err != nil {
		httputil.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
	}

	isValid, err := auth.CheckPasswordHash(b.Password, user.HashedPassword)

	if err != nil {
		httputil.RespondWithError(w, http.StatusBadGateway, "Something went wrong")
	}

	if isValid == false {
		httputil.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
	}

	// --- Token assignment ---

	// Access Token
	expiresIn := time.Hour

	token, err := auth.MakeJWT(user.ID, cfg.TokenSecret, expiresIn)

	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "Error occurred with token assignment")
	}
	// Refresh token
	rawRefreshToken := auth.MakeRefreshToken()

	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     rawRefreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	}
	refreshToken, err := cfg.DBQueries.CreateRefreshToken(r.Context(), refreshTokenParams)

	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "Error occurred with refresh token assignment")
	}

	wBody := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}

	// Respond with JSON
	err = httputil.RespondWithJSON(w, http.StatusOK, wBody)

	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "Something went wrong responding ")
	}
}
