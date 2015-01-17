package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
	//"github.com/mattn/go-sqlite3"
)

var (
	translation string
	transName   string
	base_url    []string
	doc         *goquery.Document
	footnoteMap = make(map[string]string)
	//db          *sql.DB
	//tx          *sql.Tx
)

func init() {
	defer color.Unset()
	if len(os.Args) > 1 {
		translation = os.Args[1]
	} else {
		translation = inputTranslation()
	}
	base_url = genUrl()
	copyrightFetch()
}

func genUrl() []string {
	pre_url := "https://www.biblegateway.com/passage/?version="
	mid_url := "&search="
	base_url := []string{pre_url, translation, mid_url}
	return base_url
}

func concatUrl(book, chap string) string {
	url_slice := append(base_url, book, "+", chap)
	return str(url_slice)
}

func DocIt(url string) *goquery.Document {
	theDocIsIn, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	return theDocIsIn
}

func copyrightFetch() {
	var copyrightInfo, publisherInfo string
	url := concatUrl("Genesis", "1")
	copyrightDoc := DocIt(url)
	copyrightDoc.Find(".publisher-info-bottom").Each(func(i int, s *goquery.Selection) {
		copyrightInfo = s.Find("p").Text()
		publisherInfo = s.Find("p a").Text()
		transName = s.Find("strong").Text()
	})
	copyrightAccept(copyrightInfo, publisherInfo)
}

func main() {
	defer color.Unset()
	imgINRI()
	genBible()
	progressTranslation()
	defer tx.Commit()
	defer db.Close()
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
	}
}

func chapterLoop(data BibleArchive) {
	currentBook := strconv.Itoa(data.Index)
	var cRange int = data.ChapterRange
	for c := 1; c <= cRange; c++ {
		progressChapter(c, cRange)
		currentChapter := strconv.Itoa(c)
		url := concatUrl(data.Book, currentChapter)
		chapterText := parseChapter(url)
		for i := 0; i < len(chapterText); i++ {
			verseNumber := i + 1
			verseText := chapterText[i]
			sqlInsBible(currentBook, currentChapter, verseNumber, verseText)
		}
	}
	fmt.Printf("\n")
}

func parseChapter(url string) []string {
	doc = DocIt(url)
	cleanDoc()
	queryClass, vRange := analyseDoc()
	prepDoc(queryClass)
	chapterText := verseLoop(queryClass, vRange)
	return chapterText
}

func cleanDoc() {
	rejects := []string{".chapternum", ".versenum", ".crossreference"}
	for _, reject := range rejects {
		doc.Find(reject).Remove()
	}
}

func analyseDoc() (string, int) {
	// Make a map of the footnotes together with their expanded text.
	doc.Find("ol").Contents().Each(func(i int, s *goquery.Selection) {
		id, _ := s.Attr("id")
		if strings.Contains(id, "fen") {
			fnLet := strLast(id)
			fnText := s.Find(".footnote-text").Text()
			footnoteMap[fnLet] = fnText
		}
	})
	// Note the number of verses and the html class which denotes passage text.
	lastClass, _ := doc.Find("p .text").Last().Attr("class")
	baseClass := prefixHyphenSecond(lastClass)
	queryClass := suffixSpace(baseClass)
	lastVerse := suffixHyphenSecond(lastClass)
	vRange, _ := strconv.Atoi(lastVerse)
	return queryClass, vRange
}

func prepDoc(queryClass string) {
	// Mark titles.
	doc.Find("h3 .text").Each(func(i int, s *goquery.Selection) {
		s.SetAttr("id", "title")
	})
	// Mark the opening line of each paragraph.
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		class, _ := s.Contents().Attr("class")
		if strings.Contains(class, queryClass) {
			s.Contents().First().SetAttr("id", "paragraph")
		}
	})
}

func verseLoop(queryClass string, vRange int) []string {
	var chapterText []string
	for i := 1; i <= vRange; i++ {
		v := strconv.Itoa(i)
		currClass := str([]string{queryClass, v})
		chapterText = append(chapterText, parseVerse(currClass))
	}
	return chapterText
}

func parseVerse(vClass string) string {
	var vTemp []string
	var poetry bool
	digitPrefix := unicode.IsDigit([]rune(vClass)[0])
	switch digitPrefix {
	case false:
		vClass = combine(".", vClass)
		doc.Find(vClass).Each(func(i int, s *goquery.Selection) {
			vTemp, poetry = checkPoetry(s, vTemp, poetry)
		})
	case true:
		doc.Find(".text").Each(func(i int, s *goquery.Selection) {
			class, _ := s.Attr("class")
			if strings.Contains(class, vClass) {
				vTemp, poetry = checkPoetry(s, vTemp, poetry)
			}
		})
	}
	if poetry {
		vTemp = cleanIndent(vTemp)
	}
	vFull := str(vTemp)
	return vFull
}

func checkPoetry(s *goquery.Selection, vTemp []string, poetry bool) ([]string, bool) {
	var classParents string
	classParents, _ = s.Parents().Attr("class")
	if strings.Contains(classParents, "chapter") {
		classParents, _ = s.Parents().Parents().Attr("class")
	}
	switch classParents {
	case "line":
		vTemp = parsePoetryLine(s, vTemp)
		poetry = true
	case "indent-1":
		vTemp = parsePoetryIndent(s, vTemp)
	default:
		vTemp = parseProse(s, vTemp)
	}
	return vTemp, poetry
}

func parseProse(s *goquery.Selection, vTemp []string) []string {
	id, _ := s.Attr("id")
	switch id {
	case "paragraph":
		pTemp := fmtParagraph(vTemp)
		vTemp = append(pTemp, fmtDefault(s.Contents()))
	case "title":
		vTemp = fmtTitle(s, vTemp)
	default:
		vTemp = append(vTemp, fmtDefault(s.Contents()))
	}
	return vTemp
}

func parsePoetryLine(s *goquery.Selection, vTemp []string) []string {
	id, _ := s.Attr("id")
	switch id {
	case "paragraph":
		pTemp := fmtParagraph(vTemp)
		vTemp = append(pTemp, "<PI1>", fmtDefault(s.Contents()))
	case "title":
		vTemp = fmtTitle(s, vTemp)
	default:
		vTemp = append(vTemp, "<PI1>", fmtDefault(s.Contents()))
	}
	vTemp = append(vTemp, "<CI>")
	return vTemp
}

func parsePoetryIndent(s *goquery.Selection, vTemp []string) []string {
	id, _ := s.Contents().Attr("id")
	switch id {
	case "paragraph":
		pTemp := fmtParagraph(vTemp)
		vTemp = append(pTemp, "<PI3>", fmtDefault(s.Contents()))
	case "title":
		vTemp = fmtTitle(s, vTemp)
	default:
		vTemp = append(vTemp, "<PI3>", fmtDefault(s.Contents()))
	}
	vTemp = append(vTemp, "<CI>")
	return vTemp
}

func fmtTitle(s *goquery.Selection, vTemp []string) []string {
	return append(vTemp, "<TS>", s.Text(), "<Ts>")
}

func fmtParagraph(vTemp []string) []string {
	var pTemp []string
	var formatted bool
	for _, s := range vTemp {
		if strings.Contains(s, "<TS>") {
			new := strings.Replace(s, "<TS>", "<CM><TS>", -1)
			pTemp = append(pTemp, new)
			formatted = true
		} else {
			pTemp = append(pTemp, s)
		}
	}
	if formatted == false {
		pTemp = append(vTemp, "<CM>")
	}
	return pTemp
}

func fmtTitleParagraph(vTemp []string, ts int) []string {
	var pTemp []string
	pTemp = append(vTemp[:ts], "<CM>")
	for _, s := range vTemp[ts:] {
		pTemp = append(pTemp, s)
	}
	return pTemp
}

func fmtDefault(sel *goquery.Selection) string {
	var sTemp []string
	for i := range sel.Nodes {
		s := sel.Eq(i)
		class, _ := s.Attr("class")
		woj := strings.Contains(class, "woj")
		footnote := strings.Contains(class, "footnote")
		switch {
		case woj:
			sTemp = append(sTemp, fmtWoJ(s.Contents()))
		case footnote:
			sTemp = append(sTemp, fmtFootnote(s))
		default:
			sTemp = append(sTemp, s.Text())
		}
	}
	sFull := str(sTemp)
	return sFull
}

func fmtWoJ(sel *goquery.Selection) string {
	var wojTemp []string
	wojTemp = append(wojTemp, "<FR>")
	for i := range sel.Nodes {
		s := sel.Eq(i)
		class, _ := s.Attr("class")
		footnote := strings.Contains(class, "footnote")
		if footnote {
			wojTemp = append(wojTemp, fmtFootnote(s))
		} else {
			wojTemp = append(wojTemp, s.Text())
		}
	}
	wojTemp = append(wojTemp, "<Fr>")
	wojFull := str(wojTemp)
	return wojFull
}

func fmtFootnote(s *goquery.Selection) string {
	var fTemp []string
	fnLetter := s.Find("a").Text()
	fnText := footnoteMap[fnLetter]
	fTemp = append(fTemp, "<RF>", fnText, "<Rf>")
	fFull := str(fTemp)
	return fFull
}

func cleanIndent(vClass []string) []string {
	var iTemp []string
	var foundFirst bool
	foundFirst = true
	for _, s := range vClass {
		if strings.Contains(s, "<PI1>") {
			switch foundFirst {
			case true:
				foundFirst = false
				iTemp = append(iTemp, s)
			default:
				new := strings.Replace(s, "<PI1>", "<PI2>", -1)
				iTemp = append(iTemp, new)
			}
		} else {
			iTemp = append(iTemp, s)
		}
	}
	return iTemp
}
