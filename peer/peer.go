package peer

import (
	"github.com/pion/webrtc/v2"
	"log"
	"sync"
)

type Peer struct {
	driver  Driver
	factory ConnectionFactory
	session SessionFactory
	mutex   sync.Mutex
}

func NewPeer(factory ConnectionFactory, session SessionFactory, driver Driver) *Peer {
	peer := &Peer{
		driver:  driver,
		factory: factory,
		session: session,
	}
	driver.With(peer)
	return peer
}

func (p *Peer) OnClose(id string) (name string, result string, err error) {
	if conn := p.session.Query(id); nil != conn {
		p.session.Remove(id)
		name = "CLOSE"
		result = id
		err = conn.Close()
	}
	return
}

func (p *Peer) OnOffer(offer *Description) (name string, answer *Description, err error) {
	peer, err := p.__create(offer.Session, offer.Channel)
	if nil != err {
		return
	}
	if err = peer.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: offer.SDP}); nil != err {
		log.Printf("SetRemoteDescription error: %s\n", err)
		return
	}
	sdp, err := peer.CreateAnswer(nil)
	if nil == err {
		name = "ANSWER"
		answer = &Description{
			Session:               offer.Session,
			Channel:               offer.Channel,
			SessionDescription: sdp,
		}
	}
	return
}

func (p *Peer) OnAnswer(answer *Description) error {
	if conn := p.session.Query(answer.Session); nil != conn {
		conn.state <- 1
		return conn.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeAnswer, SDP: answer.SDP})
	}
	return nil
}

func (p *Peer) __create(id string, channel string) (session *Session, err error) {
	if session = p.session.Query(id); nil != session {
		return
	}
	conn, err := p.factory.Create(channel)
	if nil == err {
		session = p.session.Attach(id, conn)
	}
	return
}

func (p *Peer) Connect(id string, channel string, init ...func(conn *webrtc.PeerConnection) error) (*Session, error) {
	conn, err := p.factory.Create(channel)
	if nil != err {
		return nil, err
	}
	for _, v := range init {
		v(conn)
	}
	offer, err := conn.CreateOffer(nil)
	if nil != err {
		conn.Close()
		log.Printf("CreateOffer error: %s\n", err)
		return nil, err
	}

	session := p.session.Create(id, conn)
	p.driver.EmitTo(id, "OFFER", &Description{
		Session:            session.id,
		Channel:            channel,
		SessionDescription: offer,
	})
	return session, nil
}
