package repository

import (
	"fmt"
	"github.com/pion/webrtc/v2"
	"github.com/technoZoomers/MasterHubBackend/models"
	"sync"
)

type VideocallsRepo struct {
	repository  *Repository
	peerConnMap map[int64]chan *Track
	tracksMap   map[int64]*Track
	isCalling   map[int64]int64
}

type Track struct {
	UserID         int64
	VideoTrack     *webrtc.Track
	AudioTrack     *webrtc.Track
	VideoTrackLock sync.RWMutex
	AudioTrackLock sync.RWMutex
}

func (vcRepo *VideocallsRepo) AddNewConnection(peerConnection *models.PeerConnection, track *Track) {
	fmt.Printf("added new connection with user:%d\n", peerConnection.UserId)
	_, ok := vcRepo.peerConnMap[peerConnection.UserId]
	if !ok {
		vcRepo.AddTrackCh(peerConnection.UserId)
	}
	vcRepo.peerConnMap[peerConnection.UserId] <- track
}

func (vcRepo *VideocallsRepo) AddTrackCh(userId int64) {
	vcRepo.peerConnMap[userId] = make(chan *Track, 1)
}

func (vcRepo *VideocallsRepo) AddTrackToMap(userId int64) {
	newTrack := Track{
		UserID: userId,
	}
	vcRepo.tracksMap[userId] = &newTrack
}

func (vcRepo *VideocallsRepo) GetTrackFromMap(userId int64) (*Track, bool) {
	track, ok := vcRepo.tracksMap[userId]
	return track, ok
}

func (vcRepo *VideocallsRepo) AddCallerState(userId int64, peerId int64) {
	vcRepo.isCalling[userId] = peerId
}

func (vcRepo *VideocallsRepo) RemoveCallerState(userId int64) {
	delete(vcRepo.isCalling, userId)
}

func (vcRepo *VideocallsRepo) GetCallerState(userId int64) int64 {
	return vcRepo.isCalling[userId]
}

func (vcRepo *VideocallsRepo) DeleteTrackFromMap(userId int64) {
	delete(vcRepo.tracksMap, userId)
}

func (vcRepo *VideocallsRepo) DeleteTrackCh(userId int64) {
	delete(vcRepo.peerConnMap, userId)

}

func (vcRepo *VideocallsRepo) GetTrack(peerConnection *models.PeerConnection) *Track {
	_, ok := vcRepo.peerConnMap[peerConnection.PeerId]
	if !ok {
		vcRepo.AddTrackCh(peerConnection.PeerId)
	}
	return <-vcRepo.peerConnMap[peerConnection.PeerId]
}
