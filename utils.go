package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

// scrubTitle replacing special characters
func scrubTitle(title string) string {
	title = strings.Replace(title, ":", "_", -1)
	title = strings.Replace(title, "'", "_", -1)
	title = strings.Replace(title, "/", "_", -1)

	return title
}

// checkIfInSlice check if a value is in the slice
func checkIfInSlice(value string, src []string) bool {
	for _, v := range src {
		if value == v {
			return true
		}
	}

	return false
}

// convertEscaped replaces escaped chars back to original values
func convertEscaped(data *[]byte) {
	*data = bytes.Replace(*data, []byte("&#39;"), []byte("'"), -1)
	*data = bytes.Replace(*data, []byte("&#34;"), []byte(`"`), -1)
	*data = bytes.Replace(*data, []byte("&#xA;"), []byte("\n"), -1)
}

// copyFile does what is says, copies files
func copyFile(src, dest string) error {
	inFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer inFile.Close()

	inputData, err := io.ReadAll(inFile)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dest, inputData, 0644); err != nil {
		log.Errorf("There was an error copying data: %v", err)
		return err
	}
	return nil
}
