package shell

import "github.com/pion/webrtc/v2"

func Command(command string) Factory{
	return Default(&Config{Command: command})
}

func Custom(config *Config) func(*Config) Pty {
	return func(*Config) Pty {
		return __create(config)
	}
}

func Default(config *Config) func(*Config) Pty {
	return func(c *Config) Pty {
		if nil == c{
			c = config
		}
		return __create(c)
	}
}

func Shell(factory Factory) func(*webrtc.PeerConnection) error {
	return shell(factory)
}

func shell(factory Factory) func(*webrtc.PeerConnection) error {
	return func(conn *webrtc.PeerConnection) error {
		return __shell(conn, factory)
	}
}
