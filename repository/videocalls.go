package repository

import (
	"github.com/google/logger"
	"github.com/pion/webrtc/v2"
	"github.com/technoZoomers/MasterHubBackend/models"
)

type VideocallsRepo struct {
	repository  *Repository
	peerConnMap map[int64]chan *webrtc.Track
}

func (vcRepo *VideocallsRepo) AddNewConnection(peerConnection *models.PeerConnection, newTrack *webrtc.Track) {
	logger.Infof("added new connection with user:%d", peerConnection.UserId)
	_, ok := vcRepo.peerConnMap[peerConnection.UserId]
	if !ok {
		vcRepo.AddTrackCh(peerConnection.UserId)
	}
	vcRepo.peerConnMap[peerConnection.UserId] <- newTrack
}

func (vcRepo *VideocallsRepo) AddTrackCh(userId int64) {
	vcRepo.peerConnMap[userId] = make(chan *webrtc.Track, 1)
}

func (vcRepo *VideocallsRepo) GetTrack(peerConnection *models.PeerConnection) *webrtc.Track {
	_, ok := vcRepo.peerConnMap[peerConnection.PeerId]
	if !ok {
		vcRepo.AddTrackCh(peerConnection.PeerId)
	}
	return <-vcRepo.peerConnMap[peerConnection.PeerId]
}
