package sftp

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Upload file to sftp server
func (s *Service) Upload(localFile []byte, remoteFile string) (err error) {
	err = s.connect()
	if err != nil {
		return err
	}

	log.Printf("Uploading to [%s] ...\n", remoteFile)

	reader := bytes.NewReader(localFile)

	// Make remote directories recursion
	parent := filepath.Dir(remoteFile)
	path := string(filepath.Separator)
	dirs := strings.Split(parent, path)
	for _, dir := range dirs {
		path = filepath.Join(path, dir)
		s.client.Mkdir(path)
	}

	// Note: SFTP To Go doesn't support O_RDWR mode
	dstFile, err := s.client.OpenFile(remoteFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		log.Printf("Unable to open remote file: %v\n", err)
		return
	}
	defer dstFile.Close()

	written, err := io.Copy(dstFile, reader)
	if err != nil {
		log.Printf("Unable to upload local file: %v\n", err)
		os.Exit(1)
	}
	log.Printf("%d bytes copied\n", written)

	return
}
