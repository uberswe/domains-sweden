package sftp

import (
	"errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"net"
	"os"
	"path/filepath"
)

// knownHosts returns host key callback from a custom known hosts path.
func (s *Service) knownHosts(file string) (ssh.HostKeyCallback, error) {
	return knownhosts.New(file)
}

// checkKnownHost checks is host in known hosts file.
// it returns is the host found in known_hosts file and error, if the host found in
// known_hosts file and error not nil that means public key mismatch, maybe MAN IN THE MIDDLE ATTACK! you should not handshake.
func (s *Service) checkKnownHost(host string, remote net.Addr, key ssh.PublicKey, knownFile string) (found bool, err error) {

	var keyErr *knownhosts.KeyError

	// Fallback to default known_hosts file
	if knownFile == "" {
		knownFile = defaultPath
	}

	// Get host key callback
	callback, err := s.knownHosts(knownFile)

	if err != nil {
		return false, err
	}

	// check if host already exists.
	err = callback(host, remote, key)

	// Known host already exists.
	if err == nil {
		return true, nil
	}

	// Make sure that the error returned from the callback is host not in file error.
	// If keyErr.Want is greater than 0 length, that means host is in file with different key.
	if errors.As(err, &keyErr) && len(keyErr.Want) > 0 {
		return true, keyErr
	}

	// Some other error occurred and safest way to handle is to pass it back to user.
	if err != nil {
		return false, err
	}

	// Key is not trusted because it is not in the file.
	return false, nil
}

func (s *Service) addKnownHost(remote net.Addr, key ssh.PublicKey, knownFile string) (err error) {

	// Fallback to default known_hosts file
	if knownFile == "" {
		knownFile = defaultPath
	}

	parent := filepath.Dir(knownFile)
	if _, err = os.Stat(knownFile); os.IsNotExist(err) {
		err = os.MkdirAll(parent, 0755)
		if err != nil {
			return err
		}
	}

	f, err := os.OpenFile(knownFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)

	if err != nil {
		return err
	}

	defer f.Close()

	knownHost := knownhosts.Normalize(remote.String())

	_, err = f.WriteString(knownhosts.Line([]string{knownHost}, key) + "\n")

	return err
}
