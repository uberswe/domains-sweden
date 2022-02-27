package sftp

import (
	"bytes"
	"fmt"
	"os"
)

func (s *Service) Fetch(remoteFile string) (localfile []byte, err error) {
	err = s.connect()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stdout, "Downloading [%s] ...\n", remoteFile)
	// Note: SFTP To Go doesn't support O_RDWR mode
	srcFile, err := s.client.OpenFile(remoteFile, os.O_RDONLY)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open remote file: %v\n", err)
		return
	}
	defer srcFile.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(srcFile)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
