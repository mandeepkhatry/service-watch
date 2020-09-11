package def

import "errors"

var (
	ErrSwaggerConfigUnregistered = errors.New("Swagger Config Unregistered")
	ErrServersUnregistered       = errors.New("Servers Config Unregistered")
)
