package main

import (
	"crypto/tls"
	"fmt"
	"github.com/google/logger"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	masterHub_handlers "github.com/technoZoomers/MasterHubBackend/handlers"
	"github.com/technoZoomers/MasterHubBackend/middleware"
	"github.com/technoZoomers/MasterHubBackend/repository"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/technoZoomers/MasterHubBackend/utils"
	gomail "gopkg.in/mail.v2"
	"net/http"
	"os"
	"time"
)

func main() {

	// logger initialization
	utils.LoggerSetup()
	defer utils.LoggerClose()

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "masterhub@mail.ru")

	// Set E-Mail receiver
	m.SetHeader("To", "spiridonovaalexis@mail.ru")

	// Set E-Mail subject
	m.SetHeader("Subject", "Gomail test subject")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", "This is Gomail test body")

	//fmt.Println(os.Getenv("MASTERHUB_MAIL_PASSWORD"))

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.mail.ru", 587, "masterhub@mail.ru", os.Getenv("MASTERHUB_MAIL_PASSWORD"))

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	//if err := d.DialAndSend(m); err != nil {
	//	fmt.Println(err)
	//	panic(err)
	//}

	//database initialization

	repo := repository.Repository{
		DbConnections: 20,
	}

	err := repo.Init(pgx.ConnConfig{
		Database: utils.DBName,
		Host:     "213.219.214.220",
		//Host:     "localhost",
		User:     "alexis",
		Password: "alexis",
	})

	fmt.Println("I'm stupid")

	//config, err := pgx.ParseConnectionString(os.Getenv("DATABASE_URL"))
	//if err != nil {
	//	logger.Fatalf("Couldn't initialize database: %v", err)
	//}
	//err = repo.Init(config)
	if err != nil {
		logger.Fatalf("Couldn't initialize database: %v", err)
	}

	//usecases initialization

	mhuseCases := useCases.UseCases{}

	err = mhuseCases.Init(repo.UsersRepo, repo.MastersRepo, repo.StudentsRepo, repo.ThemesRepo, repo.LanguagesRepo,
		repo.VideosRepo, repo.AvatarsRepo, repo.ChatsRepo, *repo.WebsocketsRepo, repo.CookiesRepo, repo.LessonsRepo, *repo.VideocallsRepo)
	if err != nil {
		logger.Fatalf("Couldn't initialize useCases: %v", err)
	}

	//handlers initialization

	mhHandlers := masterHub_handlers.Handlers{}

	err = mhHandlers.Init(mhuseCases.UsersUC, mhuseCases.MastersUC, mhuseCases.StudentsUC, mhuseCases.ThemesUC, mhuseCases.LanguagesUC,
		mhuseCases.VideosUC, mhuseCases.AvatarsUC, mhuseCases.ChatsUC, mhuseCases.WebsocketsUC, mhuseCases.LessonsUC, mhuseCases.VideocallsUC)
	if err != nil {
		logger.Fatalf("Couldn't initialize handlers: %v", err)
	}

	// middlewares initialization

	mhMiddlewares := middleware.Middlewares{}

	err = mhMiddlewares.Init(mhuseCases.UsersUC)
	if err != nil {
		logger.Fatalf("Couldn't initialize middlewares: %v", err)
	}

	// router initialization

	routerMain := mux.NewRouter()
	r := routerMain.PathPrefix(utils.Prefix).Subrouter()

	// users
	r.Handle("/users", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.UsersHandlers.CheckAuth, false)).Methods("GET")
	r.Handle("/users/login", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.UsersHandlers.Login, true)).Methods("POST")
	r.Handle("/users/logout", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.UsersHandlers.Logout, false)).Methods("DELETE")
	r.Handle("/users/{id}", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.UsersHandlers.GetUserById, false)).Methods("GET")
	r.Handle("/users/{id}/chats", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.ChatsHandlers.GetChatsByUserId, false)).Methods("GET")
	//r.Handle("/users/{id}/avatars", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.AvatarsHandlers.UploadAvatar, false)).Methods("POST")
	r.Handle("/users/{id}/avatars", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.AvatarsHandlers.ChangeAvatar, false)).Methods("PUT")
	r.HandleFunc("/users/{id}/avatars", mhHandlers.AvatarsHandlers.GetAvatar).Methods("GET")

	// interactions

	r.Handle("/interactions/users/{id}", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.WSHandlers.UpgradeConnection, false))

	//languages

	r.HandleFunc("/languages", mhHandlers.LanguagesHandlers.Get).Methods("GET")

	//themes

	r.HandleFunc("/themes", mhHandlers.ThemesHandlers.Get).Methods("GET")
	r.HandleFunc("/themes/{id}", mhHandlers.ThemesHandlers.GetThemeById).Methods("GET")

	// students
	r.HandleFunc("/students/create", mhHandlers.StudentsHandlers.Register).Methods("POST")
	r.HandleFunc("/students/{id}", mhHandlers.StudentsHandlers.GetStudentById).Methods("GET")
	r.Handle("/students/{id}", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.StudentsHandlers.ChangeStudentData, false)).Methods("PUT")
	r.Handle("/students/{id}/chats", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.ChatsHandlers.CreateChatRequest, false)).Methods("POST")
	r.Handle("/students/{id}/lessons/{lessonId}", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.LessonsHandlers.CreateLessonRequest, false)).Methods("PUT")
	r.Handle("/students/{id}/lessons/{lessonId}", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.LessonsHandlers.DeleteLessonRequest, false)).Methods("DELETE")
	r.Handle("/students/{id}/lessons", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.LessonsHandlers.GetStudentsLessons, false)).Methods("GET")
	r.Handle("/students/{id}/lessons/requests", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.LessonsHandlers.GetStudentsLessonRequests, false)).Methods("GET")

	//masters
	r.HandleFunc("/masters", mhHandlers.MastersHandlers.Get).Methods("GET")
	r.HandleFunc("/masters/create", mhHandlers.MastersHandlers.Register).Methods("POST")
	r.HandleFunc("/masters/{id}", mhHandlers.MastersHandlers.GetMasterById).Methods("GET")
	r.Handle("/masters/{id}", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.MastersHandlers.ChangeMasterData, false)).Methods("PUT")
	r.Handle("/masters/{id}/videos/create", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.VideosHandlers.Upload, false)).Methods("POST")
	r.HandleFunc("/masters/{id}/videos/{videoId}", mhHandlers.VideosHandlers.GetVideoById).Methods("GET")
	r.HandleFunc("/masters/{id}/videos/{videoId}/preview", mhHandlers.VideosHandlers.GetVideoPreviewById).Methods("GET")

	r.Handle("/masters/{id}/videos/{videoId}", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.VideosHandlers.DeleteVideoById, false)).Methods("DELETE")
	r.HandleFunc("/masters/{id}/videos/{videoId}/data", mhHandlers.VideosHandlers.GetVideoDataById).Methods("GET")
	r.Handle("/masters/{id}/videos/{videoId}/data", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.VideosHandlers.ChangeVideoData, false)).Methods("PUT")
	r.HandleFunc("/masters/{id}/videos", mhHandlers.VideosHandlers.GetVideosByMasterId).Methods("GET")
	r.Handle("/masters/{id}/intro", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.VideosHandlers.UploadIntro, false)).Methods("POST")
	r.Handle("/masters/{id}/intro", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.VideosHandlers.ChangeIntro, false)).Methods("PUT")
	r.Handle("/masters/{id}/intro", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.VideosHandlers.DeleteIntro, false)).Methods("DELETE")
	r.HandleFunc("/masters/{id}/intro", mhHandlers.VideosHandlers.GetIntro).Methods("GET")
	r.HandleFunc("/masters/{id}/intro/preview", mhHandlers.VideosHandlers.GetIntroPreview).Methods("GET")

	r.Handle("/masters/{id}/intro/data", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.VideosHandlers.ChangeIntroData, false)).Methods("PUT")
	r.HandleFunc("/masters/{id}/intro/data", mhHandlers.VideosHandlers.GetIntroData).Methods("GET")
	r.Handle("/masters/{id}/chats", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.ChatsHandlers.CreateChatByMaster, false)).Methods("POST")
	r.Handle("/masters/{id}/chats/{chatId}", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.ChatsHandlers.ChangeChatStatus, false)).Methods("PUT")
	r.Handle("/masters/{id}/lessons", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.LessonsHandlers.CreateLesson, false)).Methods("POST")
	r.HandleFunc("/masters/{id}/lessons", mhHandlers.LessonsHandlers.Get).Methods("GET")
	r.Handle("/masters/{id}/lessons/{lessonId}", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.LessonsHandlers.ChangeLessonInfo, false)).Methods("PUT")
	r.Handle("/masters/{id}/lessons/{lessonId}", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.LessonsHandlers.DeleteLesson, false)).Methods("DELETE")
	r.Handle("/masters/{id}/lessons/{lessonId}/students", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.LessonsHandlers.GetLessonStudents, false)).Methods("GET")
	r.Handle("/masters/{id}/lessons/requests", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.LessonsHandlers.GetLessonRequests, false)).Methods("GET")
	r.Handle("/masters/{id}/lessons/{lessonId}/requests", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.LessonsHandlers.ChangeLessonRequest, false)).Methods("PUT")

	//videos
	r.HandleFunc("/videos", mhHandlers.VideosHandlers.Get).Methods("GET")

	//videocalls
	r.Handle("/videocalls/users/{id}/peers/{peerId}/create", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.VCHandlers.Create, true)).Methods("POST")
	r.Handle("/videocalls/users/{id}/peers/{peerId}/connect", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.VCHandlers.Connect, true)).Methods("POST")

	//chats
	r.Handle("/chats/{id}/messages", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.ChatsHandlers.GetMessagesByChatId, false)).Methods("GET")

	//lessons
	r.HandleFunc("/lessons/{id}", mhHandlers.LessonsHandlers.GetLessonById).Methods("GET")

	cors := handlers.CORS(handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{"X-Content-Type-Options", "Access-Control-Allow-Origin", "X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "DELETE"}),
		handlers.AllowedOrigins([]string{
			"http://213.219.214.220:8080",
			"http://192.168.1.102:8080",
			"http://172.17.21.178:8080",
			"http://192.168.1.4:8083",
			"http://localhost:8080",
			"https://localhost:8080",
			"http://127.0.0.1:4200",
			"http://localhost:8083",
			"http://127.0.0.1:8083",
			"http://127.0.0.1:8080",
			"http://masterhub.site",
			"http://www.masterhub.site"}))

	// server initialization

	server := &http.Server{
		//Addr:         ":" + os.Getenv("PORT"),
		Addr:         utils.PortNum,
		Handler:      cors(r),
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}
	err = server.ListenAndServe()

	if err != nil {
		logger.Fatalf("failed to start server: %v", err)
	}
}
