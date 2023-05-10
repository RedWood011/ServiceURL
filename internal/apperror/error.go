// Package apperror apperror предназначен для кастомных ошибок сервиса
package apperror

// Кастомные ошибки
var (
	// ErrNotFound Ошибка
	ErrNotFound = NewAppError(nil, "not found")
	// ErrNoContent Ошибка
	ErrNoContent = NewAppError(nil, "no content")
	// ErrConflict Ошибка
	ErrConflict = NewAppError(nil, "conflict database")
	// ErrDataBase Ошибка
	ErrDataBase = NewAppError(nil, "error write database")
	// ErrGone Ошибка
	ErrGone = NewAppError(nil, "error not available url")
)

// AppError Кастомная структура ошибок приложения
type AppError struct {
	Err     error  `json:"-"`
	Message string `json:"Message,omitempty"`
}

// NewAppError
func NewAppError(err error, message string) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
	}
}

// - Error
func (e *AppError) Error() string {
	return e.Message
}

// Unwrap
func (e *AppError) Unwrap() error {
	return e.Err
}
