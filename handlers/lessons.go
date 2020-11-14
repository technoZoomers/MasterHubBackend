package handlers

import (
	"fmt"
	"github.com/google/logger"
	json "github.com/mailru/easyjson"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"net/http"
)

type LessonsHandlers struct {
	handlers  *Handlers
	LessonsUC useCases.LessonsUCInterface
}

func (lh *LessonsHandlers) CreateLesson(writer http.ResponseWriter, req *http.Request) {
	var err error
	sent, masterId := lh.handlers.validateMasterId(writer, req)
	if sent {
		return
	}
	sent = lh.handlers.checkUserAuth(writer, req, masterId)
	if sent {
		return
	}
	var lesson models.Lesson
	err = json.UnmarshalFromReader(req.Body, &lesson)
	if err != nil {
		jsonError := fmt.Errorf("error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(jsonError.Error()))
		return
	}
	err = lh.LessonsUC.CreateLesson(&lesson, masterId)
	lh.answerLesson(writer, lesson, http.StatusCreated, err)
}

func (lh *LessonsHandlers) Get(writer http.ResponseWriter, req *http.Request) {
	var err error
	sent, masterId := lh.handlers.validateMasterId(writer, req)
	if sent {
		return
	}
	lessons, err := lh.LessonsUC.GetMastersLessons(masterId)
	lh.answerLessons(writer, lessons, http.StatusOK, err)
}

func (lh *LessonsHandlers) ChangeLessonInfo(writer http.ResponseWriter, req *http.Request) {
	var err error
	sent, masterId := lh.handlers.validateMasterId(writer, req)
	if sent {
		return
	}
	sent = lh.handlers.checkUserAuth(writer, req, masterId)
	if sent {
		return
	}
	sent, lessonId := lh.handlers.validateLessonId(writer, req)
	if sent {
		return
	}
	var lesson models.Lesson
	err = json.UnmarshalFromReader(req.Body, &lesson)
	if err != nil {
		jsonError := fmt.Errorf("error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(jsonError.Error()))
		return
	}
	err = lh.LessonsUC.ChangeLessonInfo(&lesson, masterId, lessonId)
	lh.answerLesson(writer, lesson, http.StatusOK, err)
}

func (lh *LessonsHandlers) DeleteLesson(writer http.ResponseWriter, req *http.Request) {
	var err error
	sent, masterId := lh.handlers.validateMasterId(writer, req)
	if sent {
		return
	}
	sent = lh.handlers.checkUserAuth(writer, req, masterId)
	if sent {
		return
	}
	sent, lessonId := lh.handlers.validateLessonId(writer, req)
	if sent {
		return
	}
	err = lh.LessonsUC.DeleteMasterLesson(masterId, lessonId)
	lh.answerEmpty(writer, http.StatusOK, err)
}

func (lh *LessonsHandlers) GetLessonStudents(writer http.ResponseWriter, req *http.Request) {
	var err error
	sent, masterId := lh.handlers.validateMasterId(writer, req)
	if sent {
		return
	}
	sent = lh.handlers.checkUserAuth(writer, req, masterId)
	if sent {
		return
	}
	sent, lessonId := lh.handlers.validateLessonId(writer, req)
	if sent {
		return
	}
	students, err := lh.LessonsUC.GetMastersLessonsStudents(masterId, lessonId)
	lh.answerLessonStudents(writer, students, http.StatusOK, err)
}

func (lh *LessonsHandlers) CreateLessonRequest(writer http.ResponseWriter, req *http.Request) {
	var err error
	sent, studentId := lh.handlers.validateStudentId(writer, req)
	if sent {
		return
	}
	sent = lh.handlers.checkUserAuth(writer, req, studentId)
	if sent {
		return
	}
	sent, lessonId := lh.handlers.validateLessonId(writer, req)
	if sent {
		return
	}
	lessonRequest := models.LessonRequest{
		LessonId:  lessonId,
		StudentId: studentId,
	}

	err = lh.LessonsUC.CreateLessonRequest(&lessonRequest)
	lh.answerLessonRequest(writer, lessonRequest, http.StatusCreated, err)
}

func (lh *LessonsHandlers) DeleteLessonRequest(writer http.ResponseWriter, req *http.Request) {
	var err error
	sent, studentId := lh.handlers.validateStudentId(writer, req)
	if sent {
		return
	}
	sent = lh.handlers.checkUserAuth(writer, req, studentId)
	if sent {
		return
	}
	sent, lessonId := lh.handlers.validateLessonId(writer, req)
	if sent {
		return
	}
	err = lh.LessonsUC.DeleteLessonRequest(studentId, lessonId)
	lh.answerEmpty(writer, http.StatusOK, err)
}

func (lh *LessonsHandlers) answerLessonRequest(writer http.ResponseWriter, lessonRequest models.LessonRequest, statusCode int, err error) {
	sent := lh.handlers.handleErrorConflict(writer, err)
	if !sent {
		utils.CreateAnswerLessonRequestJson(writer, statusCode, lessonRequest)
	}
}

func (lh *LessonsHandlers) answerLesson(writer http.ResponseWriter, lesson models.Lesson, statusCode int, err error) {
	sent := lh.handlers.handleNotAcceptableError(writer, err)
	if !sent {
		utils.CreateAnswerLessonJson(writer, statusCode, lesson)
	}
}

func (lh *LessonsHandlers) answerLessonStudents(writer http.ResponseWriter, lessonStudents models.LessonStudents, statusCode int, err error) {
	sent := lh.handlers.handleForbiddenError(writer, err)
	if !sent {
		utils.CreateAnswerLessonStudentsJson(writer, statusCode, lessonStudents)
	}
}

func (lh *LessonsHandlers) answerLessons(writer http.ResponseWriter, lessons models.Lessons, statusCode int, err error) {
	sent := lh.handlers.handleError(writer, err)
	if !sent {
		utils.CreateAnswerLessonsJson(writer, statusCode, lessons)
	}
}

func (lh *LessonsHandlers) answerEmpty(writer http.ResponseWriter, statusCode int, err error) { //TODO: delete redundant functions
	sent := lh.handlers.handleError(writer, err)
	if !sent {
		utils.CreateEmptyBodyAnswer(writer, statusCode)
	}
}
