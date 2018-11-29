package shell

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"unicode/utf8"

	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/datachannel"
	"github.com/pions/webrtc/pkg/ice"
)

func Shell(cmd string) func(*webrtc.RTCPeerConnection) error {
	return shell(cmd)
}

func create(cmd string) func() Pty {
	return func() Pty {
		return __create(cmd)
	}
}

func shell(cmd string) func(*webrtc.RTCPeerConnection) error {
	get := create(cmd)
	return func(conn *webrtc.RTCPeerConnection) error {
		return __shell(conn, get)
	}
}

func __data(d *webrtc.RTCDataChannel, conn *webrtc.RTCPeerConnection, pty Pty) {
	d.OnOpen = func() {
		buf := make([]byte, 8192)
		reader := bufio.NewReader(pty)
		var buffer bytes.Buffer
		for {
			n, err := reader.Read(buf)
			if err != nil {
				log.Printf("Failed to read from pty master: %s\n", err)
				return
			}
			bufferBytes := buffer.Bytes()
			runeReader := bufio.NewReader(bytes.NewReader(append(bufferBytes[:], buf[:n]...)))
			buffer.Reset()
			i := 0
			for i < n {
				char, charLen, e := runeReader.ReadRune()
				if e != nil {
					log.Printf("Failed to read from pty master: %s\n", err)
					return
				}
				if char == utf8.RuneError {
					runeReader.UnreadRune()
					break
				}
				i += charLen
				buffer.WriteRune(char)
			}
			d.Send(datachannel.PayloadString{Data: buffer.Bytes()})
			if err != nil {
				log.Printf("Failed to send UTF8 char: %s\n", err)
				return
			}
			buffer.Reset()
			if i < n {
				buffer.Write(buf[i:n])
			}
		}
	}
	d.Onmessage = func(payload datachannel.Payload) {
		switch p := payload.(type) {
		case *datachannel.PayloadString:
			pty.Write(p.Data)
		case *datachannel.PayloadBinary:
			pty.Write(p.Data)
		default:
			log.Printf("Message '%s' from DataChannel '%s' no payload \n", p.PayloadType().String(), d.Label)
		}
	}
	d.OnMessage = d.Onmessage
}

func __control(d *webrtc.RTCDataChannel, conn *webrtc.RTCPeerConnection, pty Pty) {
	d.OnOpen = func() {
		log.Printf("Open from DataChannel '%s'\n", d.Label)
	}
	d.Onmessage = func(payload datachannel.Payload) {
		msg := &message{}
		switch p := payload.(type) {
		case *datachannel.PayloadString:
			json.Unmarshal(p.Data, msg)
		case *datachannel.PayloadBinary:
			json.Unmarshal(p.Data, msg)
		default:
			log.Printf("Message '%s' from DataChannel '%s' no payload \n", p.PayloadType().String(), d.Label)
		}
		if "resize" == msg.Type {
			log.Printf("Resize from DataChannel '%s' %v:%v \n", d.Label, msg.Data[1], msg.Data[0])
			pty.SetSize(int(msg.Data[1]), int(msg.Data[0]))
		}
	}
	d.OnMessage = d.Onmessage
}

func __shell(conn *webrtc.RTCPeerConnection, get func() Pty) error {
	var pty Pty
	conn.OnICEConnectionStateChange = func(state ice.ConnectionState) {
		if ice.ConnectionStateConnected == state {
			pty = get()
			pty.SetSize(60, 200)
		} else if state > ice.ConnectionStateConnected {
			if nil != pty {
				pty.Close()
				pty = nil
			}
			log.Printf("ICE Connection State has changed: %s\n", state.String())
		}
	}
	conn.OnDataChannel = func(d *webrtc.RTCDataChannel) {
		log.Printf("New DataChannel %s %d\n", d.Label, d.ID)

		d.Lock()
		defer d.Unlock()

		if "data" == d.Label {
			__data(d, conn, pty)
		} else if "control" == d.Label {
			__control(d, conn, pty)
		}
	}
	return nil
}
