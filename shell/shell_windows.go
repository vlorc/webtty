package shell

import (
	"log"

	"github.com/iamacarpet/go-winpty"
)

type winPtr struct {
	pty *winpty.WinPTY
}

func __create(config *Config) Pty {
	pty, err := winpty.OpenWithOptions(winpty.Options{
		Env: config.Env,
		Dir: config.Directory,
		Command:   config.Command,
	})
	if err != nil {
		log.Printf("Failed open from pty master: %s\n", err)
		panic(err)
	}

	return winPtr{
		pty: pty,
	}
}

func (p winPtr) Read(b []byte) (int, error) {
	return p.pty.StdOut.Read(b)
}

func (p winPtr) Write(b []byte) (int, error) {
	return p.pty.StdIn.Write(b)
}

func (p winPtr) Close() error {
	p.pty.Close()
	return nil
}

func (p winPtr) SetSize(r, c int) error {
	p.pty.SetSize(uint32(c), uint32(r))
	return nil
}
