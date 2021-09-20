package uuid

import (
	"github.com/SKorolchuk/dpio-workspace/internal/pkg/infra/validation"

	"github.com/google/uuid"
)

// Generate forms new unique identifier in GUID format.
func Generate() string {
	return uuid.NewString()
}

// Validate checks GUID format of a target id.
func Validate(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return validation.ErrorInvalidIdentifier
	}
	return nil
}
