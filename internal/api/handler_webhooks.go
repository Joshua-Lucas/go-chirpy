package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	httputil "github.com/Joshua-Lucas/go-chirpy/internal"
	"github.com/google/uuid"
)

func (cfg *APIConfig) PolkaHandler(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	b := body{}
	err := decoder.Decode(&b)

	if err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, "Something went wrong decoding body")
		return
	}

	if b.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	parsedId, err := uuid.Parse(b.Data.UserID)

	if err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = cfg.DBQueries.UpgradeToRed(r.Context(), parsedId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		httputil.RespondWithError(w, http.StatusInternalServerError, "couldn't upgrade user")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
