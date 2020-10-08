package def

import "errors"

var (
	ErrSwaggerConfigUnregistered = errors.New("Swagger Config Unregistered")
	ErrServersUnregistered       = errors.New("Servers Config Unregistered")
	ErrStoreUnavailable          = errors.New("Store unavailable")
	ErrInvalidFormFileSyntax     = errors.New("invalid form file syntax")
)
