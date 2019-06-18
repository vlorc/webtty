package peer

import (
	"fmt"
	"github.com/pion/webrtc/v2"
)

func NewConnectionFactory(config *webrtc.Configuration, table Table) ConnectionFactory {
	media := webrtc.MediaEngine{}
	media.RegisterDefaultCodecs()
	api := webrtc.NewAPI(webrtc.WithMediaEngine(media))
	return NewConnectionFactoryWithApi(config, table, api)
}

func NewConnectionFactoryWithApi(config *webrtc.Configuration, table Table, api *webrtc.API) ConnectionFactory {
	return &CoreConnectionFactory{
		config: config,
		table:  table,
		api:    api,
	}
}

func (f *CoreConnectionFactory) Create(channel string) (conn *webrtc.PeerConnection, err error) {
	if conn, err = f.api.NewPeerConnection(*f.config); nil != err {
		return
	}
	conn.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})
	if init, ok := f.table[channel]; ok {
		if err = init(conn); nil != err {
			conn.Close()
			conn = nil
		}
	}
	return
}
