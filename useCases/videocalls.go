package useCases

import (
	"fmt"
	"github.com/google/logger"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v2"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
	"io"
	"time"
)

type VideocallsUC struct {
	useCases        *UseCases
	VideocallsRepo  repository.VideocallsRepo // TODO: INTERFACE
	rtcpPLIInterval time.Duration
}

func (vcUC *VideocallsUC) ConnectToTrack(peerConnection *models.PeerConnection) error {
	track := vcUC.VideocallsRepo.GetTrack(peerConnection)
	senderAudio, err := peerConnection.Connection.AddTrack(track.AudioTrack)
	if err != nil {
		internalError := fmt.Errorf("couldnt add audio track: %s", err.Error())
		logger.Errorf(internalError.Error())
		return internalError
	}
	peerConnection.SenderAudio = senderAudio
	senderVideo, err := peerConnection.Connection.AddTrack(track.VideoTrack)
	if err != nil {
		internalError := fmt.Errorf("couldnt add video track: %s", err.Error())
		logger.Errorf(internalError.Error())
		return internalError
	}
	peerConnection.SenderVideo = senderVideo
	return nil
}

func (vcUC *VideocallsUC) AddTrack(peerConnection *models.PeerConnection) {
	_, err := peerConnection.Connection.AddTransceiver(webrtc.RTPCodecTypeVideo)
	if err != nil {
		internalError := fmt.Errorf("error adding video transceiver: %s", err.Error())
		logger.Errorf(internalError.Error())
	}

	_, err = peerConnection.Connection.AddTransceiver(webrtc.RTPCodecTypeAudio)
	if err != nil {
		internalError := fmt.Errorf("error adding audio transceiver: %s", err.Error())
		logger.Errorf(internalError.Error())
	}

	vcUC.VideocallsRepo.AddTrackToMap(peerConnection.UserId)

	peerConnection.Connection.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
		fmt.Println("on track")

		if remoteTrack.PayloadType() == webrtc.DefaultPayloadTypeVP8 || remoteTrack.PayloadType() == webrtc.DefaultPayloadTypeVP9 ||
			remoteTrack.PayloadType() == webrtc.DefaultPayloadTypeH264 {

			err = vcUC.sendVideoTrack(remoteTrack, peerConnection)
			if err != nil {
				return
			}

		} else {
			err = vcUC.sendAudioTrack(remoteTrack, peerConnection)
			if err != nil {
				return
			}
		}

	})
}

func (vcUC *VideocallsUC) sendVideoTrack(remoteTrack *webrtc.Track, peerConnection *models.PeerConnection) error {
	go func() {
		ticker := time.NewTicker(vcUC.rtcpPLIInterval)
		for range ticker.C {
			err := peerConnection.Connection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: remoteTrack.SSRC()}})
			if err != nil {
				internalError := fmt.Errorf("error sending data: %s", err.Error())
				logger.Errorf(internalError.Error())
			}
		}
	}()

	localTrack, err := peerConnection.Connection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), "video", "pion")
	if err != nil {
		internalError := fmt.Errorf("error creating new video track: %s", err.Error())
		logger.Errorf(internalError.Error())
		return internalError
	}
	newTrack, _ := vcUC.VideocallsRepo.GetTrackFromMap(peerConnection.UserId)
	newTrack.VideoTrack = localTrack
	if newTrack.AudioTrack != nil {
		vcUC.VideocallsRepo.AddNewConnection(peerConnection, newTrack)
	}

	go func() {
		//defer func() {
		//	err = peerConnection.Connection.RemoveTrack(peerConnection.Sender)
		//	if err != nil {
		//		internalError := fmt.Errorf("error reading: %s", err.Error())
		//		logger.Errorf(internalError.Error())
		//	}
		//	vcUC.VideocallsRepo.DeleteTrackCh(peerConnection)
		//	fmt.Println("removed track")
		//}()
		for {
			i, err := remoteTrack.ReadRTP()
			if err != nil {
				internalError := fmt.Errorf("error reading: %s", err.Error())
				logger.Errorf(internalError.Error())
				return
			}
			err = localTrack.WriteRTP(i)
			if err != nil && err != io.ErrClosedPipe {
				internalError := fmt.Errorf("error writing: %s", err.Error())
				logger.Errorf(internalError.Error())
				return
			}
		}
	}()
	return nil
}

func (vcUC *VideocallsUC) sendAudioTrack(remoteTrack *webrtc.Track, peerConnection *models.PeerConnection) error {

	localTrack, err := peerConnection.Connection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), "audio", "pion")
	if err != nil {
		internalError := fmt.Errorf("error creating new video track: %s", err.Error())
		logger.Errorf(internalError.Error())
		return internalError
	}
	newTrack, _ := vcUC.VideocallsRepo.GetTrackFromMap(peerConnection.UserId)
	newTrack.AudioTrack = localTrack
	if newTrack.VideoTrack != nil {
		vcUC.VideocallsRepo.AddNewConnection(peerConnection, newTrack)
	}

	vcUC.VideocallsRepo.AddNewConnection(peerConnection, newTrack)

	go func() {
		//defer func() {
		//	err = peerConnection.Connection.RemoveTrack(peerConnection.Sender)
		//	if err != nil {
		//		internalError := fmt.Errorf("error reading: %s", err.Error())
		//		logger.Errorf(internalError.Error())
		//	}
		//	vcUC.VideocallsRepo.DeleteTrackCh(peerConnection)
		//	fmt.Println("removed track")
		//}()
		for {
			i, err := remoteTrack.ReadRTP()
			if err != nil {
				internalError := fmt.Errorf("error reading: %s", err.Error())
				logger.Errorf(internalError.Error())
				return
			}
			err = localTrack.WriteRTP(i)
			if err != nil && err != io.ErrClosedPipe {
				internalError := fmt.Errorf("error writing: %s", err.Error())
				logger.Errorf(internalError.Error())
				return
			}
		}
	}()

	return nil
}
