package system

import (
	"time"
	"xeed/apps/cp-api/internal/usecase/contract"
)

// Clock: implementasi sederhana -> selalu UTC
type Clock struct{}

func (Clock) Now() time.Time { return time.Now().UTC() }

// Compile-time check
var _ contract.Clock = (*Clock)(nil)
