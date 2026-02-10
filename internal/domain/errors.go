package domain

type AppError struct {
	Status  string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(status, message string) *AppError {
	return &AppError{
		Status:  status,
		Message: message,
	}
}

var (
	ErrUserNotFound       = NewAppError("USER_NOT_FOUND", "user not found")
	ErrEmailAlreadyExists = NewAppError("EMAIL_EXISTS", "email already exists")
	ErrInvalidCredentials = NewAppError("INVALID_CREDENTIALS", "invalid credentials")
	ErrInvalidToken       = NewAppError("INVALID_TOKEN", "invalid token")
	ErrTokenBlacklisted   = NewAppError("TOKEN_BLACKLISTED", "token has been blacklisted")
	ErrRequestInvalid     = NewAppError("INVALID_INPUT", "request invalid")
)
