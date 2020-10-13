package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/utils"
)

type VideosRepo struct {
}

func (videosRepo *VideosRepo) InsertVideoData(video *models.VideoDB) error {
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return dbError
	}
	row := transaction.QueryRow("INSERT INTO videos (master_id, filename, extension, intro, uploaded) VALUES ($1, $2, $3, $4, $5) returning id",
		video.MasterId, video.Filename, video.Extension, video.Intro, video.Uploaded)
	err = row.Scan(&video.Id)
	if err != nil {
		logger.Errorf("failed to scan row: %v", err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			logger.Errorf("failed to rollback: %v", err)
			return errRollback
		}
		return err
	}
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("error commit: %v", err.Error())
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
		dbError := fmt.Errorf("failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return countVideo, dbError
	}
	row := transaction.QueryRow("SELECT COUNT(*) FROM videos")
	err = row.Scan(&countVideo)
	if err != nil {
		logger.Errorf("failed to scan row: %v", err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			logger.Errorf("failed to rollback: %v", err)
			return countVideo, errRollback
		}
		return countVideo, err
	}
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("error commit: %v", err.Error())
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
		dbError := fmt.Errorf("failed to start transaction: %v", err.Error())
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
		err = rows.Scan(&videoDB.Id, &videoDB.MasterId, &videoDB.Filename, &videoDB.Extension,  &videoDB.Name, &videoDB.Description, &videoDB.Intro, &theme, &videoDB.Uploaded)
		if err != nil {
			logger.Errorf("failed to retrieve video data: %v", err)
			errRollback := transaction.Rollback()
			if errRollback != nil {
				logger.Errorf("failed to rollback: %v", err)
				return videos, errRollback
			}
			return videos, err
		}
		videoDB.Theme = checkNullTheme(theme)
		videos = append(videos, videoDB)
	}
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("error commit: %v", err.Error())
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

func (videosRepo *VideosRepo) GetVideoDataById(videoDB *models.VideoDB) (int64, error) {
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return utils.SERVER_ERROR, dbError
	}
	row := transaction.QueryRow("SELECT * FROM videos WHERE id=$1 and master_id=$2", videoDB.Id, videoDB.MasterId)
	var theme sql.NullInt64
	err = row.Scan(&videoDB.Id, &videoDB.MasterId, &videoDB.Filename, &videoDB.Extension, &videoDB.Name, &videoDB.Description, &videoDB.Intro, &theme, &videoDB.Uploaded)
	if err != nil {
		logger.Errorf("failed to retrieve video data: %v", err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			logger.Errorf("failed to rollback: %v", err)
			return utils.SERVER_ERROR, errRollback
		}
		return utils.USER_ERROR, fmt.Errorf("this video doesn't exist or doesn't belong to this master")
	}
	videoDB.Theme = checkNullTheme(theme)
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return utils.SERVER_ERROR, dbError
	}
	return utils.NO_ERROR, nil
}

func (videosRepo *VideosRepo) GetVideoSubthemesById(videoId int64) ([]int64, error) {
	subthemesIds := make([]int64, 0)
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return subthemesIds, dbError
	}
	rows, err := transaction.Query(`SELECT subtheme_id FROM videos_subthemes WHERE video_id=$1`, videoId)
	if err != nil {
		return subthemesIds, nil
	}
	for rows.Next() {
		var subthemeIdFound int64
		err = rows.Scan(&subthemeIdFound)
		if err != nil {
			logger.Errorf("failed to retrieve subtheme: %v", err)
			errRollback := transaction.Rollback()
			if errRollback != nil {
				logger.Errorf("failed to rollback: %v", err)
				return subthemesIds, errRollback
			}
			return subthemesIds, err
		}
		subthemesIds = append(subthemesIds, subthemeIdFound)
	}
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return subthemesIds, err
	}
	return subthemesIds, nil
}

func (videosRepo *VideosRepo) DeleteVideoSubthemesById(videoId int64) error {
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return dbError
	}
	_, err = transaction.Exec("DELETE FROM videos_subthemes WHERE video_id=$1", videoId)
	if err != nil {
		logger.Errorf("failed to delete subthemes: %v", err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			logger.Errorf("failed to rollback: %v", err)
			return errRollback
		}
		return err
	}
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return dbError
	}
	return nil
}