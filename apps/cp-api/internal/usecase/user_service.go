package usecase

import (
	"context"
	"errors"
	"strings"

	"xeed/apps/cp-api/internal/domain"
	"xeed/apps/cp-api/internal/dto"
	"xeed/apps/cp-api/internal/usecase/contract"
)

type userService struct {
	repo   contract.UserRepository
	clock  contract.Clock
	idgen  contract.IDGen
	hasher contract.PasswordHasher
}

var _ contract.UserService = (*userService)(nil)

func NewUserService(
	repo contract.UserRepository,
	clk contract.Clock,
	idg contract.IDGen,
	hasher contract.PasswordHasher,
) contract.UserService {
	if repo == nil {
		panic("NewUserService: repo is nil")
	}
	if clk == nil {
		panic("NewUserService: clock is nil")
	}
	if idg == nil {
		panic("NewUserService: idgen is nil")
	}
	if hasher == nil {
		panic("NewUserService: hasher is nil")
	}
	return &userService{repo: repo, clock: clk, idgen: idg, hasher: hasher}
}

func (s *userService) RegisterUser(ctx context.Context, in dto.RegisterUserRequest) (*domain.User, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	if email == "" || !strings.Contains(email, "@") {
		return nil, errors.New("invalid email")
	}
	if len(in.Password) < 8 {
		return nil, errors.New("password min 8 chars")
	}

	exist, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return nil, errors.New("email already registered")
	}

	hash, alg, pwdAt, err := s.hasher.Hash(in.Password)
	if err != nil {
		return nil, err
	}

	now := s.clock.Now()
	u := domain.User{
		UserID:             s.idgen.New(),
		Email:              email,
		DisplayName:        in.DisplayName,
		PhoneE164:          in.PhoneE164,
		Locale:             def(in.Locale, "en"),
		Timezone:           def(in.Timezone, "UTC"),
		Status:             domain.UserStatus("ACTIVE"),
		IsServiceAccount:   in.IsServiceAcct,
		PasswordAlg:        alg,
		PasswordHash:       &hash,
		PasswordUpdatedAt:  &pwdAt,
		MustChangePassword: false,
		MFAEnrolled:        false,
		CreatedAt:          now,
		UpdatedAt:          now,
		CreatedBy:          in.CreatedBy, // *uuid.UUID di DTO -> cocok dengan domain
		UpdatedBy:          in.CreatedBy,
	}

	return s.repo.Create(ctx, u)
}

func def(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}
