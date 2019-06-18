package main

import (
	"flag"
	"rtclient/peer"
	"rtclient/shell"

	"github.com/pion/webrtc/v2"
)

func main() {
	id := flag.String("id", "shell", "client id")
	cmd := flag.String("cmd", "cmd", "command to execute")
	gateway := flag.String("gateway", "", "gateway url")
	stun := flag.String("stun", "", "stun url")
	flag.Parse()

	peer.NewPeer(
		peer.NewConnectionFactory(
			&webrtc.Configuration{
				ICEServers: []webrtc.ICEServer{
					{
						URLs:       []string{*stun},
					},
				},
			},
			peer.Table{
				"shell": shell.Shell(shell.Command(*cmd)),
			}),
		peer.NewSessionFactory(*id),
		peer.NewWebSocketDriver(*id, *gateway+"/"+*id, peer.NewEvent()),
	)

	select {}
}
