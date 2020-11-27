package useCases

import "github.com/technoZoomers/MasterHubBackend/models"

type VideocallsUCInterface interface {
	AddTrack(peerConnection *models.PeerConnection)
	ConnectToTrack(peerConnection *models.PeerConnection) error
	CheckIsCalling(peerId int64) int64
}
