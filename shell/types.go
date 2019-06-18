package shell

import "io"

type Pty interface {
	io.ReadWriteCloser
	SetSize(int, int) error
}

type message struct {
	Type string    `json:"type"`
	Data []float64 `json:"data"`
}

type Factory func(*Config) Pty