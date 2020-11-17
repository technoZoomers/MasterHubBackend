package handlers

import (
	"fmt"
	"github.com/google/logger"
	json "github.com/mailru/easyjson"
	"github.com/pion/webrtc/v2"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"net/http"
)

type VCHandlers struct {
	handlers     *Handlers
	videocallsUC useCases.VideocallsUCInterface
	webrtcAPI    *webrtc.API
	webrtcConfig webrtc.Configuration
}

func (vcHandlers *VCHandlers) createPeerConn(writer http.ResponseWriter, req *http.Request, creator bool) {
	var err error
	sent, userId := vcHandlers.handlers.validateUserId(writer, req)
	if sent {
		return
	}
	//sent = vcHandlers.handlers.checkUserAuth(writer, req, userId)
	//if sent {
	//	return
	//}
	sent, peerId := vcHandlers.handlers.validatePeerId(writer, req)
	if sent {
		return
	}
	var sdp models.Sdp
	err = json.UnmarshalFromReader(req.Body, &sdp)
	if err != nil {
		jsonError := fmt.Errorf("error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(jsonError.Error()))
		return
	}
	offer := webrtc.SessionDescription{}
	err = utils.Decode(sdp.Sdp, &offer)
	if err != nil {
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(err.Error()))
		return
	}

	//fmt.Println(sdp.Sdp)

	peerConnection, err := vcHandlers.webrtcAPI.NewPeerConnection(vcHandlers.webrtcConfig)
	if err != nil {
		jsonError := fmt.Errorf("error creating peer connection: %v", err.Error())
		logger.Errorf(jsonError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(jsonError.Error()))
		return
	}
	pConnection := models.PeerConnection{
		UserId:     userId,
		PeerId:     peerId,
		Connection: peerConnection,
		Sdp:        sdp,
	}

	if creator {
		vcHandlers.videocallsUC.AddTrack(&pConnection)
		fmt.Println("added new connection")
	} else {
		err = vcHandlers.videocallsUC.ConnectToTrack(&pConnection)
		if err != nil {
			vcHandlers.handlers.handleError(writer, err)
			return
		}
		fmt.Println("connected")
	}

	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		jsonError := fmt.Errorf("couldnt set remote description: %v", err.Error())
		logger.Errorf(jsonError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(jsonError.Error()))
		return
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		jsonError := fmt.Errorf("couldnt create answer: %v", err.Error())
		logger.Errorf(jsonError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(jsonError.Error()))
		return
	}

	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		jsonError := fmt.Errorf("couldnt set local description: %v", err.Error())
		logger.Errorf(jsonError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(jsonError.Error()))
		return
	}

	encodedSdp, err := utils.Encode(answer)
	if err != nil {
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(err.Error()))
		return
	}
	vcHandlers.answerSdp(writer, models.Sdp{Sdp: encodedSdp}, http.StatusOK, err)
}

func (vcHandlers *VCHandlers) Create(writer http.ResponseWriter, req *http.Request) {
	vcHandlers.createPeerConn(writer, req, true)
}

func (vcHandlers *VCHandlers) Connect(writer http.ResponseWriter, req *http.Request) {
	vcHandlers.createPeerConn(writer, req, false)
}

func (vcHandlers *VCHandlers) answerSdp(writer http.ResponseWriter, sdp models.Sdp, statusCode int, err error) {
	sent := vcHandlers.handlers.handleNotAcceptableError(writer, err)
	if !sent {
		utils.CreateAnswerSdpJson(writer, statusCode, sdp)
	}
}
