package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

func WeekOf(date time.Time) time.Time {
	return RelativeDay(date, time.Monday)
}

func RelativeDay(date time.Time, toDOW time.Weekday) (out time.Time) {
	out = date
	for out.Weekday() != time.Monday {
		out = out.AddDate(0, 0, -1)
	}
	for out.Weekday() != toDOW {
		out = out.AddDate(0, 0, 1)
	}
	return
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.Mkdir(dir, fs.FileMode(0777)); err != nil {
			return err
		}
	}
	return nil
}

func validChecksum(checksum string, payload []byte) bool {
	if checksum != fmt.Sprintf("%x", md5.Sum(payload)) {
		return false
	}
	return true
}

func (c *Config) getFromCache(filename, checksum string) ([]byte, error) {
	payload, err := os.ReadFile(filepath.Join(c.CacheLocation, filename))
	if err != nil {
		return nil, err
	}

	if !validChecksum(checksum, payload) {
		return nil, errors.New("suma de comprobación no válida en el archivo en caché")
	}

	logrus.Infof("usando caché para %s", filename)
	return payload, err
}

func (c *Config) saveToCache(f file) error {
	createDirIfNotExist(c.CacheLocation)
	err := os.WriteFile(filepath.Join(c.CacheLocation, f.Name), f.Payload, 0644)
	if err != nil {
		return err
	}
	logrus.Infof("almacenamiento en caché %s", f.Name)
	return nil
}

func (c *Config) saveAndLink(f file) error {
	err := c.saveToCache(f)
	if err != nil {
		return err
	}

	return os.Symlink(filepath.Join(c.CacheLocation, f.Name), filepath.Join(c.SaveLocation, f.Name))
}
