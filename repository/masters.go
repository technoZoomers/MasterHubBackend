package repository

import (
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/utils"
)

type MastersRepo struct {
}

func (mastersRepo *MastersRepo) GetMasterByUserId(master *models.MasterDB) (int64, error) {
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("Failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return utils.SERVER_ERROR, dbError
	}
	row := transaction.QueryRow("SELECT * FROM masters WHERE user_id=$1", master.UserId)
	err = row.Scan(&master.Id, &master.UserId, &master.Username, &master.Fullname, &master.Theme,
		&master.Description, &master.Qualification, &master.EducationFormat, &master.AveragePrice)
	if err != nil {
		logger.Errorf("Failed to retrieve master: %v", err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			logger.Errorf("Failed to rollback: %v", err)
			return utils.SERVER_ERROR, errRollback
		}
		return utils.USER_ERROR, fmt.Errorf("this master doesn't exist")
	}
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("Error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return utils.SERVER_ERROR, dbError
	}
	return utils.NO_ERROR, nil
}

func (mastersRepo *MastersRepo) GetMasterSubthemesById(master *models.MasterDB) ([]int64, error) {
	subthemesIds := make([]int64, 0)
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("Failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return subthemesIds, dbError
	}
	rows, err := transaction.Query(`SELECT subtheme_id FROM masters_subthemes WHERE master_id=$1`, master.Id)
	if err != nil {
		return subthemesIds, nil
	}
	for rows.Next() {
		var subthemeIdFound int64
		err = rows.Scan(&subthemeIdFound)
		if err != nil {
			logger.Errorf("Failed to retrieve subtheme: %v", err)
			errRollback := transaction.Rollback()
			if errRollback != nil {
				logger.Errorf("Failed to rollback: %v", err)
				return subthemesIds, errRollback
			}
			return subthemesIds, err
		}
		subthemesIds = append(subthemesIds, subthemeIdFound)
	}
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("Error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return subthemesIds, err
	}
	return subthemesIds, nil
}

func (mastersRepo *MastersRepo) GetMasterLanguagesById(master *models.MasterDB) ([]int64, error) {
	languagesIds := make([]int64, 0)
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("Failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return languagesIds, dbError
	}
	rows, err := transaction.Query(`SELECT language_id FROM masters_languages WHERE master_id=$1`, master.Id)
	if err != nil {
		return languagesIds, nil
	}
	for rows.Next() {
		var languageIdFound int64
		err = rows.Scan(&languageIdFound)
		if err != nil {
			logger.Errorf("Failed to retrieve language: %v", err)
			errRollback := transaction.Rollback()
			if errRollback != nil {
				logger.Errorf("Failed to rollback: %v", err)
				return languagesIds, errRollback
			}
			return languagesIds, err
		}
		languagesIds = append(languagesIds, languageIdFound)
	}
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("Error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return languagesIds, err
	}
	return languagesIds, nil
}