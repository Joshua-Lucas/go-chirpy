package api

//import (
//	"encoding/json"
//	"net/http"
//	"time"

//	httputil "github.com/Joshua-Lucas/go-chirpy/internal"
//	"github.com/google/uuid"
//)

//func (cfg *APIConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
//	type body struct {
//		Email string `json:"email"`
//	}

//	type User struct {
//		ID        uuid.UUID `json:"id"`
//		CreatedAt time.Time `json:"created_at"`
//		UpdatedAt time.Time `json:"updated_at"`
//		Email     string    `json:"email"`
//	}

//	decoder := json.NewDecoder(r.Body)
//	b := body{}
//	err := decoder.Decode(&b)

//	if err != nil {
//		httputil.RespondWithError(w, http.StatusBadGateway, "Something went wrong")
//	}

//	//TODO: Continue project
//}
