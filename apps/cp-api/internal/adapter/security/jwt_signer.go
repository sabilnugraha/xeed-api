package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTSigner struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTSigner(secret string, ttl time.Duration) *JWTSigner {
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	return &JWTSigner{secret: []byte(secret), ttl: ttl}
}
func (s *JWTSigner) Sign(userID uuid.UUID, email string, now time.Time) (string, error) {
	claims := jwt.MapClaims{"sub": userID.String(), "email": email, "iat": now.Unix(), "exp": now.Add(s.ttl).Unix(), "iss": "cp-api", "aud": "cp-api"}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}
