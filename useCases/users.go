package useCases

import (
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
)

type UsersUC struct {
	useCases    *UseCases
	UsersRepo   repository.UsersRepoI
	CookiesRepo repository.CookiesRepoI
}

func (usersUC *UsersUC) GetUserById(user *models.User) error {
	var userDB models.UserDB
	err := usersUC.validateUser(&userDB, user.Id)
	if err != nil {
		return err
	}
	usersUC.matchUser(&userDB, user)
	return nil
}
func (usersUC *UsersUC) matchUser(userDB *models.UserDB, user *models.User) {
	user.Id = userDB.Id
	user.Email = userDB.Email
	user.Password = userDB.Password // TODO: HASH
	user.Type = userDB.Type
	user.Created = userDB.Created
}
func (usersUC *UsersUC) validateUser(userDB *models.UserDB, userId int64) error {
	if userId == usersUC.useCases.errorId {
		return &models.BadRequestError{Message: "incorrect user id", RequestId: userId}
	}
	err := usersUC.UsersRepo.GetUserById(userDB, userId)
	if err != nil {
		return fmt.Errorf(usersUC.useCases.errorMessages.DbError)
	}
	if userDB.Id == usersUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "user doesn't exist", RequestId: userId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	return nil
}

func (usersUC *UsersUC) validateUserId(userId int64) error {
	var userDB models.UserDB
	return usersUC.validateUser(&userDB, userId)
}

func (usersUC *UsersUC) Login(user *models.User) error {
	userDB := models.UserDB{
		Email:    user.Email,
		Password: user.Password,
	}
	err := usersUC.UsersRepo.GetUserByEmailAndPassword(&userDB)
	if err != nil {
		return fmt.Errorf(usersUC.useCases.errorMessages.DbError)
	}
	if userDB.Id == usersUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "введен неправильный email или пароль"}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	user.Id = userDB.Id
	user.Type = userDB.Type
	user.Created = userDB.Created
	return nil
}

func (usersUC *UsersUC) GetUserByCookie(cookieValue string, user *models.User) error {
	var cookieDB models.CookieDB
	var userDB models.UserDB
	err := usersUC.CookiesRepo.GetUserByCookie(cookieValue, &cookieDB)
	if err != nil {
		return fmt.Errorf(usersUC.useCases.errorMessages.DbError)
	}
	if cookieDB.User == usersUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "no user with this cookie exists"}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	err = usersUC.UsersRepo.GetUserById(&userDB, cookieDB.User)
	if err != nil {
		return fmt.Errorf(usersUC.useCases.errorMessages.DbError)
	}
	if userDB.Id == usersUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "user was deleted"}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	usersUC.matchUser(&userDB, user)
	return nil
}

func (usersUC *UsersUC) InsertCookie(userId int64, cookieValue string) error {
	var cookieInfo models.CookieDB
	cookieInfo.User = userId
	cookieInfo.Cookie = cookieValue
	err := usersUC.CookiesRepo.InsertCookie(&cookieInfo)
	if err != nil {
		return fmt.Errorf(usersUC.useCases.errorMessages.DbError)
	}
	return nil
}

func (usersUC *UsersUC) DeleteCookie(cookieValue string) error {
	err := usersUC.CookiesRepo.DeleteCookie(cookieValue)
	if err != nil {
		return fmt.Errorf(usersUC.useCases.errorMessages.DbError)
	}
	return nil
}
