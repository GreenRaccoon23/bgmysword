package main

import (
	"fmt"
	//"github.com/PuerkitoBio/goquery"
	//"log"
	"os"
	"strconv"
	"strings"
	//"unicode"
	//"github.com/mattn/go-sqlite3"
	//"github.com/fatih/color"
)

var (
	translation     string
	transName       string
	BibleGatewayUrl string
	logMe           bool
	//doc            *goquery.Document
	//db          *sql.DB
	//tx          *sql.Tx
)

func init() {
	defer ColorUnset()
	ImgINRI()
	numArgs := len(os.Args)
	if numArgs > 1 {
		firstArg := os.Args[1]
		switch firstArg {
		case "h", "-h", "help", "-help":
			printHelp()
		default:
			translation = os.Args[1]
		}
	}
	if numArgs > 2 {
		logMe = true
	}
	GenBibleGatewayUrl()
	CopyrightFetch()
}

func GenBibleGatewayUrl() {
	preUrl := "https://www.biblegateway.com/passage/?version="
	midUrl := "&search="
	BibleGatewayUrl = concat(preUrl, translation, midUrl)
}

func GenFullUrl(book, chap string) (url string) {
	url = concat(BibleGatewayUrl, book, "+", chap)
	return
}

func main() {
	defer ColorUnset()
	if logMe {
		logFile := GenLog()
		defer CloseLog(logFile)
	}
	ImgSword()
	GenModule()
	progressTranslation()
	defer db.Close()
	defer tx.Commit()
	bookLoop(Bible)
}

func bookLoop(Bible []BibleArchive) {
	var bRange int
	bRange = len(Bible)
	for i := 0; i < bRange; i++ {
		b := i + 1
		title := strings.Replace(Bible[i].Book, "+", " ", -1)
		titleSpaces := 30 - len(title)
		titleSpacing := strings.Repeat(" ", titleSpaces)
		var numSpacing string
		if b < 10 {
			numSpacing = " "
		} else {
			numSpacing = ""
		}
		progressBook(title, titleSpacing, numSpacing, b, bRange)
		chapterLoop(Bible[i])
		fmt.Println()
	}
}

func chapterLoop(data BibleArchive) {
	currentBook := strconv.Itoa(data.Index)
	var cRange int = data.ChapterRange
	for c := 1; c <= cRange; c++ {
		progressChapter(c, cRange)
		currentChapter := strconv.Itoa(c)
		url := GenFullUrl(data.Book, currentChapter)
		chapterText := Chapter.Parse(url)
		Log(chapterText)
		saveChapter(currentBook, currentChapter, chapterText)
	}
}

func saveChapter(currentBook string, currentChapter string, chapterText []string) {
	for i, s := range chapterText {
		verseNumber := i + 1
		verseText := s
		sqlInsBible(currentBook, currentChapter, verseNumber, verseText)
	}
}
