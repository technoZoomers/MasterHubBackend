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

//func (vcRepo *VideocallsRepo) DeleteTrackCh(peerConnection *models.PeerConnection) {
//	_, ok := vcRepo.peerConnMap[peerConnection.PeerId]
//	if ok {
//		delete(vcRepo.peerConnMap, peerConnection.PeerId)
//	}
//}

func (vcRepo *VideocallsRepo) GetTrack(peerConnection *models.PeerConnection) *Track {
	_, ok := vcRepo.peerConnMap[peerConnection.PeerId]
	if !ok {
		vcRepo.AddTrackCh(peerConnection.PeerId)
	}
	return <-vcRepo.peerConnMap[peerConnection.PeerId]
}
