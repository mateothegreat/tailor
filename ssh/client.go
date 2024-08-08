package ssh

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/mateothegreat/tailer/util"
	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	*ssh.ClientConfig
	Name         string
	Address      string
	Port         int
	IdentityFile string
}

type HostConfig struct {
	Hostname string
}

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

func GetConfigs(path string, servers []HostConfig) (map[string]*Session, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	config, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	configs := make(map[string]*Session)

	blocks := regexp.MustCompile(`(?im)^\s*host\s+(\S+)((?:\n\s{4}\S.*)+)?`).FindAllStringSubmatch(string(config), -1)
	if len(blocks) == 0 {
		return nil, fmt.Errorf("hostname not found")
	}

	for _, block := range blocks {
		name := regexp.MustCompile(`(?im)^host\s+(\S+)`).FindStringSubmatch(block[0])
		if len(name) != 2 {
			continue
		}

		for _, server := range servers {
			if server.Hostname == name[1] {
				session := &Session{
					Config: &SSHConfig{
						ClientConfig: &ssh.ClientConfig{
							HostKeyCallback: ssh.InsecureIgnoreHostKey(),
							Timeout:         10 * time.Second,
						},
						Name: name[1],
					},
				}

				host := regexp.MustCompile(`(?im)hostname\s+(\S+)`).FindStringSubmatch(block[0])
				if len(host) == 2 {
					session.Config.Address = host[1]
				}

				port := regexp.MustCompile(`(?im)port\s+(\S+)`).FindStringSubmatch(block[0])
				if len(port) == 2 {
					p, err := strconv.Atoi(port[1])
					if err != nil {
						return nil, err
					}
					session.Config.Port = p
				} else {
					session.Config.Port = 22
				}

				user := regexp.MustCompile(`(?im)user\s+(\S+)`).FindStringSubmatch(block[0])
				if len(user) == 2 {
					session.Config.User = user[1]
				}

				identityFile := regexp.MustCompile(`(?im)identityfile\s+(\S+)`).FindStringSubmatch(block[0])
				if len(identityFile) == 2 {
					session.Config.IdentityFile = identityFile[1]
				}

				configs[server.Hostname] = session

				if session.Config.IdentityFile != "" {
					identityFile, err := os.ReadFile(util.ExpandPath(session.Config.IdentityFile))
					if err != nil {
						return nil, err
					}
					session.Config.IdentityFile = string(identityFile)

					signer, err := ssh.ParsePrivateKey([]byte(session.Config.IdentityFile))
					if err != nil {
						return nil, err
					}

					session.Config.ClientConfig.Auth = append(session.Config.ClientConfig.Auth, ssh.PublicKeys(signer))
				}
			}
		}
	}

	return configs, nil
}
