package main

import (
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
	"net/http"
	"time"
)

func main() {

	// logger initialization
	utils.LoggerSetup()
	defer utils.LoggerClose()

	//files, err := ioutil.ReadDir(".")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, f := range files {
	//	fmt.Println(f.Name())
	//}


	//database initialization

	repo := repository.Repository{
		DbConnections: 20,
	}

	err := repo.Init(pgx.ConnConfig{
		Database: utils.DBName,
		Host:     "213.219.214.220",
		//Host: "localhost",
		User:     "alexis",
		Password: "alexis",
	})

	fmt.Println("WHAT THE FUCK")

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
		repo.VideosRepo, repo.AvatarsRepo, repo.ChatsRepo, *repo.WebsocketsRepo, repo.CookiesRepo)
	if err != nil {
		logger.Fatalf("Couldn't initialize useCases: %v", err)
	}

	//handlers initialization

	mhHandlers := masterHub_handlers.Handlers{}

	err = mhHandlers.Init(mhuseCases.UsersUC, mhuseCases.MastersUC, mhuseCases.StudentsUC, mhuseCases.ThemesUC, mhuseCases.LanguagesUC,
		mhuseCases.VideosUC, mhuseCases.AvatarsUC, mhuseCases.ChatsUC, mhuseCases.WebsocketsUC)
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
	r:= routerMain.PathPrefix(utils.Prefix).Subrouter()

	// users
	r.Handle("/users", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.UsersHandlers.CheckAuth, false)).Methods("GET")
	r.Handle("/users/login", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.UsersHandlers.Login, true)).Methods("POST")
	r.Handle("/users/logout", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.UsersHandlers.Logout, false)).Methods("DELETE")
	r.Handle("/users/{id}", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.UsersHandlers.GetUserById, false)).Methods("GET")
	r.Handle("/users/{id}/chats", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.ChatsHandlers.GetChatsByUserId, false)).Methods("GET")

	r.Handle("/users/{id}/interactions",  mhMiddlewares.AuthMiddleware.Auth(mhHandlers.WSHandlers.UpgradeConnection, false))


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


	//masters
	r.HandleFunc("/masters", mhHandlers.MastersHandlers.Get).Methods("GET")
	r.HandleFunc("/masters/create", mhHandlers.MastersHandlers.Register).Methods("POST")
	r.HandleFunc("/masters/{id}", mhHandlers.MastersHandlers.GetMasterById).Methods("GET")
	r.Handle("/masters/{id}", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.MastersHandlers.ChangeMasterData, false)).Methods("PUT")
	r.Handle("/masters/{id}/videos/create", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.VideosHandlers.Upload, false)).Methods("POST")
	r.HandleFunc("/masters/{id}/videos/{videoId}", mhHandlers.VideosHandlers.GetVideoById).Methods("GET")
	r.Handle("/masters/{id}/videos/{videoId}", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.VideosHandlers.DeleteVideoById, false)).Methods("DELETE")
	r.HandleFunc("/masters/{id}/videos/{videoId}/data", mhHandlers.VideosHandlers.GetVideoDataById).Methods("GET")
	r.Handle("/masters/{id}/videos/{videoId}/data", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.VideosHandlers.ChangeVideoData, false)).Methods("PUT")
	r.HandleFunc("/masters/{id}/videos", mhHandlers.VideosHandlers.GetVideosByMasterId).Methods("GET")
	r.Handle("/masters/{id}/intro", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.VideosHandlers.UploadIntro, false)).Methods("POST")
	r.Handle("/masters/{id}/intro", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.VideosHandlers.ChangeIntro, false)).Methods("PUT")
	r.Handle("/masters/{id}/intro", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.VideosHandlers.DeleteIntro, false)).Methods("DELETE")
	r.HandleFunc("/masters/{id}/intro", mhHandlers.VideosHandlers.GetIntro).Methods("GET")
	r.Handle("/masters/{id}/intro/data", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.VideosHandlers.ChangeIntroData, false)).Methods("PUT")
	r.HandleFunc("/masters/{id}/intro/data", mhHandlers.VideosHandlers.GetIntroData).Methods("GET")
	r.Handle("/masters/{id}/chats/{chatId}", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.ChatsHandlers.ChangeChatStatus, false)).Methods("PUT")

	//videos
	r.HandleFunc("/videos", mhHandlers.VideosHandlers.Get).Methods("GET")

	//chats
	r.Handle("/chats/{id}/messages", mhMiddlewares.AuthMiddleware.Auth(mhHandlers.ChatsHandlers.GetMessagesByChatId, false)).Methods("GET")

	cors := handlers.CORS(handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{"X-Content-Type-Options", "Access-Control-Allow-Origin","X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "DELETE"}),
		handlers.AllowedOrigins([]string{
			"http://213.219.214.220:8080",
			"http://192.168.1.102:8080",
			"http://172.17.21.178:8080",
			"http://localhost:8080",
			"http://127.0.0.1:8080",
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
