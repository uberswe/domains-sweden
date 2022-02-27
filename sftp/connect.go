package sftp

import (
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
	"log"
	"net"
	"os"
)

func (s *Service) connect() error {
	// We check if there is an active connection, if there is we return early.
	if s.client != nil {
		_, err := s.client.ReadDir("/")
		if err == nil {
			return nil
		}
	}

	port := 22
	var err error
	log.Printf("Connecting to %s ...\n", s.Host)

	var auths []ssh.AuthMethod

	// Try to use $SSH_AUTH_SOCK which contains the path of the unix file socket that the sshd agent uses
	// for communication with other processes.
	if aconn, err2 := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err2 == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))
	}

	// Use password authentication if provided
	if s.Password != "" {
		auths = append(auths, ssh.Password(s.Password))
	}

	var keyErr *knownhosts.KeyError

	// Initialize client configuration
	config := ssh.ClientConfig{
		User: s.Username,
		Auth: auths,
		HostKeyCallback: ssh.HostKeyCallback(func(host string, remote net.Addr, pubKey ssh.PublicKey) error {
			found, hErr := s.checkKnownHost(host, remote, pubKey, "")
			// Reference: https://blog.golang.org/go1.13-errors
			// To understand what errors.As is.
			if errors.As(hErr, &keyErr) && len(keyErr.Want) > 0 {
				// Reference: https://www.godoc.org/golang.org/x/crypto/ssh/knownhosts#KeyError
				// if keyErr.Want slice is empty then host is unknown, if keyErr.Want is not empty
				// and if host is known then there is key mismatch the connection is then rejected.
				log.Printf("WARNING: %v is not a key of %s, either a MiTM attack or %s has reconfigured the host pub key.", string(pubKey.Marshal()), host, host)
				return keyErr
			} else if errors.As(hErr, &keyErr) && len(keyErr.Want) == 0 {
				// host key not found in known_hosts then give a warning and continue to connect.
				log.Printf("WARNING: %s is not trusted, adding this key: %q to known_hosts file.", host, string(pubKey.Marshal()))
				return s.addKnownHost(remote, pubKey, "")
			} else if !found {
				log.Printf("WARNING: %s is not found, adding this key: %q to known_hosts file.", host, string(pubKey.Marshal()))
				return s.addKnownHost(remote, pubKey, "")
			}
			log.Printf("Pub key exists for %s.", host)
			return nil
		}),
	}

	addr := fmt.Sprintf("%s:%d", s.Host, port)

	// Connect to server
	s.sshClient, err = ssh.Dial("tcp", addr, &config)
	if err != nil {
		log.Printf("Failed to connect to [%s]: %v\n", addr, err)
		return err
	}

	// Create new SFTP client
	s.client, err = sftp.NewClient(s.sshClient)
	if err != nil {
		log.Printf("Unable to start SFTP subsystem: %v\n", err)
		return err
	}
	return nil
}
