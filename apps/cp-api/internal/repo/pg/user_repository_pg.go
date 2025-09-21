// apps/cp-api/internal/repo/pg/user_repository_pg.go
package pg

import (
	"context"
	"errors"

	"xeed/apps/cp-api/internal/domain"
	"xeed/apps/cp-api/internal/usecase/contract"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepoPG struct {
	db *pgxpool.Pool
}

func NewUserRepositoryPG(db *pgxpool.Pool) contract.UserRepository {
	return &userRepoPG{db: db}
}

func (r *userRepoPG) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT
			"UserID","Email","EmailVerifiedAt","PhoneE164","PhoneVerifiedAt",
			"PasswordHash","PasswordAlg","PasswordUpdatedAt","MustChangePassword",
			"Status","IsServiceAccount","DisplayName","AvatarURL",
			"Locale","Timezone","Preferences","MFAEnrolled","MFADefaultMethod",
			"LastLoginAt","LastLoginIP","CreatedAt","CreatedBy",
			"UpdatedAt","UpdatedBy","IsDeleted"
		FROM "User"
		WHERE "Email" = $1 AND "IsDeleted" = FALSE
	`
	row := r.db.QueryRow(ctx, q, email)

	var ur UserRow
	if err := row.Scan(
		&ur.UserID, &ur.Email, &ur.EmailVerifiedAt, &ur.PhoneE164, &ur.PhoneVerifiedAt,
		&ur.PasswordHash, &ur.PasswordAlg, &ur.PasswordUpdatedAt, &ur.MustChangePassword,
		&ur.Status, &ur.IsServiceAccount, &ur.DisplayName, &ur.AvatarURL,
		&ur.Locale, &ur.Timezone, &ur.Preferences, &ur.MFAEnrolled, &ur.MFADefaultMethod,
		&ur.LastLoginAt, &ur.LastLoginIP, &ur.CreatedAt, &ur.CreatedBy,
		&ur.UpdatedAt, &ur.UpdatedBy, &ur.IsDeleted,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	u, err := ur.ToDomain()
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepoPG) Create(ctx context.Context, u domain.User) (*domain.User, error) {
	const q = `
		INSERT INTO "User" (
			"UserID","Email","EmailVerifiedAt","PhoneE164","PhoneVerifiedAt",
			"PasswordHash","PasswordAlg","PasswordUpdatedAt","MustChangePassword",
			"Status","IsServiceAccount","DisplayName","AvatarURL",
			"Locale","Timezone","Preferences","MFAEnrolled","MFADefaultMethod",
			"LastLoginAt","LastLoginIP","CreatedAt","CreatedBy",
			"UpdatedAt","UpdatedBy","IsDeleted"
		) VALUES (
			$1,$2,$3,$4,$5,
			$6,$7,$8,$9,
			$10,$11,$12,$13,
			$14,$15,COALESCE($16::jsonb, '{}'::jsonb),$17,$18,
			$19,COALESCE($20::inet, NULL),$21,$22,
			$23,$24,$25
		)
		RETURNING 
			"UserID","Email","EmailVerifiedAt","PhoneE164","PhoneVerifiedAt",
			"PasswordHash","PasswordAlg","PasswordUpdatedAt","MustChangePassword",
			"Status","IsServiceAccount","DisplayName","AvatarURL",
			"Locale","Timezone","Preferences","MFAEnrolled","MFADefaultMethod",
			"LastLoginAt","LastLoginIP","CreatedAt","CreatedBy",
			"UpdatedAt","UpdatedBy","IsDeleted"
	`

	row := r.db.QueryRow(ctx, q,
		u.UserID, u.Email, u.EmailVerifiedAt, u.PhoneE164, u.PhoneVerifiedAt,
		u.PasswordHash, u.PasswordAlg, u.PasswordUpdatedAt, u.MustChangePassword,
		u.Status, u.IsServiceAccount, u.DisplayName, u.AvatarURL,
		u.Locale, u.Timezone, u.Preferences, u.MFAEnrolled, u.MFADefaultMethod,
		u.LastLoginAt, u.LastLoginIP, u.CreatedAt, u.CreatedBy,
		u.UpdatedAt, u.UpdatedBy, u.IsDeleted,
	)

	var ur UserRow
	if err := row.Scan(
		&ur.UserID, &ur.Email, &ur.EmailVerifiedAt, &ur.PhoneE164, &ur.PhoneVerifiedAt,
		&ur.PasswordHash, &ur.PasswordAlg, &ur.PasswordUpdatedAt, &ur.MustChangePassword,
		&ur.Status, &ur.IsServiceAccount, &ur.DisplayName, &ur.AvatarURL,
		&ur.Locale, &ur.Timezone, &ur.Preferences, &ur.MFAEnrolled, &ur.MFADefaultMethod,
		&ur.LastLoginAt, &ur.LastLoginIP, &ur.CreatedAt, &ur.CreatedBy,
		&ur.UpdatedAt, &ur.UpdatedBy, &ur.IsDeleted,
	); err != nil {
		return nil, err
	}

	user, err := ur.ToDomain()
	if err != nil {
		return nil, err
	}
	return &user, nil
}
