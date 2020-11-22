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
	sender, err := peerConnection.Connection.AddTrack(vcUC.VideocallsRepo.GetTrack(peerConnection))
	if err != nil {
		internalError := fmt.Errorf("couldnt add track: %s", err.Error())
		logger.Errorf(internalError.Error())
		return internalError
	}
	peerConnection.Sender = sender
	return nil
}

func (vcUC *VideocallsUC) AddTrack(peerConnection *models.PeerConnection) {
	_, err := peerConnection.Connection.AddTransceiver(webrtc.RTPCodecTypeVideo)
	if err != nil {
		internalError := fmt.Errorf("error adding transceiver: %s", err.Error())
		logger.Errorf(internalError.Error())
	}

	peerConnection.Connection.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
		fmt.Println("on track")
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
			internalError := fmt.Errorf("error creating new track: %s", err.Error())
			logger.Errorf(internalError.Error())
		}

		vcUC.VideocallsRepo.AddNewConnection(peerConnection, localTrack)

		//_, err = peerConnection.Connection.AddTransceiverFromTrack(localTrack)
		//if err != nil {
		//	internalError := fmt.Errorf("error adding transceiver: %s", err.Error())
		//	logger.Errorf(internalError.Error())
		//}
		go func() {
			defer func() {
				err = peerConnection.Connection.RemoveTrack(peerConnection.Sender)
				if err != nil {
					internalError := fmt.Errorf("error reading: %s", err.Error())
					logger.Errorf(internalError.Error())
				}
				vcUC.VideocallsRepo.DeleteTrackCh(peerConnection)
			}()
			rtpBuf := make([]byte, 1400)
			for {
				i, err := remoteTrack.Read(rtpBuf)
				if err != nil {
					internalError := fmt.Errorf("error reading: %s", err.Error())
					logger.Errorf(internalError.Error())
				}
				_, err = localTrack.Write(rtpBuf[:i])
				if err != nil && err != io.ErrClosedPipe {
					internalError := fmt.Errorf("error writing: %s", err.Error())
					logger.Errorf(internalError.Error())
				}
			}
		}()

	})
}
