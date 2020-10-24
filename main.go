package main

import (
	"fmt"
	"github.com/google/logger"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	masterHub_handlers "github.com/technoZoomers/MasterHubBackend/handlers"
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

	fmt.Println("WHAT THE FUCK")

	//database initialization

	repo := repository.Repository{
		DbConnections: 20,
	}

	err := repo.Init(pgx.ConnConfig{
		Database: utils.DBName,
		Host:     "213.219.214.220",
		User:     "alexis",
		Password: "alexis",
	})

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
		repo.VideosRepo, repo.AvatarsRepo, repo.ChatsRepo)
	if err != nil {
		logger.Fatalf("Couldn't initialize useCases: %v", err)
	}

	//handlers initialization

	mhHandlers := masterHub_handlers.Handlers{}

	err = mhHandlers.Init(mhuseCases.UsersUC, mhuseCases.MastersUC, mhuseCases.StudentsUC, mhuseCases.ThemesUC, mhuseCases.LanguagesUC,
		mhuseCases.VideosUC, mhuseCases.AvatarsUC, mhuseCases.ChatsUC)
	if err != nil {
		logger.Fatalf("Couldn't initialize handlers: %v", err)
	}

	// router initialization

	r := mux.NewRouter()

	// users
	r.HandleFunc("/users/{id}/chats", mhHandlers.ChatsHandlers.GetChatsByUserId).Methods("GET")

	//languages

	r.HandleFunc("/languages", mhHandlers.LanguagesHandlers.Get).Methods("GET")

	//themes

	r.HandleFunc("/themes", mhHandlers.ThemesHandlers.Get).Methods("GET")
	r.HandleFunc("/themes/{id}", mhHandlers.ThemesHandlers.GetThemeById).Methods("GET")

	// students
	r.HandleFunc("/students/{id}/chats", mhHandlers.ChatsHandlers.CreateChatRequest).Methods("POST")


	//masters
	r.HandleFunc("/masters", mhHandlers.MastersHandlers.Get).Methods("GET")
	r.HandleFunc("/masters/{id}", mhHandlers.MastersHandlers.GetMasterById).Methods("GET")
	r.HandleFunc("/masters/{id}", mhHandlers.MastersHandlers.ChangeMasterData).Methods("PUT")
	r.HandleFunc("/masters/{id}/videos/create", mhHandlers.VideosHandlers.Upload).Methods("POST")
	r.HandleFunc("/masters/{id}/videos/{videoId}", mhHandlers.VideosHandlers.GetVideoById).Methods("GET")
	r.HandleFunc("/masters/{id}/videos/{videoId}", mhHandlers.VideosHandlers.DeleteVideoById).Methods("DELETE")
	r.HandleFunc("/masters/{id}/videos/{videoId}/data", mhHandlers.VideosHandlers.GetVideoDataById).Methods("GET")
	r.HandleFunc("/masters/{id}/videos/{videoId}/data", mhHandlers.VideosHandlers.ChangeVideoData).Methods("PUT")
	r.HandleFunc("/masters/{id}/videos", mhHandlers.VideosHandlers.GetVideosByMasterId).Methods("GET")
	r.HandleFunc("/masters/{id}/intro", mhHandlers.VideosHandlers.UploadIntro).Methods("POST")
	r.HandleFunc("/masters/{id}/intro", mhHandlers.VideosHandlers.ChangeIntro).Methods("PUT")
	r.HandleFunc("/masters/{id}/intro", mhHandlers.VideosHandlers.DeleteIntro).Methods("DELETE")
	r.HandleFunc("/masters/{id}/intro", mhHandlers.VideosHandlers.GetIntro).Methods("GET")
	r.HandleFunc("/masters/{id}/intro/data", mhHandlers.VideosHandlers.ChangeIntroData).Methods("PUT")
	r.HandleFunc("/masters/{id}/intro/data", mhHandlers.VideosHandlers.GetIntroData).Methods("GET")
	//r.HandleFunc("/masters/{id}/chats/{chatId}", mhHandlers.ChatsHandlers.ChangeChatStatus).Methods("PUT")

	//videos
	r.HandleFunc("/videos", mhHandlers.VideosHandlers.Get).Methods("GET")


	cors := handlers.CORS(handlers.AllowCredentials(),handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "DELETE", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))

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
