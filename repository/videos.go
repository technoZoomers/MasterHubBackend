package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
)

type VideosRepo struct {
	repository *Repository
}

func (videosRepo *VideosRepo) InsertVideoData(video *models.VideoDB) error {
	var dbError error
	transaction, err := videosRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("INSERT INTO videos (master_id, filename, extension, intro, uploaded) VALUES ($1, $2, $3, $4, $5) returning id",
		video.MasterId, video.Filename, video.Extension, video.Intro, video.Uploaded)
	err = row.Scan(&video.Id)
	if err != nil {
		dbError = fmt.Errorf("failed to insert video data: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := videosRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = videosRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (videosRepo *VideosRepo) CountVideos() (int64, error) {
	var countVideo int64 = 0
	var dbError error
	transaction, err := videosRepo.repository.startTransaction()
	if err != nil {
		return countVideo,err
	}
	row := transaction.QueryRow("SELECT COUNT(*) FROM videos")
	err = row.Scan(&countVideo)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve number of videos: %v", err.Error())
		logger.Errorf(dbError.Error())
		return countVideo, dbError
	}
	err = videosRepo.repository.commitTransaction(transaction)
	if err != nil {
		return countVideo, err
	}
	return countVideo, nil
}

func (videosRepo *VideosRepo) GetVideosByMasterId(masterId int64) ([]models.VideoDB, error) {
	var dbError error
	videos := make([]models.VideoDB, 0)
	transaction, err := videosRepo.repository.startTransaction()
	if err != nil {
		return videos,err
	}
	rows, err := transaction.Query(`SELECT * FROM videos WHERE master_id=$1`, masterId)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve videos: %v", err.Error())
		logger.Errorf(dbError.Error())
		return videos, dbError
	}
	for rows.Next() {
		var videoDB models.VideoDB
		var theme sql.NullInt64
		err = rows.Scan(&videoDB.Id, &videoDB.MasterId, &videoDB.Filename, &videoDB.Extension, &videoDB.Name, &videoDB.Description, &videoDB.Intro, &theme, &videoDB.Uploaded)
		if err != nil {
			dbError = fmt.Errorf("failed to retrieve video data: %v", err)
			logger.Errorf(dbError.Error())
			return videos, dbError
		}
		videoDB.Theme = checkNullValueInt64(theme)
		videos = append(videos, videoDB)
	}
	err = videosRepo.repository.commitTransaction(transaction)
	if err != nil {
		return videos, err
	}
	return videos, nil
}

func (videosRepo *VideosRepo) GetVideoDataByIdAndMasterId(videoDB *models.VideoDB) error {
	var dbError error
	transaction, err := videosRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT * FROM videos WHERE id=$1 and master_id=$2", videoDB.Id, videoDB.MasterId)
	var theme sql.NullInt64
	err = row.Scan(&videoDB.Id, &videoDB.MasterId, &videoDB.Filename, &videoDB.Extension, &videoDB.Name, &videoDB.Description, &videoDB.Intro, &theme, &videoDB.Uploaded)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve video: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	videoDB.Theme = checkNullValueInt64(theme)
	err = videosRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (videosRepo *VideosRepo) GetIntroByMasterId(videoDB *models.VideoDB) error {
	var dbError error
	transaction, err := videosRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT * FROM videos WHERE intro=true and master_id=$1", videoDB.MasterId)
	var theme sql.NullInt64
	err = row.Scan(&videoDB.Id, &videoDB.MasterId, &videoDB.Filename, &videoDB.Extension, &videoDB.Name, &videoDB.Description, &videoDB.Intro, &theme, &videoDB.Uploaded)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve intro: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	videoDB.Theme = checkNullValueInt64(theme)
	err = videosRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (videosRepo *VideosRepo) GetVideoSubthemesById(videoId int64) ([]int64, error) {
	var dbError error
	subthemesIds := make([]int64, 0)
	transaction, err := videosRepo.repository.startTransaction()
	if err != nil {
		return subthemesIds, err
	}
	rows, err := transaction.Query(`SELECT subtheme_id FROM videos_subthemes WHERE video_id=$1`, videoId)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve video subthemes: %v", err.Error())
		logger.Errorf(dbError.Error())
		return subthemesIds, dbError
	}
	for rows.Next() {
		var subthemeIdFound int64
		err = rows.Scan(&subthemeIdFound)
		if err != nil {
			dbError = fmt.Errorf("failed to retrieve subtheme id: %v", err)
			logger.Errorf(dbError.Error())
			return subthemesIds, dbError
		}
		subthemesIds = append(subthemesIds, subthemeIdFound)
	}
	err = videosRepo.repository.commitTransaction(transaction)
	if err != nil {
		return subthemesIds, err
	}
	return subthemesIds, nil
}

func (videosRepo *VideosRepo) DeleteVideoSubthemesById(videoId int64) error {
	var dbError error
	transaction, err := videosRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	_, err = transaction.Exec("DELETE FROM videos_subthemes WHERE video_id=$1", videoId)
	if err != nil {
		dbError = fmt.Errorf("failed to delete subthemes: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := videosRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = videosRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (videosRepo *VideosRepo) SetVideoSubthemesById(videoId int64, subthemes []int64) error {
	if subthemes == nil || len(subthemes) == 0 {
		return nil
	}
	var dbError error
	transaction, err := videosRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	var queryValues []interface{}
	insertQuery := "INSERT INTO videos_subthemes (video_id, subtheme_id) values "
	queryValues = append(queryValues, videoId)
	for i, subth := range subthemes {
		insertQuery += fmt.Sprintf("($1, $%d),", i+2)
		queryValues = append(queryValues, subth)
	}
	insertQuery = insertQuery[:len(insertQuery)-1]
	_, err = transaction.Exec(insertQuery, queryValues...)
	if err != nil {
		dbError = fmt.Errorf("failed to insert subthemes: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := videosRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = videosRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (videosRepo *VideosRepo) UpdateVideo(video *models.VideoDB) error {
	var dbError error
	transaction, err := videosRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	_, err = transaction.Exec("UPDATE videos SET (name, description, theme) = ($1, $2, nullif($3, 0)) WHERE id=$4", video.Name, video.Description, video.Theme, video.Id)
	if err != nil {
		dbError = fmt.Errorf("failed to update theme: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := videosRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = videosRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}
