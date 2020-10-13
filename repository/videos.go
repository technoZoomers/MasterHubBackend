package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
)

type VideosRepo struct {
}

func (videosRepo *VideosRepo) InsertVideoData(video *models.VideoDB) error {
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("Failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return dbError
	}
	row := transaction.QueryRow("INSERT INTO videos (master_id, filename, intro, uploaded) VALUES ($1, $2, $3, $4) returning id",
		video.MasterId, video.Filename, video.Intro, video.Uploaded)
	err = row.Scan(&video.Id)
	if err != nil {
		logger.Errorf("Failed to scan row: %v", err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			logger.Errorf("Failed to rollback: %v", err)
			return errRollback
		}
		return err
	}
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("Error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return dbError
	}
	return nil
}

func (videosRepo *VideosRepo) CountVideos() (int64, error) {
	var countVideo int64
	countVideo = 0
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("Failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return countVideo, dbError
	}
	row := transaction.QueryRow("SELECT COUNT(*) FROM videos")
	err = row.Scan(&countVideo)
	if err != nil {
		logger.Errorf("Failed to scan row: %v", err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			logger.Errorf("Failed to rollback: %v", err)
			return countVideo, errRollback
		}
		return countVideo, err
	}
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("Error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return countVideo, dbError
	}
	return countVideo, nil
}

func (videosRepo *VideosRepo) GetVideosByMasterId(masterId int64) ([]models.VideoDB, error) {
	videos := make([]models.VideoDB, 0)
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("Failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return videos, dbError
	}
	rows, err := transaction.Query(`SELECT * FROM videos WHERE master_id=$1`, masterId)
	if err != nil {
		return videos, nil
	}
	for rows.Next() {
		var videoDB models.VideoDB
		var theme sql.NullInt64
		err = rows.Scan(&videoDB.Id, &videoDB.MasterId, &videoDB.Filename, &videoDB.Name, &videoDB.Description, &videoDB.Intro, &theme, &videoDB.Uploaded)
		if err != nil {
			logger.Errorf("Failed to retrieve video data: %v", err)
			errRollback := transaction.Rollback()
			if errRollback != nil {
				logger.Errorf("Failed to rollback: %v", err)
				return videos, errRollback
			}
			return videos, err
		}
		videoDB.Theme = checkNullTheme(theme)
		videos = append(videos, videoDB)
	}
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("Error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return videos, err
	}
	return videos, nil
}

func checkNullTheme(theme sql.NullInt64) int64 {
	if theme.Valid {
		return theme.Int64
	} else {
		return 0
	}
}