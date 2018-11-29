package shell

import (
	"log"
	"os"
	"os/exec"
	"github.com/kr/pty"
)

type linuxPtr struct {
	file os.File
}

func __create(cmd string) Ptr {
	c := exec.Command(cmd)
	f, err := pty.Start(c)
	if err != nil {
		log.Printf("Failed open from pty master: %s\n", err)
		panic(err)
	}
	return linuxPtr{
		file: f
	}
}

func(p linuxPtr)Read(b []byte)(int,error){
	return p.file.Read(b)
}

func(p linuxPtr)Write(b []byte)(int,error){
	return p.file.Write(b)
}

func(p linuxPtr)Close()error{
	return p.file.Close()
}

func(p linuxPtr)SetSize(r,c int) error{
	return pty.Setsize(p.file,&pty.Winsize{
		Rows: uint16(r),
		Cols: uint16(c),
	})
}