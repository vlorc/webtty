package peer

import (
	"io"
	"time"
)

type ShellDriver struct {
	id      string
	event   Event
	fd      io.ReadWriteCloser
	url     string
	state   bool
	timeout time.Duration
	write   chan *Message
}

func NewShellDriver(id, url string, event Event) Driver {
	driver := &ShellDriver{
		id:      id,
		event:   event,
		url:     url,
		state:   true,
		timeout: 1 * time.Minute,
		write:   make(chan *Message, 256),
	}

	return driver
}


func (sh *ShellDriver) With(obj ...interface{}) Driver {
	sh.event.With(obj...)
	return sh
}

func (sh *ShellDriver) On(command string, fn ...interface{}) Driver {
	sh.event.On(command, fn...)
	return sh
}

func (sh *ShellDriver) Emit(command string, val ...interface{}) Driver {
	if nil != sh.fd {
		sh.write <- &Message{
			Command: command,
			Payload: val[0],
		}
	}
	return sh
}

func (sh *ShellDriver) Close() error {
	conn := sh.fd
	write := sh.write
	sh.fd = nil
	sh.state = false
	sh.write = nil
	if nil != conn {
		return conn.Close()
	}
	if nil != sh.write {
		close(write)
	}
	return nil
}

func (sh *ShellDriver) EmitTo(id, command string, val interface{}) Driver {
	if nil != sh.fd {
		sh.write <- &Message{
			Command: command,
			Source:  sh.id,
			Dest:    id,
			Payload: val,
		}
	}
	return sh
}
