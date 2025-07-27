package apperrors

type BusinessRuleViolationError struct {
	Msg string
}

func (e *BusinessRuleViolationError) Error() string {
	return e.Msg
}

type NotFoundError struct {
	Msg string
}

func (e *NotFoundError) Error() string {
	return e.Msg
}
