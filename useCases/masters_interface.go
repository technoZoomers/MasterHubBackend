package useCases

import "github.com/technoZoomers/MasterHubBackend/models"

type MastersUCInterface interface {
	GetMasterById(master *models.Master) error
	ChangeMasterData(master *models.Master) error
}
