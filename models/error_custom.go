package models

import "fmt"

type BadRequestError struct {
	Message string
	RequestId int64
}

func (brError* BadRequestError) Error() string {
	return fmt.Sprintf("%s with id: %d", brError.Message, brError.RequestId)
}

type ConflictError struct {
	Message string
	RequestId int64
}

func (conflictError* ConflictError) Error() string {
	return fmt.Sprintf("%s with id: %d", conflictError.Message, conflictError.RequestId)
}