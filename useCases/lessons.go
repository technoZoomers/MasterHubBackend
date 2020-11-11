package useCases

import (
	"fmt"
	"github.com/google/logger"
	"github.com/shopspring/decimal"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
	"time"
)

type LessonsUC struct {
	useCases      *UseCases
	LessonsRepo   repository.LessonsRepoI
	MastersRepo   repository.MastersRepoI
	lessonsConfig LessonsConfig
}

type LessonsConfig struct {
	zeroTime                    string
	zeroTimeParsed              time.Time
	layoutISODate               string
	layoutISOTime               string
	educationFormatMap          map[int64]string
	educationFormatMapBackwards map[string]int64
}

func (lessonsUC *LessonsUC) validateMaster(masterId int64) (int64, error) { // TODO: throw away redundant functions
	if masterId == lessonsUC.useCases.errorId {
		return lessonsUC.useCases.errorId, &models.BadRequestError{Message: "incorrect master id", RequestId: masterId}
	}
	masterDB := models.MasterDB{
		UserId: masterId,
	}
	err := lessonsUC.MastersRepo.GetMasterByUserId(&masterDB)
	if err != nil {
		return lessonsUC.useCases.errorId, fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	if masterDB.Id == lessonsUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "master doesn't exist", RequestId: masterId}
		logger.Errorf(absenceError.Error())
		return lessonsUC.useCases.errorId, absenceError
	}
	return masterDB.Id, nil
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

func (lessonsUC *LessonsUC) calculateDuration(timeStart string, timeEnd string) (string, time.Duration, error) {
	var duration string
	var durationAsDuration time.Duration
	timeStartAsTime, err := lessonsUC.validateTimeFormat(timeStart)
	if err != nil {
		return duration, durationAsDuration, err
	}
	timeEndAsTime, err := lessonsUC.validateTimeFormat(timeEnd)
	if err != nil {
		return duration, durationAsDuration, err
	}
	durationAsDuration = timeEndAsTime.Sub(timeStartAsTime)
	duration = lessonsUC.formatDuration(durationAsDuration)
	return duration, durationAsDuration, nil
}

func (lessonsUC *LessonsUC) matchLesson(lessonDB *models.LessonDB, lesson *models.Lesson, masterId int64) error {
	lesson.Id = lessonDB.Id
	lesson.MasterId = masterId
	lesson.Date = lessonDB.Date.Format(lessonsUC.lessonsConfig.layoutISODate)
	duration, _, err := lessonsUC.calculateDuration(lessonDB.TimeStart, lessonDB.TimeEnd)
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

func (lessonsUC *LessonsUC) matchLessonToDB(lesson *models.Lesson, lessonDB *models.LessonDB) error {
	dateParsed, err := time.Parse(lessonsUC.lessonsConfig.layoutISODate, lesson.Date)
	if err != nil {
		parseError := fmt.Errorf("couldnt parse date: %s", err.Error())
		logger.Errorf(parseError.Error())
		return parseError
	}
	lessonDB.Date = dateParsed
	duration, durationAsDuration, err := lessonsUC.calculateDuration(lesson.TimeStart, lesson.TimeEnd)
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
	if lesson.Price.Value.IsNegative() {
		formatError := &models.NotAcceptableError{Message: "lesson price must not be negative"}
		logger.Errorf(formatError.Error())
		return formatError
	}
	lessonDB.Price = lesson.Price.Value
	err = lessonsUC.checkLessonStatus(lesson.Status)
	if err != nil {
		return err
	}
	lessonDB.Status = lesson.Status
	return nil
}

func (lessonsUC *LessonsUC) GetMastersLessons() (models.Lessons, error) {
	panic("implement me")
}

func (lessonsUC *LessonsUC) CreateLesson(lesson *models.Lesson, masterId int64) error {
	if lesson.MasterId == lessonsUC.useCases.errorId {
		lesson.MasterId = masterId
	} else if masterId != lesson.MasterId {
		return &models.ForbiddenError{Reason: "master ids doesnt match"}
	}
	masterDBId, err := lessonsUC.validateMaster(lesson.MasterId)
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
	err = lessonsUC.matchLessonToDB(lesson, lessonDB)
	if err != nil {
		return err
	}
	err = lessonsUC.LessonsRepo.InsertLesson(lessonDB)
	if err != nil {
		return fmt.Errorf(lessonsUC.useCases.errorMessages.DbError)
	}
	lesson.Id = lessonDB.Id
	return nil
}

func (lessonsUC *LessonsUC) ChangeLessonInfo(lesson *models.Lesson) error {
	panic("implement me")
}

func (lessonsUC *LessonsUC) GetMastersLessonsRequests() (models.LessonRequests, error) {
	panic("implement me")
}

func (lessonsUC *LessonsUC) CreateLessonRequest(studentId int64, lessonId int64) error {
	panic("implement me")
}

func (lessonsUC *LessonsUC) DeleteLessonRequest(lessonRequest *models.LessonRequest) error {
	panic("implement me")
}
