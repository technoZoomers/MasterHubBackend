package useCases

import "github.com/technoZoomers/MasterHubBackend/models"

type StudentsUCInterface interface {
	Register(studentFull *models.StudentFull) error
	GetStudentById(student *models.Student) error
	ChangeStudentData(student *models.Student) error
}
