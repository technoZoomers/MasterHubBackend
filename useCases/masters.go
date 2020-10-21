package useCases

import (
	"fmt"
	"github.com/google/logger"
	"github.com/shopspring/decimal"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
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

	qualificationMapBackwards        map[string]int64
	educationFormatMapBackwards  map[string]int64
}

func (mastersUC *MastersUC) matchMaster(masterDB *models.MasterDB, master *models.Master) error {
	master.Username = masterDB.Username
	master.Fullname = masterDB.Fullname
	master.Description = masterDB.Description
	err := mastersUC.setEducationFormat(master, masterDB.EducationFormat)
	if err != nil {
		return err
	}
	err = mastersUC.setQualification(master, masterDB.Qualification)
	if err != nil {
		return err
	}
	err = mastersUC.setAveragePrice(master, masterDB.AveragePrice)
	if err != nil {
		return err
	}
	if masterDB.Theme != mastersUC.useCases.errorId {
		err = mastersUC.setTheme(master, masterDB.Theme)
		if err != nil {
			return err
		}
		err = mastersUC.setSubThemes(master, masterDB)
		if err != nil {
			return err
		}
	}
	err = mastersUC.setLanguages(master, masterDB)
	if err != nil {
		return err
	}
	return nil
}

func (mastersUC *MastersUC) validateMaster(masterDB *models.MasterDB, master *models.Master) error {
	if master.UserId == mastersUC.useCases.errorId {
		return &models.BadRequestError{Message: "incorrect master id", RequestId: master.UserId}
	}
	masterDB.UserId = master.UserId
	err := mastersUC.MastersRepo.GetMasterByUserId(masterDB)
	if err != nil {
		return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
	}
	if masterDB.Id == mastersUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "master doesn't exist", RequestId: master.UserId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	return nil
}

func (mastersUC *MastersUC) GetMasterById(master *models.Master) error {
	var masterDB models.MasterDB
	masterDB.UserId = master.UserId
	err := mastersUC.validateMaster(&masterDB, master)
	if err != nil {
		return err
	}
	err = mastersUC.matchMaster(&masterDB, master)
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

func (mastersUC *MastersUC) getTheme(themeDB *models.ThemeDB) error {
	err := mastersUC.ThemesRepo.GetThemeById(themeDB)
	if err != nil {
		return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
	}
	return nil
}

func (mastersUC *MastersUC) setTheme(master *models.Master, theme int64) error {
	var themeDB models.ThemeDB
	themeDB.Id = theme
	err := mastersUC.getTheme(&themeDB)
	if err != nil {
		return err
	}
	//master.Theme.Id = theme
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

func (mastersUC *MastersUC) setAveragePrice(master *models.Master, avgPrice decimal.Decimal) error {
	master.AveragePrice.Value = avgPrice
	master.AveragePrice.Currency = "rub" //TODO: change to different currencies
	return nil
}

func (mastersUC *MastersUC) ChangeMasterData(master *models.Master) error {
	if master.UserId == mastersUC.useCases.errorId {
		return &models.BadRequestError{Message: "incorrect master id", RequestId: master.UserId}
	}
	masterDB := models.MasterDB{
		UserId: master.UserId,
	}
	err := mastersUC.MastersRepo.GetMasterByUserId(&masterDB)
	if err != nil {
		return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
	}
	if masterDB.Id == mastersUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "master doesn't exist", RequestId: master.UserId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}

	if master.Username != "" {
		masterDBUsernameExist := models.MasterDB{
			Username: master.Username,
		}
		err = mastersUC.MastersRepo.GetMasterIdByUsername(&masterDBUsernameExist)
		if err != nil {
			return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
		}
		if masterDBUsernameExist.Id != mastersUC.useCases.errorId && masterDBUsernameExist.Id != masterDB.Id {
			absenceError := &models.ConflictError{Message: "can't update master, username is already taken", RequestId: master.UserId}
			logger.Errorf(absenceError.Error())
			return absenceError
		}
		masterDB.Username = master.Username
	} else {
		master.Username = masterDB.Username
	}

	var emptyPrice models.Price
	if master.AveragePrice != emptyPrice && !master.AveragePrice.Value.Equal(masterDB.AveragePrice) {
		fileError := fmt.Errorf("master average price can't be changed")
		logger.Errorf(fileError.Error())
		return fileError
	}
	masterDB.Fullname = master.Fullname
	masterDB.Description = master.Description

	err = mastersUC.changeMastersTheme(master, &masterDB)
	if err != nil {
		return  err
	}

	err = mastersUC.changeMastersQualification(master, &masterDB)
	if err != nil {
		return  err
	}

	err = mastersUC.changeMastersEducationFormat(master, &masterDB)
	if err != nil {
		return  err
	}

	err = mastersUC.MastersRepo.UpdateMaster(&masterDB)
	if err != nil {
		return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
	}

	err = mastersUC.changeMastersLanguages(master, &masterDB)
	if err != nil {
		return  err
	}

	err = mastersUC.changeMastersSubthemes(master, &masterDB)
	if err != nil {
		return  err
	}

	return nil
}

func (mastersUC *MastersUC) changeMastersEducationFormat(master *models.Master, masterDB *models.MasterDB) error {
	var err error
	lenEdFormat := len(master.EducationFormat)
	switch lenEdFormat {
	case 0:
		err = mastersUC.setEducationFormat(master, masterDB.EducationFormat)
		if err != nil {
			return err
		}
		return nil
	case 1:
		newEdFormat := mastersUC.mastersConfig.qualificationMapBackwards[master.EducationFormat[0]]
		if newEdFormat != mastersUC.useCases.errorId {
			if masterDB.EducationFormat != newEdFormat {
				masterDB.EducationFormat = newEdFormat
			}
			return nil
		}
	case 2:
		newEdFormatFirst := mastersUC.mastersConfig.qualificationMapBackwards[master.EducationFormat[0]]
		if newEdFormatFirst == mastersUC.useCases.errorId {
			break
		}
		newEdFormatSecond := mastersUC.mastersConfig.qualificationMapBackwards[master.EducationFormat[1]]
		if newEdFormatSecond == mastersUC.useCases.errorId {
			break
		}
		masterDB.EducationFormat = newEdFormatFirst + newEdFormatSecond
		return nil
	default:
		break
	}
	notExistError := &models.BadRequestError{Message: "cant't update master, education format doesn't exist", RequestId: master.UserId}
	logger.Errorf(notExistError.Error())
	return notExistError
}


func (mastersUC *MastersUC) changeMastersQualification(master *models.Master, masterDB *models.MasterDB) error {
	var err error
	if master.Qualification != "" {
		newQualification := mastersUC.mastersConfig.qualificationMapBackwards[master.Qualification]
		if newQualification != mastersUC.useCases.errorId {
			if masterDB.Qualification != newQualification {
				masterDB.Qualification = newQualification
			}
		} else {
			notExistError := &models.BadRequestError{Message: "cant't update master, qualification doesn't exist", RequestId: master.UserId}
			logger.Errorf(notExistError.Error())
			return notExistError
		}
	} else {
		err = mastersUC.setQualification(master, masterDB.Qualification)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mastersUC *MastersUC) changeMastersLanguages(master *models.Master, masterDB *models.MasterDB) error {
	var newLanguagesIds []int64
	for _, language := range master.Languages {
		languageDB := models.LanguageDB{Name: language}
		err := mastersUC.LanguagesRepo.GetLanguageByName(&languageDB)
		if err != nil {
			return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
		}
		if languageDB.Id == mastersUC.useCases.errorId {
			fileError := &models.BadRequestError{Message: "cant't update video, language doesn't exist", RequestId: master.UserId}
			logger.Errorf(fileError.Error())
			return fileError
		}
		newLanguagesIds = append(newLanguagesIds, languageDB.Id)
	}

	oldLanguagesIds, err := mastersUC.MastersRepo.GetMasterLanguagesById(masterDB.Id)
	if err != nil {
		return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
	}

	err = mastersUC.MastersRepo.DeleteMasterLanguagesById(masterDB.Id)
	if err != nil {
		return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
	}

	err = mastersUC.MastersRepo.SetMasterLanguagesById(masterDB.Id, newLanguagesIds)
	if err != nil {
		_ = mastersUC.MastersRepo.SetMasterLanguagesById(masterDB.Id, oldLanguagesIds)
		return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
	}
	return nil
}

func (mastersUC *MastersUC) changeMastersSubthemes(master *models.Master, masterDB *models.MasterDB) error {
	var err error

	if master.Theme.Theme == "" {
		err = mastersUC.MastersRepo.DeleteMasterSubthemesById(masterDB.Id)
		if err != nil {
			return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
		}
		return nil
	}

	var newSubthemesIds []int64
	for _, subtheme := range master.Theme.Subthemes {
		subthemeDB := models.SubthemeDB{Name: subtheme}
		err := mastersUC.ThemesRepo.GetSubthemeByName(&subthemeDB)
		if err != nil {
			return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
		}
		if subthemeDB.Id == mastersUC.useCases.errorId {
			fileError := &models.BadRequestError{Message: "cant't update video, subtheme doesn't exist", RequestId: master.UserId}
			logger.Errorf(fileError.Error())
			return fileError
		}
		newSubthemesIds = append(newSubthemesIds, subthemeDB.Id)
	}

	oldSubthemesIds, err := mastersUC.MastersRepo.GetMasterSubthemesById(masterDB.Id)
	if err != nil {
		return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
	}

	err = mastersUC.MastersRepo.DeleteMasterSubthemesById(masterDB.Id)
	if err != nil {
		return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
	}

	err = mastersUC.MastersRepo.SetMasterSubthemesById(masterDB.Id, newSubthemesIds)
	if err != nil {
		_ = mastersUC.MastersRepo.SetMasterSubthemesById(masterDB.Id, oldSubthemesIds)
		return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
	}
	return nil
}

func (mastersUC *MastersUC) changeMastersTheme(master *models.Master, masterDB *models.MasterDB) error {
	var err error

	if master.Theme.Theme == "" {
		masterDB.Theme = mastersUC.useCases.errorId
		return nil
	}

	var oldTheme models.ThemeDB
	oldTheme.Id = masterDB.Theme
	err = mastersUC.getTheme(&oldTheme)
	if err != nil {
		return err
	}

	if master.Theme.Theme != oldTheme.Name {
		newThemeDB := models.ThemeDB{
			Name: master.Theme.Theme,
		}
		err := mastersUC.ThemesRepo.GetThemeByName(&newThemeDB)
		if err != nil {
			return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
		}
		if newThemeDB.Id == mastersUC.useCases.errorId {
			fileError := &models.BadRequestError{Message: "cant't update master, theme doesn't exist", RequestId: master.UserId}
			logger.Errorf(fileError.Error())
			return fileError
		}
		masterDB.Theme = newThemeDB.Id
	}
	return nil
}

func (mastersUC *MastersUC) matchEducationFormat(educationFormat string) (int64, error) {
	var edFormatInt int64 = 0
	if educationFormat != "" {
		edFormatInt = mastersUC.mastersConfig.educationFormatMapBackwards[educationFormat]
		if edFormatInt != mastersUC.useCases.errorId {
			return edFormatInt, nil
		}
		badParamError := &models.BadQueryParameterError{Parameter:"educationFormat"}
		logger.Errorf(badParamError.Error())
		return edFormatInt, badParamError
	}
	return edFormatInt, nil
}

func (mastersUC *MastersUC) matchQualification(qualification string) (int64, error) {
	var qualifiactionInt int64 = 0
	if qualification != "" {
		qualifiactionInt = mastersUC.mastersConfig.qualificationMapBackwards[qualification]
		if qualifiactionInt != mastersUC.useCases.errorId {
			return qualifiactionInt, nil
		}
		badParamError := &models.BadQueryParameterError{Parameter: "qualification"}
		logger.Errorf(badParamError.Error())
		return qualifiactionInt, badParamError
	}
	return qualifiactionInt, nil
}

func (mastersUC *MastersUC) matchTheme(theme string, queryDB *models.MastersQueryValuesDB) error {
	if theme != "" {
		themeDB := models.ThemeDB{
			Name: theme,
		}
		err := mastersUC.ThemesRepo.GetThemeByName(&themeDB)
		if err != nil {
			return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
		}
		if themeDB.Id == mastersUC.useCases.errorId {
			badParamError := &models.BadQueryParameterError{Parameter: "theme"}
			logger.Errorf(badParamError.Error())
			return badParamError
		}
		queryDB.Theme = append(queryDB.Theme, themeDB.Id)
	}
	return nil
}

func (mastersUC *MastersUC) matchSubthemes(subthemes []string, queryDB *models.MastersQueryValuesDB) error {
	for _, subtheme := range subthemes {
		subthemeDB := models.SubthemeDB{Name: subtheme}
		err := mastersUC.ThemesRepo.GetSubthemeByName(&subthemeDB)
		if err != nil {
			return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
		}
		if subthemeDB.Id == mastersUC.useCases.errorId {
			badParamError := &models.BadQueryParameterError{Parameter: "subtheme"}
			logger.Errorf(badParamError.Error())
			return badParamError
		}
		queryDB.Subtheme = append(queryDB.Subtheme, subthemeDB.Id)
	}
	return nil
}

func (mastersUC *MastersUC) matchLanguages(languages []string, queryDB *models.MastersQueryValuesDB) error {
	for _, language := range languages {
		languageDB := models.LanguageDB{Name: language}
		err := mastersUC.LanguagesRepo.GetLanguageByName(&languageDB)
		if err != nil {
			return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
		}
		if languageDB.Id == mastersUC.useCases.errorId {
			badParamError := &models.BadQueryParameterError{Parameter: "language"}
			logger.Errorf(badParamError.Error())
			return badParamError
		}
		queryDB.Language = append(queryDB.Language, languageDB.Id)
	}
	return nil
}

func (mastersUC *MastersUC) matchMasterQuery(query *models.MastersQueryValues, queryDB *models.MastersQueryValuesDB) error {
	queryDB.Offset = query.Offset
	queryDB.Limit = query.Limit
	qualification, err := mastersUC.matchQualification(query.Qualification)
	if err != nil {
		return err
	}
	queryDB.Qualification = qualification
	educationFormat, err := mastersUC.matchEducationFormat(query.EducationFormat)
	if err != nil {
		return err
	}
	queryDB.EducationFormat = educationFormat
	queryDB.Language = make([]int64, 0)
	err = mastersUC.matchLanguages(query.Language, queryDB)
	queryDB.Theme = make([]int64, 0)
	err = mastersUC.matchTheme(query.Theme, queryDB)
	if err != nil {
		return err
	}
	queryDB.Subtheme = make([]int64, 0)
	err = mastersUC.matchSubthemes(query.Subtheme, queryDB)
	if err != nil {
		return err
	}

	if query.Search != "" {
		if query.Theme == "" {
			searchThemeIds, err := mastersUC.ThemesRepo.SearchThemeIds(query.Search)
			if err != nil {
				return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
			}
			queryDB.Theme = append(queryDB.Theme, searchThemeIds...)
			searchSubthemeIds, err := mastersUC.ThemesRepo.SearchSubthemeIdsOrThemes(query.Search, queryDB.Theme)
			if err != nil {
				return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
			}
			queryDB.Subtheme = append(queryDB.Subtheme, searchSubthemeIds...)
		} else {
			searchSubthemeIds, err := mastersUC.ThemesRepo.SearchSubthemeIdsAndThemes(query.Search, queryDB.Theme)
			if err != nil {
				return fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
			}
			queryDB.Subtheme = append(queryDB.Subtheme, searchSubthemeIds...)
		}

	}
	return nil
}

func (mastersUC *MastersUC) Get(query models.MastersQueryValues) (models.Masters, error) {
	var queryDB models.MastersQueryValuesDB
	masters := make([]models.Master, 0)
	err := mastersUC.matchMasterQuery(&query, &queryDB)
	if err != nil {
		return masters, err
	}
	if query.Search != "" && len(queryDB.Subtheme) == 0 {
		return masters, err
	}
	mastersDB, err := mastersUC.MastersRepo.GetMasters(queryDB)
	if err != nil {
		return masters, fmt.Errorf(mastersUC.useCases.errorMessages.DbError)
	}
	for _, masterDB := range mastersDB {
		master := models.Master{
			UserId: masterDB.UserId,
		}
		err = mastersUC.matchMaster(&masterDB, &master)
		if err != nil {
			return masters, err
		}
		masters = append(masters, master)
	}
	return masters, nil
}