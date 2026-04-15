package api

import (
	"net/http"
	"time"

	httputil "github.com/Joshua-Lucas/go-chirpy/internal"
	"github.com/Joshua-Lucas/go-chirpy/internal/auth"
)

type refreshTokenResponse struct {
	Token string `json:"token"`
}

func (cfg *APIConfig) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	// Valid user
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		httputil.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	refreshTokenData, err := cfg.DBQueries.GetUserFromRefreshToken(r.Context(), token)

	if err != nil {
		httputil.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if refreshTokenData.ExpiresAt.Before(time.Now()) || refreshTokenData.RevokedAt.Valid {
		httputil.RespondWithError(w, http.StatusUnauthorized, "Token is expired or revoked")
		return
	}

	jwtToken, err := auth.MakeJWT(refreshTokenData.UserID, cfg.TokenSecret, time.Hour)

	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "Error occurred with token assignment")
		return
	}

	wBody := refreshTokenResponse{
		Token: jwtToken,
	}

	// Respond with JSON
	httputil.RespondWithJSON(w, http.StatusOK, wBody)

}
