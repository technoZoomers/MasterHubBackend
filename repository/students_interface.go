package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type StudentsRepoI interface {
	GetStudentByUserId(student *models.StudentDB) error
}
