package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

const (
	dbDriver   = "sqlite3"
	tempDBFile = "mwb.db"
)

func (c *Config) getMMData() (mmd MeetingData, err error) {
	jwpubBytes, err := c.getJWPub("mwb")
	if err != nil {
		return
	}

	tempDir, err := ioutil.TempDir("", "jwpub_fetcher_")
	if err != nil {
		return
	}
	defer os.RemoveAll(tempDir)

	contents, err := unzipFile(jwpubBytes, "contents")
	if err != nil {
		return
	}

	dbBytes, err := unzipFile(contents, "mwb*.db")
	if err != nil {
		return
	}

	// write this file to disk; quick check doesn't show an easy way for sqlite to handle things in memory only
	dbFilename := filepath.Join(tempDir, tempDBFile)
	dbFile, err := os.Create(dbFilename)
	if err != nil {
		return
	}
	_, err = dbFile.Write(dbBytes)
	if err != nil {
		return
	}
	dbFile.Close()

	sqlDB, err := sql.Open(dbDriver, dbFilename)
	if err != nil {
		return
	}
	defer sqlDB.Close()

	docs, err := getMWBDocuments(sqlDB)
	if err != nil {
		return
	}
	logrus.Debug("docs >>", docs)
	var docGroups []Document

	for _, doc := range docs {
		if c.Date != doc.Date {
			continue
		}
		docGroups = append(docGroups, doc)
	}

	if len(docGroups) == 0 {
		return mmd, errors.New("no docs found!")
	}

	mmd = MeetingData{
		DateString: docGroups[0].Date.Format("2006-01-02"),
		Songs:      getMWBSongs(sqlDB, docGroups),
	}

	if c.FetchOtherMedia {
		imageNames := getImageNames(sqlDB, docGroups)

		for _, name := range imageNames {
			// skip images we already have
			for _, p := range mmd.Pictures {
				if p.Name == name {
					continue
				}
			}

			// fetch image from contents
			pic, err := unzipFile(contents, name)
			if err != nil {
				return MeetingData{}, errors.New("problem getting pic")
			}

			// queue for storage
			mmd.Pictures = append(mmd.Pictures, file{Name: name, Payload: pic})
		}

		mmd.Videos = getMWBVideos(sqlDB, docGroups)

		linkedDocs := c.getLinkedDocs(sqlDB, docGroups)
		for _, ld := range linkedDocs {
			docMedia, err := c.getDocMedia(ld)
			if err != nil {
				logrus.Warn(err)
			}
			mmd.Pictures = append(mmd.Pictures, docMedia.Pictures...)
			mmd.Videos = append(mmd.Videos, docMedia.Videos...)
		}
	}

	return
}

func (c *Config) getWMData() (wmd MeetingData, err error) {
	jwpubBytes, err := c.getJWPub("w")
	if err != nil {
		return
	}

	tempDir, err := ioutil.TempDir("", "jwpub_fetcher_")
	if err != nil {
		return
	}
	defer os.RemoveAll(tempDir)

	contents, err := unzipFile(jwpubBytes, "contents")
	if err != nil {
		return
	}

	dbBytes, err := unzipFile(contents, "w*.db")
	if err != nil {
		return
	}

	// write this file to disk; quick check doesn't show an easy way for sqlite to handle things in memory only
	dbFilename := filepath.Join(tempDir, tempDBFile)
	dbFile, err := os.Create(dbFilename)
	if err != nil {
		return
	}
	_, err = dbFile.Write(dbBytes)
	if err != nil {
		return
	}
	dbFile.Close()

	sqlDB, err := sql.Open(dbDriver, dbFilename)
	if err != nil {
		return
	}
	defer sqlDB.Close()

	// pub := getPublication(sqlDB)
	docs, err := getWTDocuments(sqlDB)
	if err != nil {
		return
	}
	logrus.Debug("docs >>", docs)

	dates, err := getWTDates(sqlDB)
	if err != nil {
		return
	}
	logrus.Debug("dates >>", dates)

	for i, doc := range docs {
		if c.Date != dates[i] {
			continue
		}

		wmd = MeetingData{
			DateString: c.Date.Format("2006-01-02"),
			Songs:      getWTSongs(sqlDB, c.Date),
		}

		if c.FetchOtherMedia {
			pics := []file{}
			imageNames := getImageNames(sqlDB, []Document{{ID: doc}})
			for _, name := range imageNames {
				pic, err := unzipFile(contents, name)
				if err != nil {
					return wmd, errors.New("problem getting pic")
				}

				pics = append(pics, file{Name: name, Payload: pic})
			}
			wmd.Pictures = pics
		}

	}

	return
}

func (c *Config) getDocMedia(ld LinkedDocument) (md MeetingData, err error) {
	jwpubBytes, err := c.getJWPub(ld.PublicationSymbol)
	if err != nil {
		return
	}

	tempDir, err := ioutil.TempDir("", "jwpub_fetcher_")
	if err != nil {
		return
	}
	defer os.RemoveAll(tempDir)

	contents, err := unzipFile(jwpubBytes, "contents")
	if err != nil {
		return
	}

	dbBytes, err := unzipFile(contents, ld.PublicationSymbol+"*.db")
	if err != nil {
		return
	}

	// write this file to disk; quick check doesn't show an easy way for sqlite to handle things in memory only
	dbFilename := filepath.Join(tempDir, tempDBFile)
	dbFile, err := os.Create(dbFilename)
	if err != nil {
		return
	}
	_, err = dbFile.Write(dbBytes)
	if err != nil {
		return
	}
	dbFile.Close()

	sqlDB, err := sql.Open(dbDriver, dbFilename)
	if err != nil {
		return
	}
	defer sqlDB.Close()

	mepsDocs, err := getMEPSDocuments(sqlDB, ld.MepsDocumentID)
	if err != nil {
		return
	}
	logrus.Debug("docs >>", mepsDocs)

	for _, d := range mepsDocs {
		switch d.MimeType {
		case "image/jpeg":

			// skip known images that are not needed from 'th'
			if d.Name == "1102018440_univ_cnt_1.jpg" ||
				d.Name == "1102018440_univ_cnt_2.jpg" {
				continue
			}

			pic, err := unzipFile(contents, d.Name)
			if err != nil {
				return MeetingData{}, errors.New("problem getting pic")
			}
			md.Pictures = append(md.Pictures, file{Name: d.Name, Payload: pic})
		case "video/mp4":
			md.Videos = append(md.Videos, d.video)
		}
	}

	// logrus.Debug("getDocMedia()", md)
	return
}

func (c *Config) getJWPub(pub string) ([]byte, error) {
	date := c.Date
	switch pub {
	case "w":
		// W is published 2 months prior to it beeing needed for the meeting
		date = c.Date.AddDate(0, -2, 0)
	case "mwb":
		// mwb is released two months at a time.
		if c.Date.Month()%2 == 0 {
			date = c.Date.AddDate(0, -1, 0)
		}
	}

	m, err := c.getJWPubInfo(date.Year(), int(date.Month()), pub)
	if err != nil {
		return nil, err
	}

	jwpubItem := m.Files[c.Language].JWPUB[0]
	filename := filepath.Base(jwpubItem.File.URL)
	payload, err := c.getFromCache(filename, jwpubItem.File.Checksum)
	if err == nil {
		return payload, err
	}

	payload, err = c.download(jwpubItem)
	if err == nil {
		c.saveToCache(file{
			Name:    filename,
			Payload: payload,
		})
	}

	return payload, err
}

func (c *Config) getJWPubInfo(year, month int, pub string) (*mediaInfo, error) {
	var str string
	switch pub {
	case "w", "mwb":
		str = fmt.Sprintf("https://pubmedia.jw-api.org/GETPUBMEDIALINKS?issue=%d%02d&output=json&pub=%s&fileformat=JWPUB&alllangs=0&langwritten=%s&txtCMSLang=%s", year, month, pub, c.Language, c.Language)
	default:
		str = fmt.Sprintf("https://pubmedia.jw-api.org/GETPUBMEDIALINKS?output=json&pub=%s&fileformat=JWPUB&alllangs=0&langwritten=%s&txtCMSLang=%s", pub, c.Language, c.Language)
	}
	logrus.Debug("getJWPubInfo()", str)

	resp, err := c.HttpClient.Get(str)
	if err != nil {
		return nil, fmt.Errorf("failed to get media info for %s-%d-%02d workbook", c.Language, year, month)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("no workbook available for %v-%02d", year, month)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error reading info for workbook")
	}

	info := new(mediaInfo)
	err = json.Unmarshal(body, info)

	return info, err
}

func (c *Config) download(jwpi JWPubItem) ([]byte, error) {
	var body []byte
	resp, err := c.HttpClient.Get(jwpi.File.URL)
	if err != nil {
		return body, errors.New("failed to download " + jwpi.File.URL)
	}

	c.Progress.ProgressBar.SetValue(0)
	c.Progress.Total = 0
	c.Progress.Title = filepath.Base(jwpi.File.URL)
	c.Progress.ProgressBar.Max = float64(jwpi.Filesize)

	data := io.TeeReader(resp.Body, c.Progress)

	body, err = ioutil.ReadAll(data)
	if err != nil {
		return body, errors.New("error reading data from " + jwpi.File.URL)
	}

	if !validChecksum(jwpi.File.Checksum, body) {
		return nil, errors.New("invalid checksum for downloaded file")
	}

	return body, nil
}
