package useCases

type ErrorMessagesUC struct {
	DbError             string
	InternalServerError string
	FileErrors          FileErrors
	MailSendError       string
}

type FileErrors struct {
	FileOpenError          string
	FileReadError          string
	FileReadExtensionError string
	FileCreateError        string
	FileRemoveError        string
	FileGenerateError      string
	FileConvertError       string
}
