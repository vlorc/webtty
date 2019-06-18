package peer

import (
	"github.com/pion/webrtc/v2"
	"math/rand"
	"sync"
)

type Event interface {
	With(...interface{}) Event
	On(string, ...interface{}) Event
	Emit(Driver, []byte) error
}

type Driver interface {
	With(...interface{}) Driver
	On(string, ...interface{}) Driver
	Emit(string, ...interface{}) Driver
	EmitTo(string, string, interface{}) Driver
}

type ConnectionFactory interface {
	Create(types string) (*webrtc.PeerConnection, error)
}

type SessionFactory interface {
	Id() string
	Create(id string, conn *webrtc.PeerConnection) *Session
	Attach(id string, conn *webrtc.PeerConnection) *Session
	Query(id string) *Session
	Remove(id string)
}

type Message struct {
	Command string      `json:"command"`
	Dest    string      `json:"dest,omitempty"`
	Source  string      `json:"source,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

type Description struct {
	Session string `json:"session"`
	Channel string `json:"channel,omitempty"`
	webrtc.SessionDescription
}

type Session struct {
	*webrtc.PeerConnection
	id    string
	state chan int
	done  chan error
}

type CoreSessionFactory struct {
	id     string
	table  sync.Map
	source rand.Source
}

type Table map[string]func(*webrtc.PeerConnection) error

type CoreConnectionFactory struct {
	config *webrtc.Configuration
	table  Table
	api *webrtc.API
}
