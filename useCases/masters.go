package useCases

import (
	"fmt"
	"github.com/google/logger"
	"github.com/shopspring/decimal"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
	"github.com/technoZoomers/MasterHubBackend/utils"
)

type MastersUC struct {
	useCases      *UseCases
	MastersRepo   repository.MastersRepoI
	ThemesRepo    repository.ThemesRepoI
	LanguagesRepo repository.LanguagesRepoI
	mastersConfig MastersConfig
}

type MastersConfig struct {
	qualificationMap        map[int64]string
	educationFormatMap   map[int64]string
}

func (mastersUC *MastersUC) GetMasterById(master *models.Master) error {
	if master.UserId == utils.ERROR_ID {
		return &models.NotFoundError{Message: "incorrect master id", RequestId: master.UserId}
	}
	var masterDB models.MasterDB
	masterDB.UserId = master.UserId
	err := mastersUC.MastersRepo.GetMasterByUserId(&masterDB)
	if err != nil {
		return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
	}
	if masterDB.Id == utils.ERROR_ID {
		absenceError := &models.NotFoundError{Message: "master doesn't exist", RequestId: master.UserId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	master.Username = masterDB.Username
	master.Fullname = masterDB.Fullname
	master.Description = masterDB.Description
	err = mastersUC.setEducationFormat(master, masterDB.EducationFormat)
	if err != nil {
		return err
	}
	err = mastersUC.setQualification(master, masterDB.Qualification)
	if err != nil {
		return err
	}
	err = setAveragePrice(master, masterDB.AveragePrice)
	if err != nil {
		return err
	}
	err = mastersUC.setTheme(master, masterDB.Theme)
	if err != nil {
		return err
	}
	err = mastersUC.setSubThemes(master, &masterDB)
	if err != nil {
		return err
	}
	err = mastersUC.setLanguages(master, &masterDB)
	if err != nil {
		return err
	}
	return nil
}

func (mastersUC *MastersUC) setEducationFormat(master *models.Master, format int64) error {
	if !(format <= 3 && format >= 1) {
		formatError := fmt.Errorf("wrong education format type")
		logger.Errorf(formatError.Error())
		return formatError
	}
	var formats []string
	if format == 3 {
		formats = append(formats, mastersUC.mastersConfig.educationFormatMap[1], mastersUC.mastersConfig.educationFormatMap[2])
	} else {
		formats = append(formats, mastersUC.mastersConfig.educationFormatMap[format])
	}
	master.EducationFormat = formats
	return nil
}

func (mastersUC *MastersUC) setQualification(master *models.Master, qualification int64) error {
	if !(qualification == 1 || qualification == 2) {
		formatError := fmt.Errorf("wrong qualification type")
		logger.Errorf(formatError.Error())
		return formatError
	}
	master.Qualification = mastersUC.mastersConfig.qualificationMap[qualification]
	return nil
}

func (mastersUC *MastersUC) setLanguages(master *models.Master, masterDB *models.MasterDB) error {
	var langs []string
	langsIds, err := mastersUC.MastersRepo.GetMasterLanguagesById(masterDB.Id)
	if err != nil {
		master.Languages = langs
		return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
	}
	for _, langId := range langsIds {
		var language models.LanguageDB
		language.Id = langId
		err = mastersUC.LanguagesRepo.GetLanguageById(&language)
		if err != nil {
			master.Languages = langs
			return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
		}
		if language.Name != "" {
			langs = append(langs, language.Name)
		}
	}
	master.Languages = langs
	return nil
}

func (mastersUC *MastersUC) setTheme(master *models.Master, theme int64) error {
	var themeDB models.ThemeDB
	themeDB.Id = theme
	err := mastersUC.ThemesRepo.GetThemeById(&themeDB)
	if err != nil {
		return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
	}
	if themeDB.Name == "" {
		absenceError := fmt.Errorf("theme doesn't exist")
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	master.Theme.Id = theme
	master.Theme.Theme = themeDB.Name
	return nil
}

func (mastersUC *MastersUC) setSubThemes(master *models.Master, masterDB *models.MasterDB) error {
	var subthemes []string
	subthemesIds, err := mastersUC.MastersRepo.GetMasterSubthemesById(masterDB.Id)
	if err != nil {
		master.Theme.Subthemes = subthemes
		return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
	}
	for _, subthemeId := range subthemesIds {
		var subtheme models.SubthemeDB
		subtheme.Id = subthemeId
		err = mastersUC.ThemesRepo.GetSubthemeById(&subtheme)
		if err != nil {
			master.Theme.Subthemes = subthemes
			return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
		}
		if subtheme.Name == "" {
			absenceError := fmt.Errorf("subtheme doesn't exist")
			logger.Errorf(absenceError.Error())
			master.Theme.Subthemes = subthemes
			return absenceError
		}
		subthemes = append(subthemes, subtheme.Name)
	}
	master.Theme.Subthemes = subthemes
	return nil
}

func setAveragePrice(master *models.Master, avgPrice decimal.Decimal) error {
	master.AveragePrice.Value = avgPrice
	master.AveragePrice.Currency = "rub" //TODO: change to different currencies
	return nil
}

func (mastersUC *MastersUC) ChangeMasterData(master *models.Master) error {
	masterDB := models.MasterDB{
		UserId: master.UserId,
	}
	err := mastersUC.MastersRepo.GetMasterByUserId(&masterDB)
	if err != nil {
		return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
	}
	if masterDB.Id == utils.ERROR_ID {
		absenceError := &models.NotFoundError{Message: "master doesn't exist", RequestId: master.UserId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	return nil
}
