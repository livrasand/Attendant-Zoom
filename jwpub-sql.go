package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

const jwpubDateFormat = "20060102"

func getMEPSDocuments(db *sql.DB, mdocid string) (mepsDocuments []mepsDocument, err error) {
	// get all docIDs
	sqlQuery := fmt.Sprintf(`SELECT Multimedia.MimeType,
																	Multimedia.FilePath,
																	Multimedia.Track,
																	Multimedia.KeySymbol,
																	Multimedia.MepsDocumentId,
																	Multimedia.IssueTagNumber
													 FROM Document
													 INNER JOIN DocumentMultimedia
													 ON Document.DocumentId = DocumentMultimedia.DocumentId
													 INNER JOIN Multimedia
													 ON DocumentMultimedia.MultimediaId = Multimedia.MultimediaId
													 WHERE Document.MepsDocumentId = %s`, mdocid)

	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("unable to get allDocs: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mepsDoc mepsDocument
		err = rows.Scan(
			&mepsDoc.MimeType,
			&mepsDoc.Name,
			&mepsDoc.Track,
			&mepsDoc.KeySymbol,
			&mepsDoc.MepsDocumentID,
			&mepsDoc.IssueTagNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to scan allDocs row: %v", err)
		}
		mepsDocuments = append(mepsDocuments, mepsDoc)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("row error after allDocs query: %v", err)
	}

	return
}

func getMWBDocuments(db *sql.DB) (mwbDocuments []Document, err error) {
	// get all docIDs
	allDocs := make([]Document, 0)
	sqlQuery := `SELECT DocumentId AS LastDocId
	             FROM Document
	             ORDER BY DocumentId ASC`

	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("unable to get allDocs: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mepsDoc Document
		err = rows.Scan(
			&mepsDoc.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to scan allDocs row: %v", err)
		}
		allDocs = append(allDocs, mepsDoc)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("row error after allDocs query: %v", err)
	}

	// Get primary docIDs w/ dates
	primaryDocs := make([]Document, 0)
	sqlQuery = `SELECT DocumentId, CAST(FirstDateOffset AS TEXT)
	            FROM DatedText`

	rows, err = db.Query(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("unable to get documents: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mepsDoc Document
		var date string
		err = rows.Scan(
			&mepsDoc.ID,
			&date,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to scan document row: %v", err)
		}

		mepsDoc.Date, err = time.Parse(jwpubDateFormat, date)
		if err != nil {
			return
		}
		primaryDocs = append(primaryDocs, mepsDoc)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("row error after document query: %v", err)
	}

	// glue it all together
	var docDate time.Time
	for _, doc := range allDocs {
		for _, pdoc := range primaryDocs {
			if pdoc.ID == doc.ID {
				docDate = pdoc.Date
			}
		}
		doc.Date = docDate
		mwbDocuments = append(mwbDocuments, doc)
	}

	return
}

func getWTDocuments(db *sql.DB) (wtDocumentIDs []int, err error) {
	wtDocumentIDs = make([]int, 0)
	sqlQuery := `SELECT DocumentId
							 FROM Document
							 WHERE Class=40
							 ORDER BY DocumentId ASC`

	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("unable to get documents: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var docID int
		err = rows.Scan(
			&docID,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to scan document row: %v", err)
		}
		wtDocumentIDs = append(wtDocumentIDs, docID)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("row error after document query: %v", err)
	}

	return
}

func getWTDates(db *sql.DB) (wtDates []time.Time, err error) {
	wtDates = make([]time.Time, 0)
	sqlQuery := `SELECT  CAST(FirstDateOffset AS TEXT) AS Date
							 FROM DatedText
							 ORDER BY DatedTextId ASC`

	rows, err := db.Query(sqlQuery)
	if err != nil {
		log.Fatal("unable to get dates: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var date string
		err = rows.Scan(
			&date,
		)
		if err != nil {
			log.Fatal("unable to scan date row: ", err)
		}

		wdate, err := time.Parse(jwpubDateFormat, date)
		if err != nil {
			return nil, err
		}
		wtDates = append(wtDates, wdate)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal("row error after date query: ", err)
	}

	return
}

func getMWBSongs(db *sql.DB, docIDs []Document) (songs []string) {
	var whereDID string
	for i, did := range docIDs {
		if i != 0 {
			whereDID += " OR "
		}
		whereDID += fmt.Sprintf("DocumentId=%v", did.ID)
	}

	sqlQuery := fmt.Sprintf(`SELECT Track
													 FROM DocumentMultimedia
													 INNER JOIN Multimedia
													 ON DocumentMultimedia.MultimediaId = Multimedia.MultimediaId
													 WHERE (%s)
													 AND BeginParagraphOrdinal IS NOT NULL
													 AND KeySymbol='sjjm'
													 ORDER BY DocumentMultimediaId ASC;`, whereDID)

	rows, err := db.Query(sqlQuery)
	if err != nil {
		log.Fatal("unable to get songs: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mm Multimedia
		err = rows.Scan(
			&mm.Track,
		)
		if err != nil {
			log.Fatal("unable to scan song row: ", err)
		}
		songs = append(songs, mm.Track)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal("row error after song query: ", err)
	}

	logrus.Debug("getMWBSongs()", songs)
	return
}

func getMWBVideos(db *sql.DB, docIDs []Document) (videos []video) {
	var whereDID string
	for i, did := range docIDs {
		if i != 0 {
			whereDID += " OR "
		}
		whereDID += fmt.Sprintf("DocumentId=%v", did.ID)
	}

	sqlQuery := fmt.Sprintf(`SELECT Track, KeySymbol, MepsDocumentId, IssueTagNumber
													 FROM DocumentMultimedia
													 INNER JOIN Multimedia
													 ON DocumentMultimedia.MultimediaId = Multimedia.MultimediaId
													 WHERE (%s)
													 AND Multimedia.MimeType="video/mp4"
 													 AND ( MepsDocumentId IS NOT NULL OR IssueTagNumber != 0)
													 ORDER BY DocumentMultimediaId ASC;`, whereDID)

	rows, err := db.Query(sqlQuery)
	if err != nil {
		log.Fatal("unable to get videos: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var v video
		err = rows.Scan(
			&v.Track,
			&v.KeySymbol,
			&v.MepsDocumentID,
			&v.IssueTagNumber,
		)
		if err != nil {
			log.Fatal("unable to scan video row: ", err)
		}
		videos = append(videos, v)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal("row error after song query: ", err)
	}

	logrus.Debug("getMWBSongs()", videos)
	return
}

func getWTSongs(db *sql.DB, date time.Time) (songs []string) {
	d := date.Format(jwpubDateFormat)
	sqlQuery := fmt.Sprintf(`SELECT Multimedia.Track
							 FROM DatedText
							 INNER JOIN Multimedia
							 ON DatedText.BeginParagraphOrdinal = Multimedia.MultimediaId
							 WHERE FirstDateOffset=%s
							 AND KeySymbol='sjjm'
							 UNION ALL
							 SELECT Multimedia.Track
							 FROM DatedText
							 INNER JOIN Multimedia
							 ON DatedText.EndParagraphOrdinal = Multimedia.MultimediaId
							 WHERE FirstDateOffset=%s
							 AND KeySymbol='sjjm';
`, d, d)

	rows, err := db.Query(sqlQuery)
	if err != nil {
		log.Fatal("unable to get wtsongs: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mm Multimedia
		err = rows.Scan(
			&mm.Track,
		)
		if err != nil {
			log.Fatal("unable to scan wtsong row: ", err)
		}
		songs = append(songs, mm.Track)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal("row error after wtsong query: ", err)
	}

	logrus.Debug("getWTSongs()", songs)
	return
}

func getImageNames(db *sql.DB, docIDs []Document) (files []string) {
	var whereDID string
	for i, did := range docIDs {
		if i != 0 {
			whereDID += " OR "
		}
		whereDID += fmt.Sprintf("DocumentId=%v", did.ID)
	}

	sqlQuery := fmt.Sprintf(`SELECT FilePath
													 FROM DocumentMultimedia
													 INNER JOIN Multimedia
													 ON DocumentMultimedia.MultimediaId = Multimedia.MultimediaId
													 WHERE (%s)
													 AND BeginParagraphOrdinal IS NOT NULL
													 AND FilePath!=''
													 ORDER BY DocumentMultimediaId ASC;`, whereDID)

	rows, err := db.Query(sqlQuery)
	if err != nil {
		log.Fatal("unable to get documents: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mm Multimedia
		err = rows.Scan(
			&mm.FilePath,
		)
		if err != nil {
			log.Fatal("unable to scan document row: ", err)
		}
		files = append(files, mm.FilePath)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal("row error after document query: ", err)
	}

	logrus.Debug("getImageNames()", files)
	return
}

func (c *Config) getLinkedDocs(db *sql.DB, docIDs []Document) (docs []LinkedDocument) {
	var whereDID string
	for i, did := range docIDs {
		if i != 0 {
			whereDID += " OR "
		}
		whereDID += fmt.Sprintf("DocumentExtract.DocumentId=%v", did.ID)
	}

	var refSymbol string
	for i, pubSymbol := range c.PubSymbols {
		if i != 0 {
			refSymbol += " OR "
		}
		refSymbol += fmt.Sprintf("RefPublication.UndatedSymbol == '%s'", pubSymbol)
	}

	sqlQuery := fmt.Sprintf(`SELECT RefPublication.UndatedSymbol, Extract.RefMepsDocumentId
													 FROM DocumentExtract
													 INNER JOIN Extract
													 ON DocumentExtract.ExtractId = Extract.ExtractId
													 INNER JOIN  RefPublication
													 ON Extract.RefPublicationId = RefPublication.RefPublicationId
													 WHERE (%s)
													 AND (%s);`,
		whereDID, refSymbol)

	rows, err := db.Query(sqlQuery)
	if err != nil {
		log.Fatal("unable to get lined documents: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ld LinkedDocument
		err = rows.Scan(
			&ld.PublicationSymbol,
			&ld.MepsDocumentID,
		)
		if err != nil {
			log.Fatal("unable to scan lined document row: ", err)
		}
		docs = append(docs, ld)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal("row error after lined document query: ", err)
	}

	logrus.Debug("getLinkedDocs()", docs)
	return
}
