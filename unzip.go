package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"io/ioutil"
	"path/filepath"
)

// unzipFile will decompress a zip archive, and return the bytes of the first file matching the pattern
// pattern is matched using filepath.Match
func unzipFile(zipBytes []byte, pattern string) ([]byte, error) {
	reader := bytes.NewReader(zipBytes)
	r, err := zip.NewReader(reader, int64(len(zipBytes)))
	if err != nil {
		return nil, err
	}

	for _, f := range r.File {
		if match, _ := filepath.Match(pattern, f.Name); !match {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, err
		}

		b, err := ioutil.ReadAll(rc)
		rc.Close()
		return b, err
	}

	return nil, errors.New("no files matched the pattern")
}
