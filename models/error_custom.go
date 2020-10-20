package models

import "fmt"

type BadRequestError struct {
	Message string
	RequestId int64
}

func (brError* BadRequestError) Error() string {
	return fmt.Sprintf("%s with id: %d", brError.Message, brError.RequestId)
}

type BadQueryParameterError struct {
	Parameter string
}

func (brError* BadQueryParameterError) Error() string {
	return fmt.Sprintf("bad query with parameter: %s", brError.Parameter)
}

type ConflictError struct {
	Message string
	RequestId int64
}

func (conflictError* ConflictError) Error() string {
	return fmt.Sprintf("%s with id: %d", conflictError.Message, conflictError.RequestId)
}

type NoContentError struct {
	Message string
	RequestId int64
}

func (noContentError* NoContentError) Error() string {
	return fmt.Sprintf("%s with id: %d", noContentError.Message, noContentError.RequestId)
}