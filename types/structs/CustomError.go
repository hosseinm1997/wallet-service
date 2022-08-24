package structs

import (
	"fmt"
)

type ErrorCategory int

var Categories = struct {
	BusinessLogic ErrorCategory
	Internal      ErrorCategory
	UnAuthorized  ErrorCategory
}{
	BusinessLogic: 1,
	Internal:      2,
	UnAuthorized:  3,
}

type CustomError struct {
	Err      string
	Category ErrorCategory
}

func (r CustomError) Error() string {
	return fmt.Sprintf(r.Err)
}
