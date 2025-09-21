package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"xeed/apps/cp-api/internal/dto"
	"xeed/apps/cp-api/internal/usecase"
	"xeed/apps/cp-api/internal/usecase/contract"

	"github.com/google/uuid"
)

type UserHandler struct {
	svc contract.UserService
}

func NewUserHandler(svc contract.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Pastikan CreatedBy bernilai nil kalau kosong (bukan Nil UUID)
	if req.CreatedBy != nil && *req.CreatedBy == uuid.Nil {
		req.CreatedBy = nil
	}

	user, err := h.svc.RegisterUser(r.Context(), req) // langsung pass DTO ke service
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := dto.UserResponse{
		UserID:   user.UserID,
		Email:    user.Email,
		Status:   string(user.Status),
		Locale:   user.Locale,
		Timezone: user.Timezone,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := h.svc.Login(r.Context(), req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, usecase.ErrInvalidCredential) {
			status = http.StatusUnauthorized
		}
		http.Error(w, err.Error(), status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
