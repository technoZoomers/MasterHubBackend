package useCases

import (
	"fmt"
	"github.com/google/logger"
	"github.com/shopspring/decimal"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
	"strconv"
	"time"
)

type LessonsUC struct {
	useCases      *UseCases
	LessonsRepo   repository.LessonsRepoI
	MastersRepo   repository.MastersRepoI
	StudentsRepo  repository.StudentsRepoI
	lessonsConfig LessonsConfig
}

type LessonsConfig struct {
	zeroTime                     string
	zeroTimeParsed               time.Time
	layoutISODate                string
	layoutISOTime                string
	educationFormatMap           map[int64]string
	educationFormatMapBackwards  map[string]int64
	lessonRequestStatus          map[int64]string
	lessonStatus                 map[int64]string
	lessonRequestStatusBackwards map[string]int64
}

func (lessonsUC *LessonsUC) validateMaster(masterId int64) (int64, error) { // TODO: throw away redundant functions
	masterId, _, err := lessonsUC.validateMasterWithPrice(masterId)
	return masterId, err
}

func (lessonsUC *LessonsUC) validateMasterWithPrice(masterId int64) (int64, decimal.Decimal, error) { // TODO: throw away redundant functions
	if masterId == lessonsUC.useCases.errorId {
		return lessonsUC.useCases.errorId, decimal.Zero, &models.BadRequestError{Message: "incorrect master id", RequestId: masterId}
	}
	masterDB := models.MasterDB{
		UserId: masterId,
	}
	err := lessonsUC.MastersRepo.GetMasterByUserId(&masterDB)
	if err != nil {
		return lessonsUC.useCases.errorId, masterDB.AveragePrice, fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	if masterDB.Id == lessonsUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "master doesn't exist", RequestId: masterId}
		logger.Errorf(absenceError.Error())
		return lessonsUC.useCases.errorId, masterDB.AveragePrice, absenceError
	}
	return masterDB.Id, masterDB.AveragePrice, nil
}

func (lessonsUC *LessonsUC) validateStudent(studentId int64) (int64, error) {
	if studentId == lessonsUC.useCases.errorId {
		return lessonsUC.useCases.errorId, &models.BadRequestError{Message: "incorrect student id", RequestId: studentId}
	}
	studentDB := models.StudentDB{
		UserId: studentId,
	}
	err := lessonsUC.StudentsRepo.GetStudentByUserId(&studentDB)
	if err != nil {
		return lessonsUC.useCases.errorId, fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	if studentDB.Id == lessonsUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "student doesn't exist", RequestId: studentId}
		logger.Errorf(absenceError.Error())
		return lessonsUC.useCases.errorId, absenceError
	}
	return studentDB.Id, nil
}

func (lessonsUC *LessonsUC) validateMasterLesson(lessonId int64, masterId int64, lessonDB *models.LessonDB) error {
	if lessonId == lessonsUC.useCases.errorId {
		return &models.BadRequestError{Message: "incorrect lesson id", RequestId: lessonId}
	}
	err := lessonsUC.LessonsRepo.GetLessonByIdAndMasterId(lessonDB, lessonId, masterId)
	if err != nil {
		return fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	if lessonDB.Id == lessonsUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "lesson doesn't exist", RequestId: lessonId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	return nil
}

func (lessonsUC *LessonsUC) validateLesson(lessonId int64, lessonDB *models.LessonDB) error {
	if lessonId == lessonsUC.useCases.errorId {
		return &models.BadRequestError{Message: "incorrect lesson id", RequestId: lessonId}
	}
	err := lessonsUC.LessonsRepo.GetLessonById(lessonDB, lessonId)
	if err != nil {
		return fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	if lessonDB.Id == lessonsUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "lesson doesn't exist", RequestId: lessonId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	return nil
}

func (lessonsUC *LessonsUC) validateLessonWithMasterUserId(lessonId int64, lessonDB *models.LessonDB) error {
	if lessonId == lessonsUC.useCases.errorId {
		return &models.BadRequestError{Message: "incorrect lesson id", RequestId: lessonId}
	}
	err := lessonsUC.LessonsRepo.GetLessonByIdWithMasterUserId(lessonDB, lessonId)
	if err != nil {
		return fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	if lessonDB.Id == lessonsUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "lesson doesn't exist", RequestId: lessonId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	return nil
}

func (lessonsUC *LessonsUC) validateLessonRequestNotExist(lessonId int64, studentDBId int64) error {
	var lessonRequestDBExist models.LessonStudentDB
	err := lessonsUC.LessonsRepo.GetLessonRequestByStudentIdAndLessonId(&lessonRequestDBExist, studentDBId, lessonId)
	if err != nil {
		return fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	if lessonRequestDBExist.LessonId != lessonsUC.useCases.errorId {
		return &models.ConflictError{Message: "lesson request already exists"}
	}
	return nil
}

func (lessonsUC *LessonsUC) validateLessonRequestExistWithStudentDBId(lessonId int64, studentDBId int64) error {
	var lessonRequestDBExist models.LessonStudentDB
	err := lessonsUC.LessonsRepo.GetLessonRequestByStudentIdAndLessonId(&lessonRequestDBExist, studentDBId, lessonId)
	if err != nil {
		return fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	if lessonRequestDBExist.LessonId == lessonsUC.useCases.errorId {
		return &models.BadRequestError{Message: "lesson request doesnt exist"}
	}
	return nil
}

func (lessonsUC *LessonsUC) validateLessonRequestExistWithStudentUserId(lessonRequest *models.LessonStudentDB, lessonId int64, studentDBId int64) error {
	err := lessonsUC.LessonsRepo.GetLessonRequestByStudentUserIdAndLessonId(lessonRequest, studentDBId, lessonId)
	if err != nil {
		return fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	if lessonRequest.LessonId == lessonsUC.useCases.errorId {
		return &models.BadRequestError{Message: "lesson request doesnt exist"}
	}
	return nil
}

func (lessonsUC *LessonsUC) matchEducationFormatToDB(format string) (int64, error) {
	formatInt := lessonsUC.lessonsConfig.educationFormatMapBackwards[format]
	if formatInt == 0 {
		formatError := &models.NotAcceptableError{Message: "wrong education format type"}
		logger.Errorf(formatError.Error())
		return formatInt, formatError
	}
	return formatInt, nil
}

func (lessonsUC *LessonsUC) matchEducationFormat(format int64) (string, error) {
	formatString := lessonsUC.lessonsConfig.educationFormatMap[format]
	if formatString == "" {
		formatError := fmt.Errorf("wrong education format type")
		logger.Errorf(formatError.Error())
		return formatString, formatError
	}
	return formatString, nil
}

func (lessonsUC *LessonsUC) matchPrice(lesson *models.Lesson, price decimal.Decimal) error {
	lesson.Price.Value = price
	lesson.Price.Currency = "rub" //TODO: change to different currencies
	return nil
}

func (lessonsUC *LessonsUC) checkLessonStatus(status int64) error {
	if !(status <= 3 && status >= 1) {
		formatError := &models.NotAcceptableError{Message: "wrong lesson status"}
		logger.Errorf(formatError.Error())
		return formatError
	}
	return nil
}

func (lessonsUC *LessonsUC) checkLessonStatusNew(status int64) error {
	if status != 1 {
		formatError := &models.NotAcceptableError{Message: "wrong lesson status"}
		logger.Errorf(formatError.Error())
		return formatError
	}
	return nil
}

func (lessonsUC *LessonsUC) formatDuration(d time.Duration) string { // TODO: make as helper
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func (lessonsUC *LessonsUC) validateTimeFormat(timeToParse string) (time.Time, error) {
	timeParsed, err := time.Parse(lessonsUC.lessonsConfig.layoutISOTime, timeToParse)
	if err != nil {
		parseError := &models.NotAcceptableError{Message: fmt.Sprintf("couldnt parse time start: %s", err.Error())}
		logger.Errorf(parseError.Error())
		return timeParsed, parseError
	}
	return timeParsed, nil
}

func (lessonsUC *LessonsUC) validateTimeFormatInsert(timeToParse string, dateEqual bool) (time.Time, error) {
	timeParsed, err := lessonsUC.validateTimeFormat(timeToParse)
	if err != nil {
		return timeParsed, err
	}
	if dateEqual {
		if timeParsed.Hour() > time.Now().Hour() {
			return timeParsed, nil
		} else if timeParsed.Hour() == time.Now().Hour() {
			if timeParsed.Minute() > time.Now().Minute() {
				return timeParsed, nil
			} else if timeParsed.Minute() == time.Now().Minute() {
				if timeParsed.Second() > time.Now().Second() {
					return timeParsed, nil
				}
			}
		}
		formatError := &models.NotAcceptableError{Message: "can't set past time"}
		logger.Errorf(formatError.Error())
		return timeParsed, formatError
	}
	return timeParsed, nil
}

func (lessonsUC *LessonsUC) calculateDuration(timeStart string, timeEnd string, insert bool, dateEqual bool) (string, time.Duration, error) {
	var duration string
	var durationAsDuration time.Duration
	var timeStartAsTime time.Time
	var timeEndAsTime time.Time
	var err error
	if insert {
		timeStartAsTime, err = lessonsUC.validateTimeFormatInsert(timeStart, dateEqual)
		if err != nil {
			return duration, durationAsDuration, err
		}
		timeEndAsTime, err = lessonsUC.validateTimeFormatInsert(timeEnd, dateEqual)
		if err != nil {
			return duration, durationAsDuration, err
		}
	} else {
		timeStartAsTime, err = lessonsUC.validateTimeFormat(timeStart)
		if err != nil {
			return duration, durationAsDuration, err
		}
		timeEndAsTime, err = lessonsUC.validateTimeFormat(timeEnd)
		if err != nil {
			return duration, durationAsDuration, err
		}
	}

	durationAsDuration = timeEndAsTime.Sub(timeStartAsTime)
	duration = lessonsUC.formatDuration(durationAsDuration)
	return duration, durationAsDuration, nil
}

func (lessonsUC *LessonsUC) matchLesson(lessonDB *models.LessonDB, lesson *models.Lesson, masterId int64) error {
	lesson.Id = lessonDB.Id
	lesson.MasterId = masterId
	lesson.Date = lessonDB.Date.Format(lessonsUC.lessonsConfig.layoutISODate)
	duration, _, err := lessonsUC.calculateDuration(lessonDB.TimeStart, lessonDB.TimeEnd, false, false)
	if err != nil {
		return err
	}
	lesson.TimeStart = lessonDB.TimeStart
	lesson.TimeEnd = lessonDB.TimeEnd
	lesson.Duration = duration

	edFormat, err := lessonsUC.matchEducationFormat(lessonDB.EducationFormat)
	if err != nil {
		return fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	lesson.EducationFormat = edFormat
	err = lessonsUC.matchPrice(lesson, lessonDB.Price)
	if err != nil {
		return err
	}
	err = lessonsUC.checkLessonStatus(lessonDB.Status)
	if err != nil {
		return fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	lesson.Status = lessonDB.Status
	return nil
}

func (lessonsUC *LessonsUC) checkLessonDate(lessonDate time.Time) (bool, error) {
	lessonYear, lessonMonth, lessonDay := lessonDate.Date()
	nowYear, nowMonth, nowDay := time.Now().Date()
	if lessonYear > nowYear {
		return false, nil
	} else if lessonYear == nowYear {
		if lessonMonth > nowMonth {
			return false, nil
		} else if lessonMonth == nowMonth {
			if lessonDay > nowDay {
				return false, nil
			} else if lessonDay == nowDay {
				return true, nil
			}
		}
	}
	formatError := &models.NotAcceptableError{Message: "can't set past date"}
	logger.Errorf(formatError.Error())
	return false, formatError
}

func (lessonsUC *LessonsUC) matchLessonToDB(lesson *models.Lesson, lessonDB *models.LessonDB, masterPrice decimal.Decimal) error {
	dateParsed, err := time.Parse(lessonsUC.lessonsConfig.layoutISODate, lesson.Date)
	if err != nil {
		parseError := fmt.Errorf("couldnt parse date: %s", err.Error())
		logger.Errorf(parseError.Error())
		return parseError
	}
	dateEqual, err := lessonsUC.checkLessonDate(dateParsed)
	if err != nil {
		return err
	}
	lessonDB.Date = dateParsed
	duration, durationAsDuration, err := lessonsUC.calculateDuration(lesson.TimeStart, lesson.TimeEnd, true, dateEqual)
	if err != nil {
		return err
	}
	lessonDB.TimeStart = lesson.TimeStart
	lessonDB.TimeEnd = lesson.TimeEnd
	lesson.Duration = duration

	if durationAsDuration <= 0 {
		formatError := &models.NotAcceptableError{Message: "lesson duration must be positive"}
		logger.Errorf(formatError.Error())
		return formatError
	}

	edFormat, err := lessonsUC.matchEducationFormatToDB(lesson.EducationFormat)
	if err != nil {
		return err
	}
	lessonDB.EducationFormat = edFormat
	masterPriceFloat, _ := masterPrice.Float64()
	if !lesson.Price.Value.IsZero() {
		formatError := &models.NotAcceptableError{Message: "lesson price can't be set"}
		logger.Errorf(formatError.Error())
		return formatError
	}
	lesson.Price.Value = decimal.NewFromFloat(durationAsDuration.Minutes() / 60 * masterPriceFloat).Truncate(2)
	lessonDB.Price = lesson.Price.Value
	err = lessonsUC.checkLessonStatus(lesson.Status)
	if err != nil {
		return err
	}
	lessonDB.Status = lesson.Status
	return nil
}

func (lessonsUC *LessonsUC) matchLessonToDBUpdate(lesson *models.Lesson, lessonDB *models.LessonDB) error {
	var dateEqual bool
	if lesson.Date != "" {
		dateParsed, err := time.Parse(lessonsUC.lessonsConfig.layoutISODate, lesson.Date)
		if err != nil {
			parseError := fmt.Errorf("couldnt parse date: %s", err.Error())
			logger.Errorf(parseError.Error())
			return parseError
		}
		dateEqual, err = lessonsUC.checkLessonDate(dateParsed)
		if err != nil {
			return err
		}
		lessonDB.Date = dateParsed
	} else {
		lesson.Date = lessonDB.Date.Format(lessonsUC.lessonsConfig.layoutISODate)
	}
	if lesson.TimeEnd == "" {
		lesson.TimeEnd = lessonDB.TimeEnd
	}
	if lesson.TimeStart == "" {
		lesson.TimeStart = lessonDB.TimeStart
	}
	duration, durationAsDuration, err := lessonsUC.calculateDuration(lesson.TimeStart, lesson.TimeEnd, true, dateEqual)
	if err != nil {
		return err
	}
	lessonDB.TimeStart = lesson.TimeStart
	lessonDB.TimeEnd = lesson.TimeEnd
	lesson.Duration = duration

	if durationAsDuration <= 0 {
		formatError := &models.NotAcceptableError{Message: "lesson duration must be positive"}
		logger.Errorf(formatError.Error())
		return formatError
	}
	if lesson.EducationFormat == "" {
		lesson.EducationFormat, _ = lessonsUC.matchEducationFormat(lessonDB.EducationFormat)
	}
	edFormat, err := lessonsUC.matchEducationFormatToDB(lesson.EducationFormat)
	if err != nil {
		return err
	}
	lessonDB.EducationFormat = edFormat
	if lesson.Price.Currency == "" {
		_ = lessonsUC.matchPrice(lesson, lessonDB.Price)
	}
	if lessonDB.Price != lesson.Price.Value {
		formatError := &models.NotAcceptableError{Message: "lesson price can't be changed"}
		logger.Errorf(formatError.Error())
		return formatError
	}
	if lesson.Status == lessonsUC.useCases.errorId {
		lesson.Status = lessonDB.Status
	} else {
		err = lessonsUC.checkLessonStatus(lesson.Status)
		if err != nil {
			return err
		}
		lessonStudents, err := lessonsUC.LessonsRepo.GetLessonStudents(lesson.Id)
		if err != nil {
			return fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
		}
		if len(lessonStudents) == 0 {
			formatError := &models.NotAcceptableError{Message: "can't book a lesson with no students"}
			logger.Errorf(formatError.Error())
			return formatError
		}
		lessonDB.Status = lesson.Status
	}
	return nil
}

func (lessonsUC *LessonsUC) matchLessonRequest(lessonRequestDB *models.LessonStudentDB, lessonRequest *models.LessonRequest, studentId int64) error {
	lessonRequest.LessonId = lessonRequestDB.LessonId
	lessonRequest.StudentId = studentId
	lessonRequest.Status = lessonRequestDB.Status
	return nil
}

func (lessonsUC *LessonsUC) CreateLesson(lesson *models.Lesson, masterId int64) error {
	if lesson.MasterId == lessonsUC.useCases.errorId {
		lesson.MasterId = masterId
	} else if masterId != lesson.MasterId {
		return &models.ForbiddenError{Reason: "master ids doesnt match"}
	}
	masterDBId, masterPrice, err := lessonsUC.validateMasterWithPrice(lesson.MasterId)
	if err != nil {
		return err
	}
	lessonDB := &models.LessonDB{
		MasterId: masterDBId,
	}
	err = lessonsUC.checkLessonStatusNew(lesson.Status)
	if err != nil {
		return err
	}
	err = lessonsUC.matchLessonToDB(lesson, lessonDB, masterPrice)
	if err != nil {
		return err
	}
	lessonsIds, err := lessonsUC.LessonsRepo.CheckLessonTimeRange(lessonDB)
	if err != nil {
		return fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	if len(lessonsIds) != 0 {
		formatError := &models.ConflictError{Message: "this lesson time is already taken", ExistingContent: strconv.FormatInt(lessonsIds[0], 10)}
		logger.Errorf(formatError.Error())
		return formatError
	}
	err = lessonsUC.LessonsRepo.InsertLesson(lessonDB)
	if err != nil {
		return fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	lesson.Id = lessonDB.Id
	return nil
}

func (lessonsUC *LessonsUC) ChangeLessonInfo(lesson *models.Lesson, masterId int64, lessonId int64) error {
	if lesson.MasterId == lessonsUC.useCases.errorId {
		lesson.MasterId = masterId
	} else if masterId != lesson.MasterId {
		return &models.ForbiddenError{Reason: "master ids doesnt match"}
	}
	masterDBId, err := lessonsUC.validateMaster(lesson.MasterId)
	if err != nil {
		return err
	}
	if lesson.Id == lessonsUC.useCases.errorId {
		lesson.Id = lessonId
	} else if lessonId != lesson.Id {
		return &models.ForbiddenError{Reason: "lesson ids doesnt match"}
	}
	var lessonDB models.LessonDB
	err = lessonsUC.validateMasterLesson(lesson.Id, masterDBId, &lessonDB)
	if err != nil {
		return err
	}
	err = lessonsUC.matchLessonToDBUpdate(lesson, &lessonDB)
	if err != nil {
		return err
	}
	lessonsIds, err := lessonsUC.LessonsRepo.CheckLessonTimeRange(&lessonDB)
	if err != nil {
		return fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	if !(len(lessonsIds) == 1 && lessonsIds[0] == lesson.Id) && len(lessonsIds) != 0 {
		formatError := &models.ConflictError{Message: "this lesson time is already taken", ExistingContent: strconv.FormatInt(lessonsIds[0], 10)}
		logger.Errorf(formatError.Error())
		return formatError
	}
	err = lessonsUC.LessonsRepo.UpdateLessonByIdAndMasterId(&lessonDB)
	if err != nil {
		return fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	return nil
}

func (lessonsUC *LessonsUC) ChangeLessonRequest(lessonRequest *models.LessonRequest, masterId int64, lessonId int64) error {
	masterDBId, err := lessonsUC.validateMaster(masterId)
	if err != nil {
		return err
	}
	if lessonRequest.LessonId == lessonsUC.useCases.errorId {
		lessonRequest.LessonId = lessonId
	} else if lessonId != lessonRequest.LessonId {
		return &models.ForbiddenError{Reason: "lesson ids doesnt match"}
	}
	var lessonDB models.LessonDB
	err = lessonsUC.validateMasterLesson(lessonRequest.LessonId, masterDBId, &lessonDB)
	if err != nil {
		return err
	}
	if lessonRequest.StudentId == lessonsUC.useCases.errorId {
		formatError := &models.NotAcceptableError{Message: "student id is not present"}
		logger.Errorf(formatError.Error())
		return formatError
	}
	if lessonsUC.lessonsConfig.lessonRequestStatus[lessonRequest.Status] == "" ||
		lessonsUC.lessonsConfig.lessonRequestStatus[lessonRequest.Status] == "unseen" {
		formatError := &models.NotAcceptableError{Message: "choose valid status: 2 or 3"}
		logger.Errorf(formatError.Error())
		return formatError
	}
	var lessonRequestDB models.LessonStudentDB
	err = lessonsUC.validateLessonRequestExistWithStudentUserId(&lessonRequestDB, lessonId, lessonRequest.StudentId)
	if err != nil {
		return err
	}
	lessonRequestDB.Status = lessonRequest.Status
	err = lessonsUC.LessonsRepo.UpdateLessonRequest(&lessonRequestDB)
	if err != nil {
		return fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	return nil
}
func (lessonsUC *LessonsUC) GetMastersLessonsRequests(masterId int64, query models.LessonsQueryValues) (models.LessonRequests, error) {
	lessonsRequests := make([]models.LessonRequest, 0)
	var queryDB models.LessonsQueryValuesDB
	err := lessonsUC.matchLessonsRequestQuery(&query, &queryDB)
	if err != nil {
		return lessonsRequests, err
	}
	masterDBId, err := lessonsUC.validateMaster(masterId)
	if err != nil {
		return lessonsRequests, err
	}
	lessonsRequestsDB, err := lessonsUC.LessonsRepo.GetMastersLessonsRequests(masterDBId, queryDB)
	if err != nil {
		return lessonsRequests, fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	for _, lessonRequestDB := range lessonsRequestsDB {
		var lessonRequest models.LessonRequest
		err = lessonsUC.matchLessonRequest(&lessonRequestDB, &lessonRequest, lessonRequestDB.StudentId)
		if err != nil {
			return lessonsRequests, err
		}
		lessonsRequests = append(lessonsRequests, lessonRequest)
	}
	return lessonsRequests, nil
}

func (lessonsUC *LessonsUC) matchLessonsRequestQuery(query *models.LessonsQueryValues, queryDB *models.LessonsQueryValuesDB) error {
	if query.Status != 0 && lessonsUC.lessonsConfig.lessonRequestStatus[query.Status] == "" {
		requestError := &models.BadRequestError{Message: "choose valid status: 1, 2 or 3"}
		logger.Errorf(requestError.Error())
		return requestError
	}
	queryDB.Status = query.Status
	return nil
}

func (lessonsUC *LessonsUC) matchLessonsQuery(query *models.LessonsQueryValues, queryDB *models.LessonsQueryValuesDB) error {
	if query.Status != 0 && lessonsUC.lessonsConfig.lessonStatus[query.Status] == "" {
		requestError := &models.BadRequestError{Message: "choose valid status: 1, 2 or 3"}
		logger.Errorf(requestError.Error())
		return requestError
	}
	queryDB.Status = query.Status
	return nil
}

func (lessonsUC *LessonsUC) GetStudentsLessonsRequests(studentId int64, query models.LessonsQueryValues) (models.LessonRequests, error) {
	lessonsRequests := make([]models.LessonRequest, 0)
	var queryDB models.LessonsQueryValuesDB
	err := lessonsUC.matchLessonsRequestQuery(&query, &queryDB)
	if err != nil {
		return lessonsRequests, err
	}
	studentDBId, err := lessonsUC.validateStudent(studentId)
	if err != nil {
		return lessonsRequests, err
	}
	lessonsRequestsDB, err := lessonsUC.LessonsRepo.GetStudentsLessonsRequests(studentDBId, queryDB)
	if err != nil {
		return lessonsRequests, fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	for _, lessonRequestDB := range lessonsRequestsDB {
		var lessonRequest models.LessonRequest
		err = lessonsUC.matchLessonRequest(&lessonRequestDB, &lessonRequest, studentId)
		if err != nil {
			return lessonsRequests, err
		}
		lessonsRequests = append(lessonsRequests, lessonRequest)
	}
	return lessonsRequests, nil
}

func (lessonsUC *LessonsUC) CreateLessonRequest(lessonRequest *models.LessonRequest) error {
	studentDBId, err := lessonsUC.validateStudent(lessonRequest.StudentId)
	if err != nil {
		return err
	}
	var lessonDB models.LessonDB
	err = lessonsUC.validateLesson(lessonRequest.LessonId, &lessonDB)
	if err != nil {
		return err
	}
	err = lessonsUC.validateLessonRequestNotExist(lessonRequest.LessonId, studentDBId)
	if err != nil {
		return err
	}

	lessonRequestDB := models.LessonStudentDB{
		StudentId: studentDBId,
		LessonId:  lessonRequest.LessonId,
		Status:    lessonsUC.lessonsConfig.lessonRequestStatusBackwards["unseen"],
	}
	err = lessonsUC.LessonsRepo.InsertLessonRequest(&lessonRequestDB)
	if err != nil {
		return fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	lessonRequest.Status = lessonRequestDB.Status
	return nil
}

func (lessonsUC *LessonsUC) DeleteLessonRequest(studentId int64, lessonId int64) error {
	studentDBId, err := lessonsUC.validateStudent(studentId)
	if err != nil {
		return err
	}
	var lessonDB models.LessonDB
	err = lessonsUC.validateLesson(lessonId, &lessonDB)
	if err != nil {
		return err
	}
	err = lessonsUC.validateLessonRequestExistWithStudentDBId(lessonId, studentDBId)
	if err != nil {
		return err
	}
	err = lessonsUC.LessonsRepo.DeleteLessonRequestByStudentIdAndLessonId(studentDBId, lessonId)
	if err != nil {
		return fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	return nil
}

func (lessonsUC *LessonsUC) DeleteMasterLesson(masterId int64, lessonId int64) error {
	masterDBId, err := lessonsUC.validateMaster(masterId)
	if err != nil {
		return err
	}
	var lessonDB models.LessonDB
	err = lessonsUC.validateMasterLesson(lessonId, masterDBId, &lessonDB)
	if err != nil {
		return err
	}
	err = lessonsUC.LessonsRepo.DeleteLessonById(lessonId)
	if err != nil {
		return fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	return nil
}

func (lessonsUC *LessonsUC) GetMastersLessons(masterId int64, query models.LessonsQueryValues) (models.Lessons, error) {
	lessons := make([]models.Lesson, 0)
	var queryDB models.LessonsQueryValuesDB
	err := lessonsUC.matchLessonsQuery(&query, &queryDB)
	if err != nil {
		return lessons, err
	}
	masterDBId, err := lessonsUC.validateMaster(masterId)
	if err != nil {
		return lessons, err
	}
	lessonsDB, err := lessonsUC.LessonsRepo.GetMastersLessons(masterDBId, queryDB)
	if err != nil {
		return lessons, fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	for _, lessonDB := range lessonsDB {
		var lesson models.Lesson
		err = lessonsUC.matchLesson(&lessonDB, &lesson, masterId)
		if err != nil {
			return lessons, err
		}
		lessons = append(lessons, lesson)
	}
	return lessons, nil
}

func (lessonsUC *LessonsUC) GetStudentsLessons(studentId int64, query models.LessonsQueryValues) (models.Lessons, error) {
	lessons := make([]models.Lesson, 0)
	var queryDB models.LessonsQueryValuesDB
	err := lessonsUC.matchLessonsQuery(&query, &queryDB)
	if err != nil {
		return lessons, err
	}
	studentDBId, err := lessonsUC.validateStudent(studentId)
	if err != nil {
		return lessons, err
	}
	lessonsDB, err := lessonsUC.LessonsRepo.GetStudentsLessons(studentDBId, queryDB)
	if err != nil {
		return lessons, fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	for _, lessonDB := range lessonsDB {
		var lesson models.Lesson
		err = lessonsUC.matchLesson(&lessonDB, &lesson, lessonDB.MasterId)
		if err != nil {
			return lessons, err
		}
		lessons = append(lessons, lesson)
	}
	return lessons, nil
}

func (lessonsUC *LessonsUC) GetMastersLessonsStudents(masterId int64, lessonId int64) (models.LessonStudents, error) {
	lessonStudents := make([]int64, 0)
	masterDBId, err := lessonsUC.validateMaster(masterId)
	if err != nil {
		return lessonStudents, err
	}
	var lessonDB models.LessonDB
	err = lessonsUC.validateMasterLesson(lessonId, masterDBId, &lessonDB)
	if err != nil {
		return lessonStudents, err
	}
	lessonStudents, err = lessonsUC.LessonsRepo.GetLessonStudents(lessonId)
	if err != nil {
		return lessonStudents, fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	return lessonStudents, nil
}
func (lessonsUC *LessonsUC) GetLessonById(lesson *models.Lesson) error {
	var lessonDB models.LessonDB
	err := lessonsUC.validateLessonWithMasterUserId(lesson.Id, &lessonDB)
	if err != nil {
		return err
	}
	err = lessonsUC.matchLesson(&lessonDB, lesson, lessonDB.MasterId)
	if err != nil {
		return err
	}
	return nil
}
