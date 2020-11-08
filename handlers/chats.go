package handlers

import (
	"fmt"
	"github.com/google/logger"
	json "github.com/mailru/easyjson"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"net/http"
	"net/url"
	"strconv"
)

type ChatsHandlers struct {
	handlers       *Handlers
	ChatsUC        useCases.ChatsUCInterface
	ChatsQueryKeys ChatsQueryKeys
}

type ChatsQueryKeys struct {
	Type   string
	Limit  string
	Offset string
}

func (ch *ChatsHandlers) parseChatsQuery(query url.Values, chatsQuery *models.ChatsQueryValues) error {
	offsetString := query.Get(ch.ChatsQueryKeys.Offset)
	if offsetString != "" {
		offset, err := strconv.ParseInt(offsetString, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing query parameter %s: %v", ch.ChatsQueryKeys.Offset, err.Error())
		}
		chatsQuery.Offset = offset
	}
	limitString := query.Get(ch.ChatsQueryKeys.Limit)
	if limitString != "" {
		limit, err := strconv.ParseInt(limitString, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing query parameter %s: %v", ch.ChatsQueryKeys.Limit, err.Error())
		}
		chatsQuery.Limit = limit
	}
	typeString := query.Get(ch.ChatsQueryKeys.Type)
	if typeString != "" {
		typeInt, err := strconv.ParseInt(typeString, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing query parameter %s: %v", ch.ChatsQueryKeys.Type, err.Error())
		}
		chatsQuery.Type = typeInt
	}
	return nil
}

func (ch *ChatsHandlers) GetChatsByUserId(writer http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	sent, userId := ch.handlers.validateUserId(writer, req)
	if sent {
		return
	}
	sent = ch.handlers.checkUserAuth(writer, req, userId)
	if sent {
		return
	}
	var chatsQuery models.ChatsQueryValues
	err := ch.parseChatsQuery(query, &chatsQuery)
	if err != nil {
		logger.Errorf(err.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(err.Error()))
		return
	}
	chats, err := ch.ChatsUC.GetUserChatsById(userId, chatsQuery)
	ch.answerChatsQuery(writer, chats, err)
}

func (ch *ChatsHandlers) CreateChatRequest(writer http.ResponseWriter, req *http.Request) {
	var err error
	sent, studentId := ch.handlers.validateStudentId(writer, req)
	if sent {
		return
	}
	sent = ch.handlers.checkUserAuth(writer, req, studentId)
	if sent {
		return
	}
	var chatRequest models.Chat
	err = json.UnmarshalFromReader(req.Body, &chatRequest)
	if err != nil {
		jsonError := fmt.Errorf("error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(jsonError.Error()))
		return
	}
	err = ch.ChatsUC.CreateChatRequest(&chatRequest, studentId)
	ch.answerChat(writer, chatRequest, http.StatusCreated, err)
}

func (ch *ChatsHandlers) ChangeChatStatus(writer http.ResponseWriter, req *http.Request) {
	var err error
	sent, masterId := ch.handlers.validateMasterId(writer, req)
	if sent {
		return
	}
	sent = ch.handlers.checkUserAuth(writer, req, masterId)
	if sent {
		return
	}
	sent, chatId := ch.handlers.validateChatId(writer, req)
	if sent {
		return
	}
	var chatRequest models.Chat
	err = json.UnmarshalFromReader(req.Body, &chatRequest)
	if err != nil {
		jsonError := fmt.Errorf("error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(jsonError.Error()))
		return
	}
	err = ch.ChatsUC.ChangeChatStatus(&chatRequest, masterId, chatId)
	ch.answerChat(writer, chatRequest, http.StatusOK, err)
}

func (ch *ChatsHandlers) GetMessagesByChatId(writer http.ResponseWriter, req *http.Request) {
	sent, chatId := ch.handlers.validateChatIdSimple(writer, req)
	if sent {
		return
	}
	sent = ch.handlers.checkChatAuth(writer, req, chatId)
	if sent {
		return
	}
	messages, err := ch.ChatsUC.GetMessagesByChatId(chatId)
	ch.answerMessages(writer, messages, err)
}

func (ch *ChatsHandlers) answerChatsQuery(writer http.ResponseWriter, chats []models.Chat, err error) {
	sent := ch.handlers.handleErrorBadQueryParameter(writer, err)
	if !sent {
		utils.CreateAnswerChatsJson(writer, http.StatusOK, chats)
	}
}

func (ch *ChatsHandlers) answerMessages(writer http.ResponseWriter, messages []models.Message, err error) {
	sent := ch.handlers.handleError(writer, err)
	if !sent {
		utils.CreateAnswerMessagesJson(writer, http.StatusOK, messages)
	}
}
func (ch *ChatsHandlers) answerChat(writer http.ResponseWriter, chat models.Chat, statusCode int, err error) {
	sent := ch.handlers.handleErrorConflict(writer, err)
	if !sent {
		utils.CreateAnswerChatJson(writer, statusCode, chat)
	}
}
