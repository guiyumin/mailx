//go:build !darwin || !cgo

package calendar

// AuthStatus represents the current authorization status
type AuthStatus int

const (
	AuthNotDetermined AuthStatus = iota
	AuthRestricted
	AuthDenied
	AuthAuthorized
)

// GetAuthStatus returns the current calendar authorization status
func GetAuthStatus() AuthStatus {
	return AuthDenied
}

func NewClient() (Client, error) {
	return nil, ErrNotSupported
}
