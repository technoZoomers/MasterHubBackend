package repository

import (
	"github.com/pion/webrtc/v2"
	"github.com/technoZoomers/MasterHubBackend/models"
)

type VideocallsRepo struct {
	repository  *Repository
	peerConnMap map[int64]chan *webrtc.Track
}

func (vcRepo *VideocallsRepo) AddNewConnection(peerConnection *models.PeerConnection, newTrack *webrtc.Track) {
	vcRepo.peerConnMap[peerConnection.UserId] = make(chan *webrtc.Track, 1)
	vcRepo.peerConnMap[peerConnection.UserId] <- newTrack
}

func (vcRepo *VideocallsRepo) GetTrack(peerConnection *models.PeerConnection) *webrtc.Track {
	return <-vcRepo.peerConnMap[peerConnection.PeerId]
}
