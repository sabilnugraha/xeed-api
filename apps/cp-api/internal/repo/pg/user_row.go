// apps/cp-api/internal/repo/pg/user_row.go
package pg

import (
	"encoding/json"
	"net" // ganti dari "net/netip" ke "net"
	"time"
	"xeed/apps/cp-api/internal/domain"

	"github.com/google/uuid"
)

type UserRow struct {
	UserID             uuid.UUID
	Email              string
	EmailVerifiedAt    *time.Time
	PhoneE164          *string
	PhoneVerifiedAt    *time.Time
	PasswordHash       *string
	PasswordAlg        string
	PasswordUpdatedAt  *time.Time
	MustChangePassword bool
	Status             string
	IsServiceAccount   bool
	DisplayName        *string
	AvatarURL          *string
	Locale             string
	Timezone           string
	Preferences        []byte // JSONB
	MFAEnrolled        bool
	MFADefaultMethod   *string
	LastLoginAt        *time.Time
	LastLoginIP        *string // simpan string inet (pgx bisa map langsung)
	CreatedAt          time.Time
	CreatedBy          *uuid.UUID
	UpdatedAt          time.Time
	UpdatedBy          *uuid.UUID
	IsDeleted          bool
}

func (r *UserRow) ToDomain() (domain.User, error) {
	var prefs domain.Preferences
	if len(r.Preferences) > 0 {
		_ = json.Unmarshal(r.Preferences, &prefs)
	}

	// Samakan tipe dengan domain.User.LastLoginIP: *net.IP
	var ip *net.IP
	if r.LastLoginIP != nil {
		if parsed := net.ParseIP(*r.LastLoginIP); parsed != nil {
			ip = &parsed
		}
	}

	var mfa *domain.MFAMethod
	if r.MFADefaultMethod != nil {
		m := domain.MFAMethod(*r.MFADefaultMethod)
		mfa = &m
	}

	return domain.User{
		UserID:             r.UserID,
		Email:              r.Email,
		EmailVerifiedAt:    r.EmailVerifiedAt,
		PhoneE164:          r.PhoneE164,
		PhoneVerifiedAt:    r.PhoneVerifiedAt,
		PasswordHash:       r.PasswordHash,
		PasswordAlg:        domain.PasswordAlg(r.PasswordAlg),
		PasswordUpdatedAt:  r.PasswordUpdatedAt,
		MustChangePassword: r.MustChangePassword,
		Status:             domain.UserStatus(r.Status),
		IsServiceAccount:   r.IsServiceAccount,
		DisplayName:        r.DisplayName,
		AvatarURL:          r.AvatarURL,
		Locale:             r.Locale,
		Timezone:           r.Timezone,
		Preferences:        prefs,
		MFAEnrolled:        r.MFAEnrolled,
		MFADefaultMethod:   mfa,
		LastLoginAt:        r.LastLoginAt,
		LastLoginIP:        ip, // sekarang *net.IP
		CreatedAt:          r.CreatedAt,
		CreatedBy:          r.CreatedBy,
		UpdatedAt:          r.UpdatedAt,
		UpdatedBy:          r.UpdatedBy,
		IsDeleted:          r.IsDeleted,
	}, nil
}
