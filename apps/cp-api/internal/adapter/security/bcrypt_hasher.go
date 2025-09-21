// internal/adapters/security/bcrypt_hasher.go
package security

import (
	"time"
	"xeed/apps/cp-api/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

func (BcryptHasher) Hash(plain string) (string, domain.PasswordAlg, time.Time, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost) // cost 10-12 ok
	if err != nil {
		return "", "", time.Time{}, err
	}
	return string(b), domain.PasswordAlg("bcrypt"), time.Now().UTC(), nil
}

// (opsional) verifikasi untuk login
func (BcryptHasher) Verify(plain, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
