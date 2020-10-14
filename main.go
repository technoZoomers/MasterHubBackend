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
	//err := repository.Init(pgx.ConnConfig{
	//	Database: utils.DBName,
	//	Host:     "localhost",
	//	User:     "alexis",
	//	Password: "sinope27",
	//})

	config, err := pgx.ParseConnectionString(os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatalf("Couldn't initialize database: %v", err)
	}
	err = repository.Init(config)
	if err != nil {
		logger.Fatalf("Couldn't initialize database: %v", err)
	}

	err = useCases.Init(repository.GetUsersRepo(), repository.GetMastersRepo(), repository.GetStudentsRepo(), repository.GetThemesRepo(), repository.GetLanguagesRepo(),
		repository.GetVideosRepo(), repository.GetAvatarsRepo())
	if err != nil {
		logger.Fatalf("Couldn't initialize useCases: %v", err)
	}

	err = masterHub_handlers.Init(useCases.GetUsersUC(), useCases.GetMastersUC(), useCases.GetStudentsUC(), useCases.GetThemesUC(), useCases.GetLanguagesUC(),
		useCases.GetVideosUC(), useCases.GetAvatarsUC())
	if err != nil {
		logger.Fatalf("Couldn't initialize handlers: %v", err)
	}

	// router initialization

	r := mux.NewRouter()

	//languages

	r.HandleFunc("/languages", masterHub_handlers.GetLanguagesH().Get).Methods("GET")

	//themes

	r.HandleFunc("/themes", masterHub_handlers.GetThemesH().Get).Methods("GET")
	r.HandleFunc("/themes/{id}", masterHub_handlers.GetThemesH().GetThemeById).Methods("GET")

	//masters

	r.HandleFunc("/masters/{id}", masterHub_handlers.GetMastersH().GetMasterById).Methods("GET")
	r.HandleFunc("/masters/{id}/videos/create", masterHub_handlers.GetVideosH().Upload).Methods("POST")
	r.HandleFunc("/masters/{id}/videos/{videoId}", masterHub_handlers.GetVideosH().GetVideoById).Methods("GET")
	r.HandleFunc("/masters/{id}/videos/{videoId}/data", masterHub_handlers.GetVideosH().GetVideoDataById).Methods("GET")
	r.HandleFunc("/masters/{id}/videos/{videoId}/data", masterHub_handlers.GetVideosH().ChangeVideoData).Methods("PUT")
	r.HandleFunc("/masters/{id}/videos", masterHub_handlers.GetVideosH().GetVideosByMasterId).Methods("GET")



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
