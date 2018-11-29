package main

import (
	"flag"
	"github.com/pions/webrtc"
	"github.com/vlorc/webtty/peer"
	"github.com/vlorc/webtty/shell"
)

func main() {
	id := flag.String("id", "shell", "client id")
	cmd := flag.String("cmd", "cmd", "command to execute")
	gateway := flag.String("gateway", "", "gateway url")
	stun := flag.String("stun", "", "stun url")
	flag.Parse()

	webrtc.RegisterDefaultCodecs()
	peer.NewPeer(
		peer.NewConnectionFactory(
			&webrtc.RTCConfiguration{
				IceServers: []webrtc.RTCIceServer{
					{
						URLs:       []string{*stun},
					},
				},
			},
			peer.Table{
				"shell": shell.Shell(*cmd),
			}),
		peer.NewSessionFactory(*id),
		peer.NewWebSocketDriver(*id, *gateway+"/"+*id, peer.NewEvent()),
	)

	select {}
}
