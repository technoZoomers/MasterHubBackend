package useCases

import (
	"github.com/technoZoomers/MasterHubBackend/models"
)

type MastersUCInterface interface {
	Get(query models.MastersQueryValues) (models.Masters, error)
	GetMasterById(master *models.Master) error
	ChangeMasterData(master *models.Master, masterId int64) error
	Register(masterFull *models.MasterFull) error
}
