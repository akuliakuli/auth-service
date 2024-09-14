package handlers

import (
	"github.com/akuliakuli/auth-service/internal/services"
	"github.com/akuliakuli/auth-service/internal/db"
	"encoding/json"
	"net/http"
)

type TokenRequest struct {
	UserID string `json:"userId"`
	IP     string `json:"ip"`
}

type RefreshRequest struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	IP           string `json:"ip"`
}

func GenerateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req TokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	accessToken, err := services.GenerateAccessToken(req.UserID, req.IP)
	if err != nil {
		http.Error(w, "Could not generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, hashedRefreshToken, err := services.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "Could not generate refresh token", http.StatusInternalServerError)
		return
	}

	db.StoreRefreshToken(req.UserID, hashedRefreshToken)

	json.NewEncoder(w).Encode(map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	claims, err := services.ValidateAccessToken(req.AccessToken)
	if err != nil {
		http.Error(w, "Invalid access token", http.StatusUnauthorized)
		return
	}

	valid, err := db.ValidateRefreshToken(claims.UserID, req.RefreshToken)
	if !valid || err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	if claims.IP != req.IP {
		services.SendEmailWarning(claims.UserID, claims.IP, req.IP)
	}

	newAccessToken, err := services.GenerateAccessToken(claims.UserID, req.IP)
	if err != nil {
		http.Error(w, "Could not generate access token", http.StatusInternalServerError)
		return
	}

	newRefreshToken, hashedNewRefreshToken, err := services.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "Could not generate refresh token", http.StatusInternalServerError)
		return
	}

	db.StoreRefreshToken(claims.UserID, hashedNewRefreshToken)

	json.NewEncoder(w).Encode(map[string]string{
		"accessToken":  newAccessToken,
		"refreshToken": newRefreshToken,
	})
}
