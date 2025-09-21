package contract

import (
	"context"
	"time"

	"xeed/apps/cp-api/internal/domain"
	"xeed/apps/cp-api/internal/dto"

	"github.com/google/uuid"
)

// Repository interface yang harus diimplementasikan infra (pg, mongo, dll)
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error) // nil,nil kalau tidak ada
	Create(ctx context.Context, u domain.User) (*domain.User, error)
}

// Service interface untuk layer bisnis
type UserService interface {
	RegisterUser(ctx context.Context, in dto.RegisterUserRequest) (*domain.User, error)
}

// Adapter utilitas (Clock, UUID, PasswordHasher)
type Clock interface{ Now() time.Time }
type IDGen interface{ New() uuid.UUID }
type PasswordHasher interface {
	Hash(plain string) (hash string, alg domain.PasswordAlg, updatedAt time.Time, err error)
}
