package ssh

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type Manager struct {
	Sessions map[string]*Session
}

type Session struct {
	Config *SSHConfig
	Client *ssh.Client
}

func (s *Session) Connect() error {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.Config.Address, s.Config.Port), s.Config.ClientConfig)
	if err != nil {
		return err
	}
	s.Client = client
	return nil
}

func (s *Session) Close() error {
	return s.Client.Close()
}

func NewManager(hosts []HostConfig) (*Manager, error) {
	configs, err := GetConfigs(os.Getenv("HOME")+"/.ssh/config", hosts)
	if err != nil {
		return nil, err
	}
	return &Manager{Sessions: configs}, nil
}
