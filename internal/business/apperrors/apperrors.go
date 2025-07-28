package apperrors

type AppError struct {
	Msg string
}

func (e *AppError) Error() string {
	return e.Msg
}

type NotFoundError struct {
	Msg string
}

func (e *NotFoundError) Error() string {
	return e.Msg
}

type ServerBusyError struct {
	Msg string
}

func (e *ServerBusyError) Error() string {
	return e.Msg
}
