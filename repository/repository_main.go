package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/logger"
	"github.com/jackc/pgx"
	"github.com/pion/webrtc/v2"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

type Repository struct {
	DbConnections  int
	mongoDB        *mongo.Database
	pool           *pgx.ConnPool
	UsersRepo      *UsersRepo
	StudentsRepo   *StudentsRepo
	MastersRepo    *MastersRepo
	ThemesRepo     *ThemesRepo
	LanguagesRepo  *LanguagesRepo
	VideosRepo     *VideosRepo
	AvatarsRepo    *AvatarsRepo
	ChatsRepo      *ChatsRepo
	WebsocketsRepo *WebsocketsRepo
	CookiesRepo    *CookiesRepo
	LessonsRepo    *LessonsRepo
	VideocallsRepo *VideocallsRepo
}

func (repository *Repository) Init(config pgx.ConnConfig) error {
	var err error
	repository.pool, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: repository.DbConnections,
	})
	if err != nil {
		return err
	}

	repository.StudentsRepo = &StudentsRepo{repository}
	repository.MastersRepo = &MastersRepo{repository}
	repository.UsersRepo = &UsersRepo{repository}
	repository.LanguagesRepo = &LanguagesRepo{repository}
	repository.ThemesRepo = &ThemesRepo{repository}
	repository.VideosRepo = &VideosRepo{
		repository: repository,
		mutex:      sync.Mutex{},
	}
	repository.ChatsRepo = &ChatsRepo{
		repository: repository,
		userMap: map[string]int64{
			"master":  1,
			"student": 2,
		}}
	repository.WebsocketsRepo = &WebsocketsRepo{
		repository:     repository,
		userConnMap:    make(map[int64]string),
		clientsMap:     make(map[string]*models.WebsocketConnection),
		NewClients:     make(chan *models.WebsocketConnection),
		DroppedClients: make(chan *models.WebsocketConnection),
		Messages:       make(chan models.WebsocketMessage)}
	repository.CookiesRepo = &CookiesRepo{
		repository:     repository,
		userKey:        "user",
		cookieKey:      "cookie",
		collectionName: "cookies",
	}
	repository.AvatarsRepo = &AvatarsRepo{
		repository:     repository,
		userKey:        "user",
		avatarKey:      "filename",
		collectionName: "avatars",
		extKey:         "extension",
	}
	repository.LessonsRepo = &LessonsRepo{repository}
	repository.VideocallsRepo = &VideocallsRepo{
		repository:  repository,
		peerConnMap: make(map[int64]chan *webrtc.Track),
	}
	err = repository.dropTables()
	if err != nil {
		return err
	}
	//err = repository.createTables()
	//if err != nil {
	//	return err
	//}
	//err = repository.fillTables()
	//if err != nil {
	//	return err
	//}
	err = repository.InitMongoDB(config.Host)
	if err != nil {
		return err
	}
	return nil
}

func (repository *Repository) InitMongoDB(host string) error {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", host))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		dbError := fmt.Errorf("can't connect to mongodb: %v", err.Error())
		logger.Errorf(dbError.Error())
		return err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		dbError := fmt.Errorf("can't ping mongodb: %v", err.Error())
		logger.Errorf(dbError.Error())
		return err
	}
	repository.mongoDB = client.Database(utils.DBName)
	//err = repository.dropCollections()
	//if err != nil {
	//	return err
	//}
	return nil
}

// mongo collection drop

func (repository *Repository) dropCollections() error {
	err := repository.mongoDB.Collection(repository.CookiesRepo.collectionName).Drop(context.TODO())
	if err != nil {
		dbError := fmt.Errorf("can't drop cookies collection: %v", err.Error())
		logger.Errorf(dbError.Error())
		return err
	}
	err = repository.mongoDB.Collection(repository.AvatarsRepo.collectionName).Drop(context.TODO())
	if err != nil {
		dbError := fmt.Errorf("can't drop avatars collection: %v", err.Error())
		logger.Errorf(dbError.Error())
		return err
	}
	return nil
}

// relation style tables creation

func (repository *Repository) createTables() error {
	_, err := repository.pool.Exec(TABLES_CREATION)
	if err != nil {
		return err
	}
	return nil
}

func (repository *Repository) dropTables() error {
	_, err := repository.pool.Exec(drop_videos)
	if err != nil {
		return err
	}
	return nil
}

func (repository *Repository) fillTables() error {
	_, err := repository.pool.Exec(TABLES_FILLING_RU)
	if err != nil {
		return err
	}
	//repository.VideosRepo.videosCount = 9 //TODO: HARDCODED VIDEOS
	return nil
}

func (repository *Repository) GetPool() *pgx.ConnPool {
	return repository.pool
}

func (repository *Repository) startTransaction() (*pgx.Tx, error) {
	db := repository.GetPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return transaction, dbError
	}
	return transaction, err
}

func (repository *Repository) commitTransaction(transaction *pgx.Tx) error {
	err := transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return dbError
	}
	return nil
}

func (repository *Repository) rollbackTransaction(transaction *pgx.Tx) error {
	err := transaction.Rollback()
	if err != nil {
		dbError := fmt.Errorf("failed to rollback: %v", err.Error())
		logger.Errorf(dbError.Error())
		return err
	}
	return err
}

func checkNullValueInt64(value sql.NullInt64) int64 {
	if value.Valid {
		return value.Int64
	} else {
		return 0
	}
}

func checkNullValueString(value sql.NullString) string {
	if value.Valid {
		return value.String
	} else {
		return ""
	}
}
