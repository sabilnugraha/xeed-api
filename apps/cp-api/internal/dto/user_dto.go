package dto

import "github.com/google/uuid"

// DTO murni untuk transport layer (HTTP JSON, gRPC, dsb).
// Tidak bawa logic bisnis, hanya data binding.

type RegisterUserRequest struct {
	Email         string     `json:"email"`
	Password      string     `json:"password"`
	DisplayName   *string    `json:"displayName,omitempty"`
	PhoneE164     *string    `json:"phoneE164,omitempty"`
	Locale        string     `json:"locale,omitempty"`
	Timezone      string     `json:"timezone,omitempty"`
	IsServiceAcct bool       `json:"isServiceAccount,omitempty"`
	CreatedBy     *uuid.UUID `json:"createdBy,omitempty"`
}

type UserResponse struct {
	UserID      uuid.UUID `json:"userId"`
	DisplayName *string   `json:"displayName,omitempty"`
	PhoneE164   *string   `json:"phoneE164,omitempty"`
	Email       string    `json:"email"`
	Status      string    `json:"status"`
	Locale      string    `json:"locale"`
	Timezone    string    `json:"timezone"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string       `json:"accessToken"`
	User        UserResponse `json:"user"`
}
