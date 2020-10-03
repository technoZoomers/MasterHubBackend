package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type MastersRepoI interface {
	GetMasterByUserId(master *models.MasterDB) (int64, error)
	GetMasterSubthemesById(master *models.MasterDB) ([]int64, error)
	GetMasterLanguagesById(master *models.MasterDB) ([]int64, error)
}

