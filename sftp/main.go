package sftp

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"os"
)

// Default known hosts path.
var defaultPath = os.ExpandEnv("$HOME/.ssh/known_hosts")

type Service struct {
	Host      string
	Username  string
	Password  string
	sshClient *ssh.Client
	client    *sftp.Client
}

func New(host, username, password string) *Service {
	s := Service{
		Host:     host,
		Username: username,
		Password: password,
	}
	s.connect()
	return &s
}

func (s *Service) Close() error {
	err := s.sshClient.Close()
	if err != nil {
		return err
	}
	err = s.client.Close()
	if err != nil {
		return err
	}
	return nil
}
