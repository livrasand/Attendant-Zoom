package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
)

func NewConfig() *Config {
	c := Config{}
	c.readConfigFromFile()
	c.HttpClient = newHttpClient()

	if err := createDirIfNotExist(c.SaveLocation); err != nil {
		logrus.Warn(err)
	}

	return &c
}

func newHttpClient() *retryablehttp.Client {
	client := retryablehttp.NewClient()
	client.RetryWaitMin = 5 * time.Second
	client.RetryWaitMax = 60 * time.Second
	client.RetryMax = 20
	client.Logger = nil
	return client
}

func (c *Config) LoadDefaults() {
	homeDir := os.Getenv("HOME")

	c.AutoFetchMeetingData = true
	c.FetchOtherMedia = true
	c.CreatePlaylist = true
	c.SaveLocation = filepath.Join(homeDir, "Downloads/meetings")
	c.Resolution = RES720
	c.Language = "S"
	c.PubSymbols = []string{"th", "lff"}
	c.CacheLocation = filepath.Join(homeDir, "Downloads/meetings_cache")
}

func (c *Config) readConfigFromFile() {
	homeDir := os.Getenv("HOME")
	c.LoadDefaults()

	data, err := os.ReadFile(filepath.Join(homeDir, CONFIG_FILE))
	if err != nil {
		c.writeConfigToFile()
	}

	err = toml.Unmarshal(data, c)
	if err != nil {
		logrus.Warn(err)
	}
}

func (c *Config) writeConfigToFile() {
	homeDir := os.Getenv("HOME")

	logrus.Info("Guardando configuraciones")

	config := struct {
		AutoFetchMeetingData bool
		FetchOtherMedia      bool
		CreatePlaylist       bool
		PurgeSaveDir         bool
		Resolution           string
		SaveLocation         string
		Language             string
		CacheLocation        string
		PubSymbols           []string
	}{
		AutoFetchMeetingData: c.AutoFetchMeetingData,
		FetchOtherMedia:      c.FetchOtherMedia,
		CreatePlaylist:       c.CreatePlaylist,
		PurgeSaveDir:         c.PurgeSaveDir,
		Resolution:           c.Resolution,
		SaveLocation:         c.SaveLocation,
		Language:             c.Language,
		PubSymbols:           c.PubSymbols,
		CacheLocation:        c.CacheLocation,
	}

	configToml, err := toml.Marshal(config)
	if err != nil {
		logrus.Warn(err)
	}

	err = os.WriteFile(filepath.Join(homeDir, CONFIG_FILE), configToml, 0644)
	if err != nil {
		logrus.Warn(err)
	}
}
