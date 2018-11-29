package peer

import (
	"fmt"
	"github.com/pions/webrtc"
	"math/rand"
	"time"
)

func NewSessionFactory(id string) SessionFactory {
	return &CoreSessionFactory{
		id:     id,
		source: rand.NewSource(time.Now().Unix()),
	}
}

func (f *CoreSessionFactory) Id() string {
	return f.id
}

func (f *CoreSessionFactory) __id(id string) string {
	return fmt.Sprintf("%s.%s.%d", f.id, id, f.source.Int63())
}

func (f *CoreSessionFactory) Query(id string) *Session {
	it, ok := f.table.Load(id)
	if ok {
		return it.(*Session)
	}
	return nil
}

func (f *CoreSessionFactory) Remove(id string) {
	f.table.Delete(id)
}

func (f *CoreSessionFactory) Attach(id string, conn *webrtc.RTCPeerConnection) *Session {
	session := NewSession(id, conn)
	f.table.Store(id, session)
	return session
}

func (f *CoreSessionFactory) Create(id string, conn *webrtc.RTCPeerConnection) *Session {
	return f.Attach(f.__id(id), conn)
}

func (s *Session) Close() error {
	err := s.RTCPeerConnection.Close()
	close(s.done)
	return err
}

func NewSession(id string, conn *webrtc.RTCPeerConnection) *Session {
	return &Session{
		RTCPeerConnection: conn,
		id:                id,
		state:             make(chan int, 1),
		done:              make(chan error, 1),
	}
}
func (s *Session) Id() string {
	return s.id
}

func (s *Session) Done() error {
	return <-s.done
}

func (s *Session) State() int {
	return <-s.state
}
