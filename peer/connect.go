package peer

import (
	"fmt"
	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/ice"
)

func NewConnectionFactory(config *webrtc.RTCConfiguration, table Table) ConnectionFactory {
	return &CoreConnectionFactory{
		config: config,
		table:  table,
	}
}

func (f *CoreConnectionFactory) Create(channel string) (conn *webrtc.RTCPeerConnection, err error) {
	if conn, err = webrtc.New(*f.config); nil != err {
		return
	}
	conn.OnICEConnectionStateChange = func(connectionState ice.ConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	}
	if init, ok := f.table[channel]; ok {
		if err = init(conn); nil != err {
			conn.Close()
			conn = nil
		}
	}
	return
}
