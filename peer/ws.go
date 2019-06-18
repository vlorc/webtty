package peer

import (
	"github.com/json-iterator/go"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketDriver struct {
	id      string
	event   Event
	conn    *websocket.Conn
	url     string
	state   bool
	timeout time.Duration
	write   chan *Message
}

func NewWebSocketDriver(id, url string, event Event) Driver {
	driver := &WebSocketDriver{
		id:      id,
		event:   event,
		url:     url,
		state:   true,
		timeout: 1 * time.Minute,
		write:   make(chan *Message, 256),
	}

	go driver.__loop()
	go driver.__write()
	return driver
}

func __write(conn *websocket.Conn, msg *Message) error {
	w, err := conn.NextWriter(websocket.TextMessage)
	if err == nil {
		jsoniter.ConfigFastest.NewEncoder(w).Encode(msg)
		err = w.Close()
	}
	return err
}

func (ws *WebSocketDriver) __write() {
	ticker := time.NewTicker(ws.timeout - time.Second*5)
	defer ticker.Stop()

	var buf []*Message
	for {
		select {
		case msg, ok := <-ws.write:
			if !ok {
				// ws.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			if nil == ws.conn {
				buf = append(buf, msg)
				break
			}
			ws.conn.SetWriteDeadline(time.Now().Add(time.Second * 15))
			__write(ws.conn, msg)
			for _, b := range buf {
				__write(ws.conn, b)
			}
			buf = nil
			ws.conn.SetWriteDeadline(time.Time{})
		case <-ticker.C:
			if nil != ws.conn {
				ws.conn.SetWriteDeadline(time.Now().Add(time.Second * 5))
				ws.conn.WriteMessage(websocket.PingMessage, nil)
				ws.conn.SetWriteDeadline(time.Time{})
			}
		}
	}
}

func (ws *WebSocketDriver) __read(conn *websocket.Conn) {
	conn.SetReadDeadline(time.Now().Add(ws.timeout))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(ws.timeout))
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Websocket readMessage error:", err.Error())
			break
		}
		ws.event.Emit(ws, message)
	}
}

func (ws *WebSocketDriver) __loop() {
	for i := 0; ws.state; i++ {
		conn, _, err := websocket.DefaultDialer.Dial(ws.url, nil)
		if nil != err {
			log.Println("Websocket dial error:", err.Error(), "try:", i)
			time.Sleep(time.Duration(5+(i&15)) + time.Second)
			continue
		}
		ws.conn = conn
		ws.__read(conn)
		ws.Close()
	}
}

func (ws *WebSocketDriver) With(obj ...interface{}) Driver {
	ws.event.With(obj...)
	return ws
}

func (ws *WebSocketDriver) On(command string, fn ...interface{}) Driver {
	ws.event.On(command, fn...)
	return ws
}

func (ws *WebSocketDriver) Emit(command string, val ...interface{}) Driver {
	if nil != ws.conn {
		ws.write <- &Message{
			Command: command,
			Payload: val[0],
		}
	}
	return ws
}

func (ws *WebSocketDriver) Close() error {
	conn := ws.conn
	write := ws.write
	ws.conn = nil
	ws.state = false
	ws.write = nil
	if nil != conn {
		return conn.Close()
	}
	if nil != ws.write {
		close(write)
	}
	return nil
}

func (ws *WebSocketDriver) EmitTo(id, command string, val interface{}) Driver {
	if nil != ws.conn {
		ws.write <- &Message{
			Command: command,
			Source:  ws.id,
			Dest:    id,
			Payload: val,
		}
	}
	return ws
}
