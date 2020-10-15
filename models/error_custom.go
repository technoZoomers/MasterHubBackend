package models

import "fmt"

type NotFoundError struct {
	Message string
	RequestId int64
}

func (nfError* NotFoundError) Error() string {
	return fmt.Sprintf("%s with id: %d", nfError.Message, nfError.RequestId)
}