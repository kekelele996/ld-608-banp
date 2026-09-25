package services

import "groundTurn/src/types"

// BusinessError carries a stable error code plus rendered detail so the
// controller can surface exact flight/task/resource identifiers. Services
// wrap repository errors into this; controllers wrap it again for HTTP.
type BusinessError struct {
	Code       string
	Message    string
	Violation  *types.ReleaseViolation
	Violations []types.ReleaseViolation
}

func (e *BusinessError) Error() string { return e.Message }

func bizError(code, message string) *BusinessError {
	return &BusinessError{Code: code, Message: message}
}
