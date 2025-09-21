// apps/cp-api/internal/domain/user.go
package domain

import (
	"errors"
	"net"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	UserPending   UserStatus = "PENDING"
	UserActive    UserStatus = "ACTIVE"
	UserLocked    UserStatus = "LOCKED"
	UserSuspended UserStatus = "SUSPENDED"
	UserDeleted   UserStatus = "DELETED"
)

type MFAMethod string

const (
	MFATOTP     MFAMethod = "totp"
	MFAWebAuthn MFAMethod = "webauthn"
	MFASMS      MFAMethod = "sms"
	MFAEmail    MFAMethod = "email"
)

type PasswordAlg string

const (
	AlgArgon2id PasswordAlg = "argon2id"
	AlgBcrypt   PasswordAlg = "bcrypt"
	AlgScrypt   PasswordAlg = "scrypt"
	AlgExternal PasswordAlg = "external"
	AlgNone     PasswordAlg = "none"
)

var rxEmail = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
var rxE164 = regexp.MustCompile(`^\+\d{6,15}$`)

type Preferences map[string]any

type User struct {
	UserID          uuid.UUID
	Email           string
	EmailVerifiedAt *time.Time

	PhoneE164       *string
	PhoneVerifiedAt *time.Time

	PasswordHash       *string
	PasswordAlg        PasswordAlg
	PasswordUpdatedAt  *time.Time
	MustChangePassword bool

	Status           UserStatus
	IsServiceAccount bool

	DisplayName *string
	AvatarURL   *string

	Locale   string // default 'id-ID'
	Timezone string // default 'Asia/Jakarta'

	Preferences Preferences // JSONB '{}'

	MFAEnrolled      bool
	MFADefaultMethod *MFAMethod

	LastLoginAt *time.Time
	LastLoginIP *net.IP

	CreatedAt time.Time
	CreatedBy *uuid.UUID
	UpdatedAt time.Time
	UpdatedBy *uuid.UUID

	IsDeleted bool
}

func NewUser(email, displayName string) (User, error) {
	if !rxEmail.MatchString(email) {
		return User{}, errors.New("invalid email")
	}
	u := User{
		UserID:      uuid.New(),
		Email:       email,
		DisplayName: &displayName,
		Status:      UserPending,
		PasswordAlg: AlgArgon2id,
		Locale:      "id-ID",
		Timezone:    "Asia/Jakarta",
		Preferences: map[string]any{},
	}
	return u, nil
}

// === Invariants / transitions ===

func (u *User) ChangeEmail(newEmail string) error {
	if !rxEmail.MatchString(newEmail) {
		return errors.New("invalid email")
	}
	u.Email = newEmail
	u.EmailVerifiedAt = nil
	return nil
}

func (u *User) VerifyEmail(at time.Time) { u.EmailVerifiedAt = &at }

func (u *User) SetPhoneE164(e164 string) error {
	if !rxE164.MatchString(e164) {
		return errors.New("invalid phone (E.164)")
	}
	u.PhoneE164 = &e164
	u.PhoneVerifiedAt = nil
	return nil
}

func (u *User) VerifyPhone(at time.Time) { u.PhoneVerifiedAt = &at }

func (u *User) SetPasswordHash(hash string, at time.Time, mustChange bool) {
	u.PasswordHash = &hash
	u.PasswordUpdatedAt = &at
	u.MustChangePassword = mustChange
}

func (u *User) RequirePasswordChange() { u.MustChangePassword = true }

// Status transitions (kontrol sesuai enum DB)
func (u *User) Activate()   { u.Status = UserActive }
func (u *User) Lock()       { u.Status = UserLocked }
func (u *User) Suspend()    { u.Status = UserSuspended }
func (u *User) SoftDelete() { u.Status = UserDeleted; u.IsDeleted = true }

func (u *User) EnableMFA(method MFAMethod) {
	u.MFAEnrolled = true
	u.MFADefaultMethod = &method
}
func (u *User) DisableMFA() {
	u.MFAEnrolled = false
	u.MFADefaultMethod = nil
}

func (u *User) SetLocaleTimezone(locale, tz string) {
	if locale != "" {
		u.Locale = locale
	}
	if tz != "" {
		u.Timezone = tz
	}
}

func (u *User) SetProfile(displayName, avatarURL *string) {
	u.DisplayName = displayName
	u.AvatarURL = avatarURL
}

func (p Preferences) GetString(key, def string) string {
	if v, ok := p[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return def
}
func (p Preferences) Set(key string, val any) { p[key] = val }
