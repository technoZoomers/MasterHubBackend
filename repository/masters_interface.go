package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type MastersRepoI interface {
	GetMasterByUserId(master *models.MasterDB) error
	GetMasterIdByUsername(master *models.MasterDB) error
	GetMasterSubthemesById(masterId int64) ([]int64, error)
	DeleteMasterSubthemesById(masterId int64) error
	SetMasterSubthemesById(masterId int64, subthemes []int64) error
	GetMasterLanguagesById(masterId int64) ([]int64, error)
	DeleteMasterLanguagesById(masterId int64) error
	SetMasterLanguagesById(masterId int64, languages []int64) error
	UpdateMaster(master *models.MasterDB) error
}
