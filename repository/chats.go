package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
)

type ChatsRepo struct {
	repository *Repository
	userMap map[string]int64
}

func (chatsRepo *ChatsRepo) addWhereParamsToQuery (query models.ChatsQueryValuesDB, selectQuery string,
	queryValues *[]interface{}, queryCount int) (string, int, error) {
	queryCount++
	if query.User == chatsRepo.userMap["master"] {
		selectQuery += fmt.Sprintf(" user_id_master = $%d", queryCount)
	} else if query.User == chatsRepo.userMap["student"] {
		selectQuery += fmt.Sprintf(" user_id_student = $%d", queryCount)
	} else {
		return selectQuery, queryCount, fmt.Errorf("wrong user type")
	}
	*queryValues = append(*queryValues, query.UserId)

	if query.Type != 0 {
		queryCount++
		selectQuery += fmt.Sprintf(" AND type = $%d", queryCount)
		*queryValues = append(*queryValues, query.Type)
	}
	return selectQuery, queryCount, nil
}

func (chatsRepo *ChatsRepo) GetChatsByUserId(query models.ChatsQueryValuesDB) ([]models.ChatDB, error) {
	var dbError error
	chats := make([]models.ChatDB, 0)
	transaction, err := chatsRepo.repository.startTransaction()
	if err != nil {
		return chats, err
	}

	var queryValues []interface{}
	queryCount := 0
	var selectQuery string
	if query.Limit == 0 && query.Offset == 0 {
		selectQuery = "SELECT chats.id, chats.user_id_master, user_id_student, chats.type, chats.created FROM chats INNER JOIN messages m on chats.id = m.chat_id WHERE"
		selectQuery, queryCount, err = chatsRepo.addWhereParamsToQuery(query, selectQuery, &queryValues, queryCount)
		if err != nil {
			return chats, err
		}
		selectQuery += " GROUP BY (chats.id) ORDER BY MAX(m.created) DESC;"
	} else {
		selectQuery = "SELECT id, user_id_master, user_id_student, type, created FROM (SELECT row_number() over (ORDER BY MAX(m.created) DESC) " +
			"as select_id, chats.id, chats.user_id_master, user_id_student, chats.type, chats.created FROM chats " +
			"INNER JOIN messages m on chats.id = m.chat_id WHERE"
		selectQuery, queryCount, err = chatsRepo.addWhereParamsToQuery(query, selectQuery, &queryValues, queryCount)
		selectQuery+= " GROUP BY chats.id"
		if err != nil {
			return chats, err
		}
		selectQuery += ") as i"
		if query.Limit == 0 {
			queryCount++
			selectQuery += fmt.Sprintf(" WHERE i.select_id > $%d", queryCount)
			queryValues = append(queryValues, query.Offset)
		}else {
			queryCount++
			selectQuery += fmt.Sprintf(" WHERE i.select_id BETWEEN $%d", queryCount)
			queryCount++
			selectQuery += fmt.Sprintf(" AND $%d", queryCount)
			queryValues = append(queryValues, query.Offset+1, query.Offset+query.Limit)
		}
	}
	rows, err := transaction.Query(selectQuery, queryValues...)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve chats: %v", err.Error())
		logger.Errorf(dbError.Error())
		return chats, dbError
	}
	for rows.Next() {
		var chatFound models.ChatDB
		err = rows.Scan(&chatFound.Id, &chatFound.MasterId, &chatFound.StudentId, &chatFound.Type, &chatFound.Created)
		if err != nil {
			dbError = fmt.Errorf("failed to retrieve chat: %v", err)
			logger.Errorf(dbError.Error())
			return chats, dbError
		}
		chats = append(chats, chatFound)
	}
	err = chatsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return chats, err
	}
	return chats, nil
}

func  (chatsRepo *ChatsRepo) InsertChatRequest(chat *models.ChatDB) error {
	var dbError error
	transaction, err := chatsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("INSERT INTO chats (user_id_master, user_id_student, type, created) VALUES ($1, $2, $3, $4) returning id",
		chat.MasterId, chat.StudentId, chat.Type, chat.Created)
	err = row.Scan(&chat.Id)
	if err != nil {
		dbError = fmt.Errorf("failed to insert chat: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := chatsRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = chatsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (chatsRepo *ChatsRepo) GetChatByStudentIdAndMasterId(chat *models.ChatDB) error {
	var dbError error
	transaction, err := chatsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT * FROM chats WHERE user_id_student=$1 AND user_id_master=$2", chat.StudentId, chat.MasterId)
	err = row.Scan(&chat.Id, &chat.MasterId, &chat.StudentId, &chat.Type, &chat.Created)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve chat: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = chatsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (chatsRepo *ChatsRepo) GetChatById(chat *models.ChatDB, chatId int64) error { //TODO: refactor!!!
	var dbError error
	transaction, err := chatsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT * FROM chats WHERE id=$1", chatId)
	err = row.Scan(&chat.Id, &chat.MasterId, &chat.StudentId, &chat.Type, &chat.Created)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve chat: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = chatsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (chatsRepo *ChatsRepo) GetChatByIdAndMasterOrStudentId(chat *models.ChatDB, chatId int64, userId int64) error { //TODO: refactor!!!
	var dbError error
	transaction, err := chatsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT * FROM chats WHERE id = $1 AND (user_id_master=$2 OR user_id_student = $2)", chatId, userId)
	err = row.Scan(&chat.Id, &chat.MasterId, &chat.StudentId, &chat.Type, &chat.Created)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve chat: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = chatsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}



func (chatsRepo *ChatsRepo) ChangeChatType(chat *models.ChatDB) error {
	var dbError error
	transaction, err := chatsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	_, err = transaction.Exec("UPDATE chats SET type = $1 WHERE id=$2", chat.Type, chat.Id)
	if err != nil {
		dbError = fmt.Errorf("failed to update chat type: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := chatsRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = chatsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}
func(chatsRepo *ChatsRepo) GetMessagesByChatId(chatId int64) ([]models.MessageDB, error) {
	var dbError error
	messages := make([]models.MessageDB, 0)
	transaction, err := chatsRepo.repository.startTransaction()
	if err != nil {
		return messages, err
	}
	rows, err := transaction.Query(`SELECT * FROM messages WHERE chat_id = $1`, chatId)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve messages: %v", err.Error())
		logger.Errorf(dbError.Error())
		return messages, dbError
	}
	for rows.Next() {
		var userId sql.NullInt64
		var messageFound models.MessageDB
		err = rows.Scan(&messageFound.Id, &messageFound.Info, &userId, &messageFound.ChatId, &messageFound.Text, &messageFound.Created)
		if err != nil {
			dbError = fmt.Errorf("failed to retrieve message: %v", err)
			logger.Errorf(dbError.Error())
			return messages, dbError
		}
		messageFound.UserId = checkNullValueInt64(userId)
		messages = append(messages, messageFound)
	}
	err = chatsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return messages, err
	}
	return messages, nil
}

func (chatsRepo *ChatsRepo) InsertMessage(message *models.MessageDB) error {
	var dbError error
	transaction, err := chatsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("INSERT INTO messages (info, user_id, chat_id, text, created) values ($1, $2, $3, $4, $5) returning id",
		message.Info, message.UserId, message.ChatId, message.Text, message.Created)
	err = row.Scan(&message.Id)
	if err != nil {
		dbError = fmt.Errorf("failed to insert chat: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := chatsRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = chatsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}