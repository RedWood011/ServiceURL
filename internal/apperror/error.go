package apperror

var (
	ErrNotFound = NewAppError(nil, "not found")
)

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
