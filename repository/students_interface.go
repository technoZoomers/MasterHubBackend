package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type StudentsRepoI interface {
	GetStudentByUserId(student *models.StudentDB) error
	GetStudentIdByUsername(student *models.StudentDB) error
	InsertStudent(student *models.StudentDB) error
	UpdateStudent(student *models.StudentDB) error
	GetStudentByUserIdWithEmail(student *models.StudentDB) (string, error)
}
