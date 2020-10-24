package useCases

import (
	"errors"
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
	"time"
)

type ChatsUC struct {
	useCases     *UseCases
	ChatsRepo   repository.ChatsRepoI
	MastersRepo repository.MastersRepoI
	StudentsRepo repository.StudentsRepoI
	badRequestError *models.BadRequestError
	chatsConfig ChatsConfig
}

type ChatsConfig struct {
	userMap map[string]int64
	chatTypes map[string]int64

	userMapBackwards map[int64]string
	chatTypesBackwards  map[int64]string
}

func (chatsUC *ChatsUC) validateMaster(masterId int64) error {
	if masterId == chatsUC.useCases.errorId {
		return &models.BadRequestError{Message: "incorrect master id", RequestId: masterId}
	}
	masterDB := models.MasterDB{
		UserId: masterId,
	}
	err := chatsUC.MastersRepo.GetMasterByUserId(&masterDB)
	if err != nil {
		return fmt.Errorf(chatsUC.useCases.errorMessages.DbError)
	}
	if masterDB.Id == chatsUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "master doesn't exist", RequestId: masterId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	return nil
}

func (chatsUC *ChatsUC) validateStudent(studentId int64) error {
	if studentId == chatsUC.useCases.errorId {
		return &models.BadRequestError{Message: "incorrect master id", RequestId: studentId}
	}
	studentDB := models.StudentDB{
		UserId: studentId,
	}
	err := chatsUC.StudentsRepo.GetStudentByUserId(&studentDB)
	if err != nil {
		return fmt.Errorf(chatsUC.useCases.errorMessages.DbError)
	}
	if studentDB.Id == chatsUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "student doesn't exist", RequestId: studentId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	return nil
}

func (chatsUC *ChatsUC) matchChat (chatDB *models.ChatDB, chat *models.Chat) error {
	chat.Id = chatDB.Id
	chat.Type = chatDB.Type
	chat.Created = chatDB.Created
	chat.StudentId = chatDB.StudentId
	chat.MasterId = chatDB.MasterId
	return nil
}

func (chatsUC *ChatsUC) matchChatsQuery(userId int64, userType int64, query *models.ChatsQueryValues, queryDB *models.ChatsQueryValuesDB) {
	queryDB.Offset = query.Offset
	queryDB.Limit = query.Limit
	queryDB.UserId = userId
	queryDB.Type = query.Type
	queryDB.User = userType
}

func (chatsUC *ChatsUC) checkDeleted(userType int64, chatType int64) bool {
	if (chatsUC.chatsConfig.userMap["student"] == userType &&
		chatsUC.chatsConfig.chatTypes["deleted by student"] == chatType) ||
		(chatsUC.chatsConfig.userMap["master"] == userType &&
			chatsUC.chatsConfig.chatTypes["deleted by master"] == chatType){
		return true
	} else {
		return  false
	}
}

func (chatsUC *ChatsUC) GetUserChatsById(userId int64, query models.ChatsQueryValues) (models.Chats, error) {
	var queryDB models.ChatsQueryValuesDB
	chats := make([]models.Chat, 0)
	var userType int64 = 0
	err := chatsUC.validateMaster(userId)
	if err != nil {
		if errors.As(err, &chatsUC.badRequestError) {
			err := chatsUC.validateStudent(userId)
			if err != nil {
				return chats, err
			}
			userType = chatsUC.chatsConfig.userMap["student"]
		} else {
			return chats, err
		}
	} else {
		userType = chatsUC.chatsConfig.userMap["master"]
	}

	chatsUC.matchChatsQuery(userId, userType, &query, &queryDB)
	chatsDB, err := chatsUC.ChatsRepo.GetChatsByUserId(queryDB)
	if err != nil {
		return chats, fmt.Errorf(chatsUC.useCases.errorMessages.DbError)
	}
	for _, chatDB := range chatsDB {
		if !chatsUC.checkDeleted(userType, chatDB.Type) {
			var chat models.Chat
			err = chatsUC.matchChat(&chatDB, &chat)
			if err != nil {
				return chats, err
			}
			chats = append(chats, chat)
		}
	}
	return chats, nil
}

func (chatsUC *ChatsUC) CreateChatRequest(chatRequest *models.Chat, studentId int64) error {
	err := chatsUC.validateStudent(chatRequest.StudentId)
	if err != nil {
		return err
	}
	err = chatsUC.validateMaster(chatRequest.MasterId)
	if err != nil {
		return err
	} // TODO: check if chat exists
	if chatRequest.StudentId != studentId {
		accessError := &models.BadRequestError{Message: "can't create other student's chat", RequestId: studentId}
		logger.Errorf(accessError.Error())
		return accessError
	}
	chatDB := models.ChatDB{
		Type:  1,
		StudentId: chatRequest.StudentId,
		MasterId:  chatRequest.MasterId,
		Created:  time.Now(),
	}
	err = chatsUC.ChatsRepo.InsertChatRequest(&chatDB)
	if err != nil {
		return fmt.Errorf(chatsUC.useCases.errorMessages.DbError)
	}
	chatRequest.Type = chatDB.Type
	chatRequest.Created = chatDB.Created
	chatRequest.Id = chatDB.Id
	return  nil
}