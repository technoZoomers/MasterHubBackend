package useCases

type ErrorMessagesUC struct {
	DbError             string
	InternalServerError string
	FileErrors          FileErrors
}

type FileErrors struct {
	FileOpenError          string
	FileReadError          string
	FileReadExtensionError string
	FileCreateError        string
	FileRemoveError        string
}
