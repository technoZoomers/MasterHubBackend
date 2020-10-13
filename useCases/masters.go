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
	MastersRepo   repository.MastersRepoI
	ThemesRepo    repository.ThemesRepoI
	LanguagesRepo repository.LanguagesRepoI
}

var qualificationMap = map[int64]string{1: "self-educated", 2: "professional"}
var educationFormatMap = map[int64]string{1: "online", 2: "live"}

func (mastersUC *MastersUC) GetMasterById(master *models.Master) (bool, error) {
	if master.UserId == utils.ERROR_ID {
		return true, fmt.Errorf("incorrect master id")
	}
	var masterDB models.MasterDB
	masterDB.UserId = master.UserId
	errType, err := mastersUC.MastersRepo.GetMasterByUserId(&masterDB)
	if err != nil {
		if errType == utils.USER_ERROR {
			return true, err
		} else if errType == utils.SERVER_ERROR {
			return false, err
		}
	}
	master.Username = masterDB.Username
	master.Fullname = masterDB.Fullname
	master.Description = masterDB.Description
	err = setEducationFormat(master, masterDB.EducationFormat)
	if err != nil {
		return false, err
	}
	err = setQualification(master, masterDB.Qualification)
	if err != nil {
		return false, err
	}
	err = setAveragePrice(master, masterDB.AveragePrice)
	if err != nil {
		return false, err
	}
	err = mastersUC.setTheme(master, masterDB.Theme)
	if err != nil {
		return false, err
	}
	err = mastersUC.setSubThemes(master, &masterDB)
	if err != nil {
		return false, err
	}
	err = mastersUC.setLanguages(master, &masterDB)
	if err != nil {
		return false, err
	}
	return false, nil
}

func setEducationFormat(master *models.Master, format int64) error {
	if !(format <= 3 && format >= 1) {
		formatError := fmt.Errorf("wrong education format type")
		logger.Errorf(formatError.Error())
		return formatError
	}
	var formats []string
	if format == 3 {
		formats = append(formats, educationFormatMap[1], educationFormatMap[2])
	} else {
		formats = append(formats, educationFormatMap[format])
	}
	master.EducationFormat = formats
	return nil
}

func setQualification(master *models.Master, qualification int64) error {
	if !(qualification == 1 || qualification == 2) {
		formatError := fmt.Errorf("wrong qualification type")
		logger.Errorf(formatError.Error())
		return formatError
	}
	master.Qualification = qualificationMap[qualification]
	return nil
}

func (mastersUC *MastersUC) setLanguages(master *models.Master, masterDB *models.MasterDB) error {
	var langs []string
	langsIds, err := mastersUC.MastersRepo.GetMasterLanguagesById(masterDB.Id)
	if err != nil {
		master.Languages = langs
		return err
	}
	for _, langId := range langsIds {
		var language models.LanguageDB
		language.Id = langId
		_, err = mastersUC.LanguagesRepo.GetLanguageById(&language)
		if err != nil {
			master.Languages = langs
			return err
		}
		langs = append(langs, language.Name)
	}
	master.Languages = langs
	return nil
}

func (mastersUC *MastersUC) setTheme(master *models.Master, theme int64) error {
	var themeDB models.ThemeDB
	themeDB.Id = theme
	_, err := mastersUC.ThemesRepo.GetThemeById(&themeDB)
	if err != nil {
		return err
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
		return err
	}
	for _, subthemeId := range subthemesIds {
		var subtheme models.SubthemeDB
		subtheme.Id = subthemeId
		_, err = mastersUC.ThemesRepo.GetSubthemeById(&subtheme)
		if err != nil {
			master.Theme.Subthemes = subthemes
			return err
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
