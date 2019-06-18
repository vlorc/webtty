package shell

import (
	"bufio"
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"github.com/pion/ice"
	"log"
	"unicode/utf8"

	"github.com/pion/webrtc/v2"
)

func __data(d *webrtc.DataChannel, conn *webrtc.PeerConnection, pty Pty) {
	d.OnOpen(func() {
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
			d.SendText(string(buffer.Bytes()))
			if err != nil {
				log.Printf("Failed to send UTF8 char: %s\n", err)
				return
			}
			buffer.Reset()
			if i < n {
				buffer.Write(buf[i:n])
			}
		}
	})

	d.OnMessage(func(payload webrtc.DataChannelMessage) {
		pty.Write(payload.Data)
	})
}

func __control(d *webrtc.DataChannel, conn *webrtc.PeerConnection, pty Pty) {
	d.OnOpen(func() {
		log.Printf("Open from DataChannel '%s'\n", d.Label())
	})
	d.OnMessage((func(payload webrtc.DataChannelMessage) {
		msg := &message{}
		jsoniter.Unmarshal(payload.Data, msg)
		if "resize" == msg.Type {
			log.Printf("Resize from DataChannel '%s' %v:%v \n", d.Label(), msg.Data[1], msg.Data[0])
			pty.SetSize(int(msg.Data[1]), int(msg.Data[0]))
		}
	}))
}

func __shell(conn *webrtc.PeerConnection, factory Factory) error {
	var pty Pty = factory(nil) //use default config
	pty.SetSize(60, 200)
	conn.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		if state > ice.ConnectionStateConnected {
			if nil != pty {
				pty.Close()
			}
		}
		log.Printf("ICE Connection State has changed: %s\n", state.String())
	})
	conn.OnDataChannel(func(d *webrtc.DataChannel) {
		log.Printf("New DataChannel %s %d\n", d.Label(), d.ID())
		if label := d.Label(); "data" == label {
			__data(d, conn, pty)
		} else if "control" == label {
			__control(d, conn, pty)
		}
	})
	return nil
}
