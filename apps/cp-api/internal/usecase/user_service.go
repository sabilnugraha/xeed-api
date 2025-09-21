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
	signer contract.TokenSigner // ← baru (untuk JWT)
}

var _ contract.UserService = (*userService)(nil)

func NewUserService(
	repo contract.UserRepository,
	clk contract.Clock,
	idg contract.IDGen,
	hasher contract.PasswordHasher,
	signer contract.TokenSigner, // ← baru
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
	if signer == nil {
		panic("NewUserService: signer is nil")
	}
	return &userService{repo: repo, clock: clk, idgen: idg, hasher: hasher, signer: signer}
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

var ErrInvalidCredential = errors.New("invalid email or password")

func (s *userService) Login(ctx context.Context, in dto.LoginRequest) (*dto.LoginResponse, error) {
	if s.signer == nil {
		return nil, errors.New("internal: token signer not configured")
	}

	email := strings.ToLower(strings.TrimSpace(in.Email))
	if email == "" || !strings.Contains(email, "@") {
		return nil, ErrInvalidCredential
	}
	if len(in.Password) == 0 {
		return nil, ErrInvalidCredential
	}

	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if u == nil || u.PasswordHash == nil {
		return nil, ErrInvalidCredential
	}

	// verify password (bcrypt)
	if !s.hasher.Verify(in.Password, *u.PasswordHash) {
		return nil, ErrInvalidCredential
	}

	// optional: cek status BLOCKED, dsb
	// if u.Status == domain.UserBlocked { return nil, errors.New("user blocked") }

	now := s.clock.Now()
	tok, err := s.signer.Sign(u.UserID, u.Email, now)
	if err != nil {
		return nil, err
	}

	resp := dto.LoginResponse{
		AccessToken: tok,
		User: dto.UserResponse{
			UserID:      u.UserID,
			Email:       u.Email,
			DisplayName: u.DisplayName,
			PhoneE164:   u.PhoneE164,
			Locale:      u.Locale,
			Timezone:    u.Timezone,
			Status:      string(u.Status),
		},
	}
	return &resp, nil
}
