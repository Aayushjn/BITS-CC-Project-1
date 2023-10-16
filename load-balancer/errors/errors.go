package errors

import "fmt"

type ErrBackendLimitOverflow struct{}

func (err *ErrBackendLimitOverflow) Error() string {
	return "max number of backends reached"
}

type ErrBackendLimitUnderflow struct{}

func (err *ErrBackendLimitUnderflow) Error() string {
	return "min number of backends reached"
}

type ErrBackendNotFound struct {
	BackendUrl string
}

func (err *ErrBackendNotFound) Error() string {
	return fmt.Sprintf("no backend found for %s", err.BackendUrl)
}
