package apperror

var (
	ErrNotFound  = NewAppError(nil, "not found")
	ErrNoContent = NewAppError(nil, "no content")
	ErrConflict  = NewAppError(nil, "conflict database")
	ErrDataBase  = NewAppError(nil, "error write database")
	ErrGone      = NewAppError(nil, "error not available url")
)

// AppError Кастомная структура ошибок приложения
type AppError struct {
	Err     error  `json:"-"`
	Message string `json:"Message,omitempty"`
}

func NewAppError(err error, message string) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}
