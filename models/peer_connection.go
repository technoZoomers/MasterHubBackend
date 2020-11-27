package models

import "github.com/pion/webrtc/v2"

type PeerConnection struct {
	UserId      int64
	PeerId      int64
	Sdp         Sdp
	SenderAudio *webrtc.RTPSender
	SenderVideo *webrtc.RTPSender
	Connection  *webrtc.PeerConnection
}
