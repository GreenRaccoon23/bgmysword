package main

import (
	"fmt"
	"os"
	"strconv"
	//"github.com/PuerkitoBio/goquery"
	//"github.com/mattn/go-sqlite3"
	//"github.com/fatih/color"
)

var (
	translation     string
	transName       string
	BibleGatewayUrl string
	logMe           bool
	//db          *sql.DB
	//tx          *sql.Tx
)

func init() {
	defer ColorUnset()
	ImgINRI()
	AnalyseArgs()
	GenBibleGatewayUrl()
	CopyrightFetch()
}

func AnalyseArgs() {
	numArgs := len(os.Args)
	if numArgs > 1 {
		switch os.Args[1] {
		case "h", "-h", "help", "-help":
			printHelp()
		default:
			translation = os.Args[1]
		}

	}
	if numArgs > 2 {
		logMe = true
		GenLog()
	}
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
	defer CloseLog()
	defer ColorUnset()
	ImgSword()
	GenModule()
	defer CloseModule()
	progressTranslation()
	BookLoop(Bible)
}

func BookLoop(Bible []BibleArchive) {
	for i := 0; i < 66; i++ {
		progressBook(i)
		ChapterLoop(Bible[i])
		fmt.Println()
	}
}

func ChapterLoop(data BibleArchive) {
	currentBook := strconv.Itoa(data.Index)
	cRange := data.ChapterRange
	for c := 1; c <= cRange; c++ {
		progressChapter(c, cRange)

		currentChapter := strconv.Itoa(c)
		url := GenFullUrl(data.Book, currentChapter)

		chapterText := Chapter.Parse(url)

		SaveChapter(currentBook, currentChapter, chapterText)
	}
}

func SaveChapter(currentBook string, currentChapter string, chapterText []string) {
	for i, s := range chapterText {
		verseNumber := i + 1
		verseText := s

		sqlInsBible(currentBook, currentChapter, verseNumber, verseText)
	}
}
