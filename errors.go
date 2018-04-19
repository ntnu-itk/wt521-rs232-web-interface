package main

type customError struct {
	msg string
}

func NewError(msg string) error {
	return error(customError{msg: msg})
}
func (ce customError) Error() string {
	return ce.msg
}
