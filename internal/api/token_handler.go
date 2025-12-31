package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Anezz12/femProject/internal/store"
	"github.com/Anezz12/femProject/internal/utils"
)

type ToekenHandler struct {
	tokenSotre store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *ToekenHandler {
	return &ToekenHandler{
		tokenSotre: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *ToekenHandler) HandlerCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Println("ERROR: decodeCreateTokenRequest:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.
			Envelope{"error": "Invalid request payload"})
		return
	}

	// let get user by username
	user, err := h.userStore.GetUserByName(req.Username)
	if err != nil || user == nil {
		h.logger.Println("ERROR: getUserByName:", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Internal server error"})
		return
	}

	passwordDoesMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		h.logger.Println("ERROR: passwordMatches:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	if !passwordDoesMatch {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid credentials"})
		return
	}

	token, err := h.tokenSotre.CreateToken(user.ID, 24*time.Hour, "authentication")
	if err != nil {
		h.logger.Println("ERROR: createToken:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"auth_token": token.Plaintext})

}
