package useCases

import (
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
	"time"
)

type StudentsUC struct {
	useCases      *UseCases
	UsersRepo     repository.UsersRepoI
	MastersRepo   repository.MastersRepoI
	StudentsRepo  repository.StudentsRepoI
	LanguagesRepo repository.LanguagesRepoI
}

func (studentsUC *StudentsUC) insertStudentsLanguages(languages []string, studentDB *models.StudentDB) error {
	var newLanguagesIds []int64
	for _, language := range languages {
		languageDB := models.LanguageDB{Name: language}
		err := studentsUC.LanguagesRepo.GetLanguageByName(&languageDB)
		if err != nil {
			return fmt.Errorf(studentsUC.useCases.errorMessages.DbError)
		}
		if languageDB.Id == studentsUC.useCases.errorId {
			fileError := &models.BadRequestError{Message: "can't register student, language doesn't exist"}
			logger.Errorf(fileError.Error())
			return fileError
		}
		newLanguagesIds = append(newLanguagesIds, languageDB.Id)
	}
	err := studentsUC.UsersRepo.SetUserLanguagesById(studentDB.UserId, newLanguagesIds)
	if err != nil {
		return fmt.Errorf(studentsUC.useCases.errorMessages.DbError)
	}
	return nil
}

func (studentsUC *StudentsUC) Register(studentFull *models.StudentFull) error {
	var err error
	var studentDB models.StudentDB
	var userDB models.UserDB

	if studentFull.Email == "" {
		reqError := &models.BadRequestError{Message: "email can't be empty"}
		logger.Errorf(reqError.Error())
		return reqError
	} else {
		userDBEmailExists := models.UserDB{
			Email: studentFull.Email,
		}
		err = studentsUC.UsersRepo.GetUserByEmail(&userDBEmailExists)
		if err != nil {
			return fmt.Errorf(studentsUC.useCases.errorMessages.DbError)
		}
		if userDBEmailExists.Id != studentsUC.useCases.errorId {
			conflictError := &models.ConflictError{Message: "can't register student, email is already taken", ExistingContent: studentFull.Email}
			logger.Errorf(conflictError.Error())
			return conflictError
		}
	}
	userDB.Email = studentFull.Email
	userDB.Password = studentFull.Password
	userDB.Created = time.Now()
	userDB.Type = 2
	if studentFull.Username == "" {
		reqError := &models.BadRequestError{Message: "username can't be empty", RequestId: studentFull.UserId}
		logger.Errorf(reqError.Error())
		return reqError
	} else {
		masterDBUsernameExist := models.MasterDB{
			Username: studentFull.Username,
		}
		err = studentsUC.MastersRepo.GetMasterIdByUsername(&masterDBUsernameExist)
		if err != nil {
			return fmt.Errorf(studentsUC.useCases.errorMessages.DbError)
		}
		if masterDBUsernameExist.Id != studentsUC.useCases.errorId {
			conflictError := &models.ConflictError{Message: "can't register student, username is already taken", ExistingContent: studentFull.Username}
			logger.Errorf(conflictError.Error())
			return conflictError
		}
		studentDBUsernameExist := models.StudentDB{
			Username: studentFull.Username,
		}
		err = studentsUC.StudentsRepo.GetStudentIdByUsername(&studentDBUsernameExist)
		if err != nil {
			return fmt.Errorf(studentsUC.useCases.errorMessages.DbError)
		}
		if studentDBUsernameExist.Id != studentsUC.useCases.errorId {
			conflictError := &models.ConflictError{Message: "can't register student, username is already taken", ExistingContent: studentFull.Username}
			logger.Errorf(conflictError.Error())
			return conflictError
		}
	}
	studentDB.Username = studentFull.Username
	studentDB.Fullname = studentFull.Fullname
	err = studentsUC.UsersRepo.InsertUser(&userDB)
	if err != nil {
		return fmt.Errorf(studentsUC.useCases.errorMessages.DbError)
	}
	studentDB.UserId = userDB.Id
	err = studentsUC.StudentsRepo.InsertStudent(&studentDB)
	if err != nil {
		_ = studentsUC.UsersRepo.DeleteUserWithId(userDB.Id)
		return fmt.Errorf(studentsUC.useCases.errorMessages.DbError)
	}
	studentFull.Password = ""
	studentFull.UserId = userDB.Id

	err = studentsUC.insertStudentsLanguages(studentFull.Languages, &studentDB)
	if err != nil {
		return err
	}
	return nil
}

func (studentsUC *StudentsUC) setLanguages(student *models.Student, studentDB *models.StudentDB) error {
	var langs []string
	langsIds, err := studentsUC.UsersRepo.GetUserLanguagesById(studentDB.UserId)
	if err != nil {
		student.Languages = langs
		return fmt.Errorf(studentsUC.useCases.errorMessages.DbError)
	}
	for _, langId := range langsIds {
		var language models.LanguageDB
		language.Id = langId
		err = studentsUC.LanguagesRepo.GetLanguageById(&language)
		if err != nil {
			student.Languages = langs
			return fmt.Errorf(studentsUC.useCases.errorMessages.DbError)
		}
		if language.Name != "" {
			langs = append(langs, language.Name)
		}
	}
	student.Languages = langs
	return nil
}

func (studentsUC *StudentsUC) GetStudentById(student *models.Student) error {
	var studentDB models.StudentDB
	err := studentsUC.validateStudent(&studentDB, student)
	if err != nil {
		return err
	}
	err = studentsUC.matchStudent(&studentDB, student)
	if err != nil {
		return err
	}
	return nil
}
func (studentsUC *StudentsUC) matchStudent(studentDB *models.StudentDB, student *models.Student) error {
	student.Username = studentDB.Username
	student.Fullname = studentDB.Fullname

	err := studentsUC.setLanguages(student, studentDB)
	if err != nil {
		return err
	}
	return nil
}
func (studentsUC *StudentsUC) validateStudent(studentDB *models.StudentDB, student *models.Student) error {
	if student.UserId == studentsUC.useCases.errorId {
		return &models.BadRequestError{Message: "incorrect student id", RequestId: student.UserId}
	}
	studentDB.UserId = student.UserId
	err := studentsUC.StudentsRepo.GetStudentByUserId(studentDB)
	if err != nil {
		return fmt.Errorf(studentsUC.useCases.errorMessages.DbError)
	}
	if studentDB.Id == studentsUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "student doesn't exist", RequestId: student.UserId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	return nil
}

func (studentsUC *StudentsUC) ChangeStudentData(student *models.Student) error {
	if student.UserId == studentsUC.useCases.errorId {
		return &models.BadRequestError{Message: "incorrect student id", RequestId: student.UserId}
	}
	studentDB := models.StudentDB{
		UserId: student.UserId,
	}
	err := studentsUC.StudentsRepo.GetStudentByUserId(&studentDB)
	if err != nil {
		return fmt.Errorf(studentsUC.useCases.errorMessages.DbError)
	}
	if studentDB.Id == studentsUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "student doesn't exist", RequestId: student.UserId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}

	if student.Username != "" {
		masterDBUsernameExist := models.MasterDB{
			Username: student.Username,
		}
		err = studentsUC.MastersRepo.GetMasterIdByUsername(&masterDBUsernameExist)
		if err != nil {
			return fmt.Errorf(studentsUC.useCases.errorMessages.DbError)
		}
		if masterDBUsernameExist.Id != studentsUC.useCases.errorId && masterDBUsernameExist.Id != studentDB.Id {
			absenceError := &models.ConflictError{Message: "can't update student, username is already taken", ExistingContent: student.Username}
			logger.Errorf(absenceError.Error())
			return absenceError
		}
		studentDBUsernameExist := models.StudentDB{
			Username: student.Username,
		}
		err = studentsUC.StudentsRepo.GetStudentIdByUsername(&studentDBUsernameExist)
		if err != nil {
			return fmt.Errorf(studentsUC.useCases.errorMessages.DbError)
		}
		if studentDBUsernameExist.Id != studentsUC.useCases.errorId && studentDBUsernameExist.Id != studentDB.Id {
			conflictError := &models.ConflictError{Message: "can't update student, username is already taken", ExistingContent: student.Username}
			logger.Errorf(conflictError.Error())
			return conflictError
		}
		studentDB.Username = student.Username
	} else {
		student.Username = studentDB.Username
	}

	studentDB.Fullname = student.Fullname

	err = studentsUC.StudentsRepo.UpdateStudent(&studentDB)
	if err != nil {
		return fmt.Errorf(studentsUC.useCases.errorMessages.DbError)
	}

	err = studentsUC.changeStudentsLanguages(student, &studentDB)
	if err != nil {
		return err
	}

	return nil
}

func (studentsUC *StudentsUC) changeStudentsLanguages(student *models.Student, studentDB *models.StudentDB) error {
	var newLanguagesIds []int64
	for _, language := range student.Languages {
		languageDB := models.LanguageDB{Name: language}
		err := studentsUC.LanguagesRepo.GetLanguageByName(&languageDB)
		if err != nil {
			return fmt.Errorf(studentsUC.useCases.errorMessages.DbError)
		}
		if languageDB.Id == studentsUC.useCases.errorId {
			fileError := &models.BadRequestError{Message: "cant't update student, language doesn't exist", RequestId: student.UserId}
			logger.Errorf(fileError.Error())
			return fileError
		}
		newLanguagesIds = append(newLanguagesIds, languageDB.Id)
	}

	oldLanguagesIds, err := studentsUC.UsersRepo.GetUserLanguagesById(studentDB.UserId)
	if err != nil {
		return fmt.Errorf(studentsUC.useCases.errorMessages.DbError)
	}

	err = studentsUC.UsersRepo.DeleteUserLanguagesById(studentDB.UserId)
	if err != nil {
		return fmt.Errorf(studentsUC.useCases.errorMessages.DbError)
	}

	err = studentsUC.UsersRepo.SetUserLanguagesById(studentDB.UserId, newLanguagesIds)
	if err != nil {
		_ = studentsUC.UsersRepo.SetUserLanguagesById(studentDB.UserId, oldLanguagesIds)
		return fmt.Errorf(studentsUC.useCases.errorMessages.DbError)
	}
	return nil
}
