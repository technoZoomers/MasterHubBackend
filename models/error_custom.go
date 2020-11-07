package models

import "fmt"

type BadRequestError struct {
	Message   string
	RequestId int64
}

func (brError *BadRequestError) Error() string {
	return fmt.Sprintf("%s with request id: %d", brError.Message, brError.RequestId)
}

type BadQueryParameterError struct {
	Parameter string
}

func (brError *BadQueryParameterError) Error() string {
	return fmt.Sprintf("bad query with parameter: %s", brError.Parameter)
}

type ConflictError struct {
	Message         string
	ExistingContent string
}

func (conflictError *ConflictError) Error() string {
	return fmt.Sprintf("%s with conflict content: %s", conflictError.Message, conflictError.ExistingContent)
}

type NoContentError struct {
	Message   string
	RequestId int64
}

func (noContentError *NoContentError) Error() string {
	return fmt.Sprintf("%s with request id: %d", noContentError.Message, noContentError.RequestId)
}

type ForbiddenError struct {
	Reason string
}

func (forbiddenError *ForbiddenError) Error() string {
	return fmt.Sprintf("forbidden resource, reason: %s ", forbiddenError.Reason)
}
