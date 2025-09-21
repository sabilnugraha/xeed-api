package system

import (
	"xeed/apps/cp-api/internal/usecase/contract"

	"github.com/google/uuid"
)

// IDGen: generator UUID v4
type IDGen struct{}

func (IDGen) New() uuid.UUID { return uuid.New() }

// Compile-time check
var _ contract.IDGen = (*IDGen)(nil)
