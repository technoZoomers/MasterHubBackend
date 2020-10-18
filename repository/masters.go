package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
)

type MastersRepo struct {
	repository *Repository
}

func (mastersRepo *MastersRepo) GetMasterByUserId(master *models.MasterDB) error {
	var dbError error
	transaction, err := mastersRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	var theme sql.NullInt64
	row := transaction.QueryRow("SELECT * FROM masters WHERE user_id=$1", master.UserId)
	err = row.Scan(&master.Id, &master.UserId, &master.Username, &master.Fullname, &theme,
		&master.Description, &master.Qualification, &master.EducationFormat, &master.AveragePrice)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve master: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	master.Theme = checkNullValueInt64(theme)
	err = mastersRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (mastersRepo *MastersRepo) GetMasterIdByUsername(master *models.MasterDB) error {
	var dbError error
	transaction, err := mastersRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT id FROM masters WHERE username=$1", master.Username)
	err = row.Scan(&master.Id)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve master id: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = mastersRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (mastersRepo *MastersRepo) GetMasterSubthemesById(masterId int64) ([]int64, error) {
	subthemesIds := make([]int64, 0)
	var dbError error
	transaction, err := mastersRepo.repository.startTransaction()
	if err != nil {
		return subthemesIds, err
	}
	rows, err := transaction.Query(`SELECT subtheme_id FROM masters_subthemes WHERE master_id=$1`, masterId)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve subtheme ids: %v", err.Error())
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
	err = mastersRepo.repository.commitTransaction(transaction)
	if err != nil {
		return subthemesIds, err
	}
	return subthemesIds, nil
}

func (mastersRepo *MastersRepo) GetMasterLanguagesById(masterId int64) ([]int64, error) {
	languagesIds := make([]int64, 0)
	var dbError error
	transaction, err := mastersRepo.repository.startTransaction()
	if err != nil {
		return languagesIds, err
	}
	rows, err := transaction.Query(`SELECT language_id FROM masters_languages WHERE master_id=$1`, masterId)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve language ids: %v", err.Error())
		logger.Errorf(dbError.Error())
		return languagesIds, dbError
	}
	for rows.Next() {
		var languageIdFound int64
		err = rows.Scan(&languageIdFound)
		if err != nil {
			dbError = fmt.Errorf("failed to retrieve language id: %v", err)
			logger.Errorf(dbError.Error())
			return languagesIds, dbError
		}
		languagesIds = append(languagesIds, languageIdFound)
	}
	err = mastersRepo.repository.commitTransaction(transaction)
	if err != nil {
		return languagesIds, err
	}
	return languagesIds, nil
}

func (mastersRepo *MastersRepo) DeleteMasterSubthemesById(masterId int64) error {
	var dbError error
	transaction, err := mastersRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	_, err = transaction.Exec("DELETE FROM masters_subthemes WHERE master_id=$1", masterId)
	if err != nil {
		dbError = fmt.Errorf("failed to delete subthemes: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := mastersRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = mastersRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (mastersRepo *MastersRepo) SetMasterSubthemesById(masterId int64, subthemes []int64) error {
	if subthemes == nil || len(subthemes) == 0 {
		return nil
	}
	var dbError error
	transaction, err := mastersRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	var queryValues []interface{}
	insertQuery := "INSERT INTO masters_subthemes (master_id, subtheme_id) values "
	queryValues = append(queryValues, masterId)
	for i, subth := range subthemes {
		insertQuery += fmt.Sprintf("($1, $%d),", i+2)
		queryValues = append(queryValues, subth)
	}
	insertQuery = insertQuery[:len(insertQuery)-1]
	_, err = transaction.Exec(insertQuery, queryValues...)
	if err != nil {
		dbError = fmt.Errorf("failed to insert subthemes: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := mastersRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = mastersRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (mastersRepo *MastersRepo) DeleteMasterLanguagesById(masterId int64) error {
	var dbError error
	transaction, err := mastersRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	_, err = transaction.Exec("DELETE FROM masters_languages WHERE master_id=$1", masterId)
	if err != nil {
		dbError = fmt.Errorf("failed to delete languages: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := mastersRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = mastersRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (mastersRepo *MastersRepo) SetMasterLanguagesById(masterId int64, languages []int64) error {
	if languages == nil || len(languages) == 0 {
		return nil
	}
	var dbError error
	transaction, err := mastersRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	var queryValues []interface{}
	insertQuery := "INSERT INTO masters_languages (master_id, language_id) values "
	queryValues = append(queryValues, masterId)
	for i, subth := range languages {
		insertQuery += fmt.Sprintf("($1, $%d),", i+2)
		queryValues = append(queryValues, subth)
	}
	insertQuery = insertQuery[:len(insertQuery)-1]
	_, err = transaction.Exec(insertQuery, queryValues...)
	if err != nil {
		dbError = fmt.Errorf("failed to insert languages: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := mastersRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = mastersRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (mastersRepo *MastersRepo) UpdateMaster(master *models.MasterDB) error {
	var dbError error
	transaction, err := mastersRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	_, err = transaction.Exec("UPDATE masters SET (username, fullname, theme, description, qualification, education_format) = ($1, $2, nullif($3, 0), $4, $5, $6) where id = $7",
		master.Username, master.Fullname, master.Theme, master.Description, master.Qualification, master.EducationFormat, master.Id)
	if err != nil {
		dbError = fmt.Errorf("failed to update master: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := mastersRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = mastersRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}


func (mastersRepo *MastersRepo) GetMasters(query models.MastersQueryValuesDB) ([]models.MasterDB, error) {
	var dbError error
	masters := make([]models.MasterDB, 0)
	transaction, err := mastersRepo.repository.startTransaction()
	if err != nil {
		return masters, err
	}

	var queryValues []interface{}
	queryCount := 0
	selectQuery := "SELECT id, user_id, username, fullname, theme, description, qualification, education_format, avg_price FROM " +
		"(SELECT row_number() over (ORDER BY id) as select_id, * FROM masters "


	if len(query.Subtheme) > 0 {
		selectQuery += "  INNER JOIN (SELECT DISTINCT master_id FROM masters_subthemes WHERE subtheme_id in ("
		for _, subth := range query.Subtheme {
			queryCount++
			selectQuery += fmt.Sprintf("$%d,", queryCount)
			queryValues = append(queryValues, subth)
		}
		selectQuery = selectQuery[:len(selectQuery)-1]
		selectQuery += ")) as s on s.master_id = id"
	}

	if len(query.Language) > 0 {
		selectQuery += " INNER JOIN (SELECT DISTINCT master_id FROM masters_languages WHERE language_id in ("
		for _, lang := range query.Language {
			queryCount++
			selectQuery += fmt.Sprintf("$%d,", queryCount)
			queryValues = append(queryValues, lang)
		}
		selectQuery = selectQuery[:len(selectQuery)-1]
		selectQuery += ")) as l on l.master_id = id"
	}
	if len(query.Theme) > 0 {
		selectQuery += " WHERE theme in ("
		for _, th := range query.Theme {
			queryCount++
			selectQuery += fmt.Sprintf("$%d,", queryCount)
			queryValues = append(queryValues, th)
		}
		selectQuery = selectQuery[:len(selectQuery)-1]
		selectQuery += ")"
	}
	selectQuery += ") as i"
	if query.Limit == 0 {
		if query.Offset != 0 {
			queryCount++
			selectQuery += fmt.Sprintf(" WHERE i.select_id > $%d", queryCount)
			queryValues = append(queryValues, query.Offset)
		}
	} else {
		queryCount++
		selectQuery += fmt.Sprintf(" WHERE i.select_id BETWEEN $%d", queryCount)
		queryCount++
		selectQuery += fmt.Sprintf(" AND $%d", queryCount)
		queryValues = append(queryValues, query.Offset+1, query.Offset+query.Limit)
	}
	rows, err := transaction.Query(selectQuery, queryValues...)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve masters: %v", err.Error())
		logger.Errorf(dbError.Error())
		return masters, dbError
	}

	for rows.Next() {
		var theme sql.NullInt64
		var masterFound models.MasterDB
		err = rows.Scan(&masterFound.Id, &masterFound.UserId, &masterFound.Username, &masterFound.Fullname, &theme,
			&masterFound.Description, &masterFound.Qualification, &masterFound.EducationFormat, &masterFound.AveragePrice)
		if err != nil {
			dbError = fmt.Errorf("failed to retrieve master: %v", err)
			logger.Errorf(dbError.Error())
			return masters, dbError
		}
		masterFound.Theme = checkNullValueInt64(theme)
		masters = append(masters, masterFound)
	}
	err = mastersRepo.repository.commitTransaction(transaction)
	if err != nil {
		return masters, err
	}
	return masters, nil
}