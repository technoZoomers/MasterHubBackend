package main

import (
	"github.com/google/logger"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	masterHub_handlers "github.com/technoZoomers/MasterHubBackend/handlers"
	"github.com/technoZoomers/MasterHubBackend/repository"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"net/http"
	"os"
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

	//err := repo.Init(pgx.ConnConfig{
	//	Database: utils.DBName,
	//	Host:     "localhost",
	//	User:     "alexis",
	//	Password: "sinope27",
	//})

	config, err := pgx.ParseConnectionString(os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatalf("Couldn't initialize database: %v", err)
	}
	err = repo.Init(config)
	if err != nil {
		logger.Fatalf("Couldn't initialize database: %v", err)
	}

	//usecases initialization

	mhuseCases := useCases.UseCases{}

	err = mhuseCases.Init(repo.UsersRepo, repo.MastersRepo, repo.StudentsRepo, repo.ThemesRepo, repo.LanguagesRepo,
		repo.VideosRepo, repo.AvatarsRepo)
	if err != nil {
		logger.Fatalf("Couldn't initialize useCases: %v", err)
	}

	//handlers initialization

	mhHandlers := masterHub_handlers.Handlers{}

	err = mhHandlers.Init(mhuseCases.UsersUC, mhuseCases.MastersUC, mhuseCases.StudentsUC, mhuseCases.ThemesUC, mhuseCases.LanguagesUC,
		mhuseCases.VideosUC, mhuseCases.AvatarsUC)
	if err != nil {
		logger.Fatalf("Couldn't initialize handlers: %v", err)
	}

	// router initialization

	r := mux.NewRouter()

	//languages

	r.HandleFunc("/languages", mhHandlers.LanguagesHandlers.Get).Methods("GET")

	//themes

	r.HandleFunc("/themes", mhHandlers.ThemesHandlers.Get).Methods("GET")
	r.HandleFunc("/themes/{id}", mhHandlers.ThemesHandlers.GetThemeById).Methods("GET")

	//masters

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

	cors := handlers.CORS(handlers.AllowCredentials(), handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE"}))

	// server initialization

	server := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		//Addr:         utils.PortNum,
		Handler:      cors(r),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = server.ListenAndServe()

	if err != nil {
		logger.Fatalf("failed to start server: %v", err)
	}
}
