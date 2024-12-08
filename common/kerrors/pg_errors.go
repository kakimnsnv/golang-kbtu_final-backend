package pg_errors

const (
	PgErrorNoRowsAffected PgErrorType = iota
)

const ()

type PgErrorType int

type PgError struct {
	Message   string `json:"message"`
	ErrorType PgErrorType
}

var _ error = (*PgError)(nil)

func (e *PgError) Error() string {
	return e.Message
}

func NewPgError(err PgErrorType) *PgError {
	var message string
	switch err {
	case PgErrorNoRowsAffected:
		message = "no rows affected"
	}

	return &PgError{
		Message:   message,
		ErrorType: err,
	}
}
