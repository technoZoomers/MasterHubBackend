package useCases

import (
	"errors"
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
	"strconv"
	"time"
)

type ChatsUC struct {
	useCases        *UseCases
	ChatsRepo       repository.ChatsRepoI
	MastersRepo     repository.MastersRepoI
	StudentsRepo    repository.StudentsRepoI
	badRequestError *models.BadRequestError
	chatsConfig     ChatsConfig
}

type ChatsConfig struct {
	userMap   map[string]int64
	chatTypes map[string]int64

	userMapBackwards   map[int64]string
	chatTypesBackwards map[int64]string

	messagesTypesMap          map[int64]bool
	messagesTypesMapBackwards map[bool]int64
}

func (chatsUC *ChatsUC) validateChat(chatExists *models.ChatDB, chatId int64) error {
	if chatId == chatsUC.useCases.errorId {
		return &models.BadRequestError{Message: "incorrect chat id", RequestId: chatExists.Id}
	}
	err := chatsUC.ChatsRepo.GetChatById(chatExists, chatId)
	if err != nil {
		return fmt.Errorf(chatsUC.useCases.errorMessages.DbError)
	}
	if chatExists.Id == chatsUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "chat doesn't exist", RequestId: chatId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	return nil
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

func (chatsUC *ChatsUC) matchMessage(messageDB *models.MessageDB, message *models.Message) error {
	message.Id = messageDB.Id
	if messageDB.UserId != chatsUC.useCases.errorId {
		message.AuthorId = messageDB.UserId
	}
	message.ChatId = messageDB.ChatId
	message.Created = messageDB.Created
	message.Text = messageDB.Text
	message.Type = chatsUC.chatsConfig.messagesTypesMapBackwards[messageDB.Info]
	return nil
}

func (chatsUC *ChatsUC) matchChat(chatDB *models.ChatDB, chat *models.Chat) error {
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
			chatsUC.chatsConfig.chatTypes["deleted by master"] == chatType) {
		return true
	} else {
		return false
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
	if chatRequest.StudentId != studentId {
		accessError := &models.ForbiddenError{Reason: "can't create other student's chat"}
		logger.Errorf(accessError.Error())
		return accessError
	}
	chatDB, err := chatsUC.createChat(chatRequest.StudentId, chatRequest.MasterId, chatsUC.chatsConfig.chatTypes["unseen"])
	if err != nil {
		return err
	}
	err = chatsUC.matchChat(&chatDB, chatRequest)
	if err != nil {
		return err
	}
	return nil
}

func (chatsUC *ChatsUC) CreateChatByMaster(chatRequest *models.Chat, masterId int64) error {
	if chatRequest.MasterId != masterId {
		accessError := &models.ForbiddenError{Reason: "can't create other master's chat"}
		logger.Errorf(accessError.Error())
		return accessError
	}
	chatDB, err := chatsUC.createChat(chatRequest.StudentId, chatRequest.MasterId, chatsUC.chatsConfig.chatTypes["approved"])
	if err != nil {
		return err
	}
	err = chatsUC.matchChat(&chatDB, chatRequest)
	if err != nil {
		return err
	}
	return nil
}

func (chatsUC *ChatsUC) createChat(studentId int64, masterId int64, chatType int64) (models.ChatDB, error) {
	var chatDB models.ChatDB
	err := chatsUC.validateStudent(studentId)
	if err != nil {
		return chatDB, err
	}
	err = chatsUC.validateMaster(masterId)
	if err != nil {
		return chatDB, err
	}

	if chatsUC.chatsConfig.chatTypesBackwards[chatType] == "" {
		reqError := &models.BadRequestError{Message: "wrong chat type provided"}
		logger.Errorf(reqError.Error())
		return chatDB, reqError
	}
	chatDB = models.ChatDB{
		Type:      chatType,
		StudentId: studentId,
		MasterId:  masterId,
		Created:   time.Now(),
	}
	err = chatsUC.ChatsRepo.GetChatByStudentIdAndMasterId(&chatDB)
	if err != nil {
		return chatDB, fmt.Errorf(chatsUC.useCases.errorMessages.DbError)
	}
	if chatDB.Id != chatsUC.useCases.errorId {
		existsError := &models.ConflictError{Message: "chat already exists", ExistingContent: strconv.FormatInt(chatDB.Id, 10)}
		logger.Errorf(existsError.Error())
		return chatDB, existsError
	}
	err = chatsUC.ChatsRepo.InsertChatRequest(&chatDB)
	if err != nil {
		return chatDB, fmt.Errorf(chatsUC.useCases.errorMessages.DbError)
	}
	return chatDB, nil
}

func (chatsUC *ChatsUC) ChangeChatStatus(chat *models.Chat, masterId int64, chatId int64) error {
	if chat.Id == chatsUC.useCases.errorId {
		chat.Id = chatId
	} else {
		if chat.Id != chatId {
			matchError := &models.BadRequestError{Message: "chat id doesn't match", RequestId: chat.Id}
			logger.Errorf(matchError.Error())
			return matchError
		}
	}
	if chat.MasterId == chatsUC.useCases.errorId {
		chat.MasterId = masterId
	} else {
		if chat.MasterId != masterId {
			matchError := &models.ForbiddenError{Reason: "can't change other masters's chat status"}
			logger.Errorf(matchError.Error())
			return matchError
		}
	}
	if chat.Type != chatsUC.chatsConfig.chatTypes["approved"] &&
		chat.Type != chatsUC.chatsConfig.chatTypes["disapproved"] {
		reqError := &models.BadRequestError{Message: "wrong chat type provided", RequestId: chat.Id}
		logger.Errorf(reqError.Error())
		return reqError
	}
	var chatExists models.ChatDB
	err := chatsUC.validateChat(&chatExists, chat.Id)
	if err != nil {
		return err
	}
	if chat.StudentId != chatsUC.useCases.errorId {
		if chatExists.StudentId != chat.StudentId {
			if chat.MasterId != masterId {
				reqError := &models.BadRequestError{Message: "can't change chat's student", RequestId: chat.Id}
				logger.Errorf(reqError.Error())
				return reqError
			}
		}
	} else {
		chat.StudentId = chatExists.StudentId
	}
	if chatExists.MasterId != chat.MasterId {
		if chat.MasterId != masterId {
			reqError := &models.BadRequestError{Message: "can't change chat's master", RequestId: chat.Id}
			logger.Errorf(reqError.Error())
			return reqError
		}
	}
	if chatExists.Type != chatsUC.chatsConfig.chatTypes["unseen"] {
		reqError := &models.BadRequestError{Message: "can't change chat's type", RequestId: chat.Id}
		logger.Errorf(reqError.Error())
		return reqError
	}
	if !chat.Created.IsZero() {
		if !chatExists.Created.Equal(chat.Created) {
			requestError := fmt.Errorf("chat creation time can't be changed")
			logger.Errorf(requestError.Error())
			return requestError
		}
	} else {
		chat.Created = chatExists.Created
	}
	chatExists.Type = chat.Type
	err = chatsUC.ChatsRepo.ChangeChatType(&chatExists)
	if err != nil {
		return err
	}
	return nil
}

func (chatsUC *ChatsUC) GetMessagesByChatId(chatId int64) (models.Messages, error) {
	messages := make([]models.Message, 0)
	//var chatExists models.ChatDB
	//err := chatsUC.validateChat(&chatExists, chatId)
	//if err != nil {
	//	return messages, err
	//}
	messagesDB, err := chatsUC.ChatsRepo.GetMessagesByChatId(chatId)
	if err != nil {
		return messages, fmt.Errorf(chatsUC.useCases.errorMessages.DbError)
	}
	for _, messageDB := range messagesDB {
		var message models.Message
		err = chatsUC.matchMessage(&messageDB, &message)
		if err != nil {
			return messages, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (chatsUC *ChatsUC) CheckChatByUserId(chatId int64, userId int64) error {
	var chatExists models.ChatDB
	if chatId == chatsUC.useCases.errorId {
		return &models.BadRequestError{Message: "incorrect chat id", RequestId: chatId}
	}
	err := chatsUC.ChatsRepo.GetChatByIdAndMasterOrStudentId(&chatExists, chatId, userId)
	if err != nil {
		return fmt.Errorf(chatsUC.useCases.errorMessages.DbError)
	}
	if chatExists.Id == chatsUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "chat doesn't exist", RequestId: chatId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	return nil
}
