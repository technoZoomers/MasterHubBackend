package useCases

import (
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
	gomail "gopkg.in/mail.v2"
	"time"
)

type LessonNotificationsUC struct {
	useCases                   *UseCases
	websocketUC                WebsocketsUCInterface
	MastersRepo                repository.MastersRepoI
	StudentsRepo               repository.StudentsRepoI
	LessonsRepo                repository.LessonsRepoI
	CheckInterval              time.Duration
	dialer                     *gomail.Dialer
	lessonsNotificationsConfig LessonsNotificationsConfig
}

type LessonsNotificationsConfig struct {
	layoutISODate  string
	layoutISOTime  string
	masterhubEmail string
	masterhubTheme string
}

func (lnUC *LessonNotificationsUC) Start() {
	ticker := time.NewTicker(lnUC.CheckInterval)
	for range ticker.C {
		lnUC.getSoonLessons()
	}
}
func (lnUC *LessonNotificationsUC) CheckFinished() {
	ticker := time.NewTicker(lnUC.CheckInterval)
	for range ticker.C {
		lnUC.updateFinishedLessons()
	}
}

func (lnUC *LessonNotificationsUC) getSoonLessons() {
	lessonsDB, err := lnUC.LessonsRepo.GetSoonLessons(time.Now().Format(lnUC.lessonsNotificationsConfig.layoutISODate),
		time.Now().Format(lnUC.lessonsNotificationsConfig.layoutISOTime))
	if err != nil {
		logger.Error(fmt.Errorf(lnUC.useCases.errorMessages.DbError))
		return
	}
	fmt.Println(lessonsDB)

	for _, lessonDB := range lessonsDB {
		var masterDB models.MasterDB
		masterDB.Id = lessonDB.MasterId
		emailMaster, err := lnUC.MastersRepo.GetMasterByIdWithEmail(&masterDB)
		if err != nil {
			logger.Error(fmt.Errorf(lnUC.useCases.errorMessages.DbError))
			return
		}
		var lesson models.Lesson
		err = lnUC.useCases.LessonsUC.matchLesson(&lessonDB, &lesson, masterDB.UserId)
		if err != nil {
			return
		}

		lessonStudents, err := lnUC.LessonsRepo.GetLessonStudents(lessonDB.Id)
		if err != nil {
			logger.Error(fmt.Errorf(lnUC.useCases.errorMessages.DbError))
			return
		}
		students := make([]models.StudentDB, 0)
		for _, lessonStudent := range lessonStudents {
			var studentDB models.StudentDB
			studentDB.UserId = lessonStudent
			emailStudent, err := lnUC.StudentsRepo.GetStudentByUserIdWithEmail(&studentDB)
			if err != nil {
				logger.Error(fmt.Errorf(lnUC.useCases.errorMessages.DbError))
				return
			}
			studentMessage := lnUC.createStudentMessage(&lesson, &masterDB)
			lnUC.sendToWebsocket(lessonStudent, studentMessage)
			lnUC.sendEmail(emailStudent, studentMessage)
			students = append(students, studentDB)
		}

		masterMessage := lnUC.createMasterMessage(&lesson, students)
		lnUC.sendToWebsocket(masterDB.UserId, masterMessage)
		lnUC.sendEmail(emailMaster, masterMessage)
	}
}

func (lnUC *LessonNotificationsUC) createStudentMessage(lesson *models.Lesson, master *models.MasterDB) string {
	priceFloat, _ := lesson.Price.Value.Float64()
	priceFloatString := fmt.Sprintf("%.2f", priceFloat)
	return fmt.Sprintf("У вас назначен урок с %s (%s) в %s продолжительностью %s. Он пройдет в формате %s и будет стоить: %s %s. Не забудьте посетить!",
		master.Fullname, master.Username, lesson.TimeStart, lesson.Duration, lesson.EducationFormat, priceFloatString, lesson.Price.Currency)
}

func (lnUC *LessonNotificationsUC) createMasterMessage(lesson *models.Lesson, students []models.StudentDB) string {
	priceFloat, _ := lesson.Price.Value.Float64()
	priceFloatString := fmt.Sprintf("%.2f", priceFloat)
	var studentsString string
	for _, student := range students {
		studentsString += fmt.Sprintf("%s (%s), ", student.Fullname, student.Username)
	}
	if len(studentsString) > 0 {
		studentsString = studentsString[:len(studentsString)-1]
	}
	return fmt.Sprintf("У вас назначен урок с %s в %s продолжительностью %s. Он пройдет в формате %s и будет стоить: %s %s. Не забудьте провести его!",
		studentsString, lesson.TimeStart, lesson.Duration, lesson.EducationFormat, priceFloatString, lesson.Price.Currency)
}

func (lnUC *LessonNotificationsUC) sendEmail(email string, messageText string) {
	fmt.Println(messageText, email)
	message := gomail.NewMessage()
	message.SetHeader("From", lnUC.lessonsNotificationsConfig.masterhubEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", lnUC.lessonsNotificationsConfig.masterhubTheme)
	message.SetBody("text/plain", messageText)

	if err := lnUC.dialer.DialAndSend(message); err != nil {
		logger.Error(fmt.Errorf(lnUC.useCases.errorMessages.MailSendError))
		return
	}
}

func (lnUC *LessonNotificationsUC) sendToWebsocket(userId int64, messageText string) {
	lnUC.websocketUC.SendNotification(models.WebsocketNotification{Type: 4, Notification: models.Notification{UserId: userId, Text: messageText}})
}

func (lnUC *LessonNotificationsUC) updateFinishedLessons() {
	err := lnUC.LessonsRepo.UpdateFinishedLessons(time.Now().Format(lnUC.lessonsNotificationsConfig.layoutISODate),
		time.Now().Format(lnUC.lessonsNotificationsConfig.layoutISOTime))
	if err != nil {
		logger.Error(fmt.Errorf(lnUC.useCases.errorMessages.DbError))
		return
	}
}
