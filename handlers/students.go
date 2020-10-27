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

type StudentsHandlers struct {
	handlers   *Handlers
	StudentsUC useCases.StudentsUCInterface
}

func (sh *StudentsHandlers) Register(writer http.ResponseWriter, req *http.Request) {
	var newStudent models.StudentFull
	err := json.UnmarshalFromReader(req.Body, &newStudent)
	if err != nil {
		jsonError := fmt.Errorf("error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(jsonError.Error()))
		return
	}
	err = sh.StudentsUC.Register(&newStudent)
	sh.answerStudentFull(writer, newStudent, err)

}
func (sh *StudentsHandlers) GetStudentById(writer http.ResponseWriter, req *http.Request) {
	sent, studentId := sh.handlers.validateStudentId(writer, req)
	if sent {
		return
	}
	var student models.Student
	student.UserId = studentId
	err := sh.StudentsUC.GetStudentById(&student)
	sh.answerStudent(writer, student, err)
}

func (sh *StudentsHandlers) ChangeStudentData(writer http.ResponseWriter, req *http.Request) {
	sent, studentId := sh.handlers.validateStudentId(writer, req)
	if sent {
		return
	}
	var student models.Student
	err := json.UnmarshalFromReader(req.Body, &student)
	if err != nil {
		jsonError := fmt.Errorf("error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(jsonError.Error()))
		return
	}
	if studentId != student.UserId {
		paramError := fmt.Errorf("wrong student id parameter")
		logger.Errorf(paramError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(paramError.Error()))
		return
	}
	err = sh.StudentsUC.ChangeStudentData(&student)
	sh.answerStudent(writer, student, err)
}


func (sh *StudentsHandlers) answerStudent(writer http.ResponseWriter, student models.Student,  err error) {
	sent := sh.handlers.handleErrorConflict(writer, err)
	if !sent {
		utils.CreateAnswerStudentJson(writer, http.StatusOK, student)
	}
}

func (sh *StudentsHandlers) answerStudentFull(writer http.ResponseWriter, studentFull models.StudentFull,  err error) {
	sent := sh.handlers.handleErrorConflict(writer, err)
	if !sent {
		utils.CreateAnswerStudentFullJson(writer, http.StatusCreated, studentFull)
	}
}
