package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type MastersRepoI interface {
	GetMasterByUserId(master *models.MasterDB) (int64, error)
	GetMasterSubthemesById(masterId int64) ([]int64, error)
	DeleteMasterSubthemesById(masterId int64) error
	SetMasterSubthemesById(masterId int64, subthemes []int64) error
	GetMasterLanguagesById(masterId int64) ([]int64, error)
}

