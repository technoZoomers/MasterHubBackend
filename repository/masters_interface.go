package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type MastersRepoI interface {
	GetMasterByUserId(master *models.MasterDB) error
	GetMasterByIdWithEmail(master *models.MasterDB) (string, error)
	GetMasterUserIdById(id int64) (int64, error)
	GetMasterIdByUsername(master *models.MasterDB) error
	GetMasterSubthemesById(masterId int64) ([]int64, error)
	DeleteMasterSubthemesById(masterId int64) error
	SetMasterSubthemesById(masterId int64, subthemes []int64) error
	UpdateMaster(master *models.MasterDB) error
	InsertMaster(master *models.MasterDB) error
	GetMasters(query models.MastersQueryValuesDB) ([]models.MasterDB, error)
}
