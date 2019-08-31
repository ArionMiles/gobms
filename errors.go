package gobms

import "fmt"

type goBMSError struct {
	Message string
}

func (e goBMSError) Error() string {
	return fmt.Sprintf("%v", e.Message)
}
