package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
	//"github.com/mattn/go-sqlite3"
	//"github.com/fatih/color"
)

var (
	translation    string
	transName      string
	baseUrl        string
	chapterText    []string
	doc            *goquery.Document
	footnoteMap    = make(map[string]string)
	containsPoetry bool
	//db          *sql.DB
	//tx          *sql.Tx
)

func init() {
	defer ColorUnset()
	ImgINRI()
	if len(os.Args) > 1 {
		translation = os.Args[1]
	} else {
		translation = inputTranslation()
	}
	genBaseUrl()
	copyrightFetch()
}

func genBaseUrl() {
	preUrl := "https://www.biblegateway.com/passage/?version="
	midUrl := "&search="
	baseUrl = combine(preUrl, translation, midUrl)
	return
}

func genFullUrl(book, chap string) (url string) {
	url = combine(baseUrl, book, "+", chap)
	return
}

func genDoc(url string) (doc *goquery.Document) {
	var err error
	var tried bool
	doc, err = goquery.NewDocument(url)
	if err != nil {
		if tried {
			log.Fatal(err)
		} else {
			tried = true
			doc = genDoc(url)
			return
		}
	}
	return
}

func copyrightFetch() {
	var copyrightInfo, publisherInfo string
	url := genFullUrl("Genesis", "1")
	copyrightDoc := genDoc(url)
	copyrightInfo = copyrightDoc.Find(".publisher-info-bottom p").Text()
	publisherInfo = copyrightDoc.Find(".publisher-info-bottom p a").Text()
	transName = copyrightDoc.Find(".publisher-info-bottom strong").Text()
	copyrightAccept(copyrightInfo, publisherInfo)
}

func main() {
	defer ColorUnset()
	ImgSword()
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
		url := genFullUrl(data.Book, currentChapter)
		parseChapter(url)
		for i := 0; i < len(chapterText); i++ {
			verseNumber := i + 1
			verseText := chapterText[i]
			sqlInsBible(currentBook, currentChapter, verseNumber, verseText)
		}
		chapterText = nil
	}
	fmt.Printf("\n")
}

func parseChapter(url string) {
	// Fetch webpage for current chapter.
	doc = genDoc(url)
	// Filter out unneeded text.
	cleanDoc()
	// Find footnotes, number of verses, and the html class for the current chapter.
	chapterClass, vRange := analyseDoc()
	// Mark titles and paragraphs.
	prepDoc(chapterClass)
	// Parse the webpage.
	verseLoop(chapterClass, vRange)
}

func cleanDoc() {
	// Remove chapter numbers, verse numbers, and cross references
	//    (because MySword formats them differently/automatically).
	rejects := []string{".chapternum", ".versenum", ".crossreference"}
	for _, reject := range rejects {
		doc.Find(reject).Remove()
	}
}

func analyseDoc() (chapterClass string, vRange int) {
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
	doubleClass := prefixHyphenSecond(lastClass)
	chapterClass = suffixSpace(doubleClass)
	lastVerse := suffixHyphenSecond(lastClass)
	vRange, _ = strconv.Atoi(lastVerse)
	return
}

func prepDoc(chapterClass string) {
	// Mark titles.
	doc.Find("h3 .text").Each(func(i int, s *goquery.Selection) {
		s.SetAttr("id", "title")
	})
	doc.Find("h4 .text").Each(func(i int, s *goquery.Selection) {
		s.SetAttr("id", "title")
	})
	// Mark the opening line of each paragraph.
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		class, _ := s.Contents().Attr("class")
		if strings.Contains(class, chapterClass) {
			s.Contents().First().SetAttr("id", "paragraph")
		}
	})
}

func verseLoop(chapterClass string, vRange int) {
	for i := 1; i <= vRange; i++ {
		// For each verse:
		v := strconv.Itoa(i)
		currClass := str([]string{chapterClass, v})
		// Each element of chapterText will hold the text
		//    for each verse of the current chapter.
		chapterText = append(chapterText, parseVerse(currClass))
	}
}

func parseVerse(vClass string) (fullVerseText string) {
	var targetClass string
	var vTemp []string
	containsPoetry = false
	// For books that begin with a number, do a workaround.
	//    I.e., rather than searching the chapter for the current verse,
	//          iterate through ALL of the chapter's text elements,
	//          and for each one, test whether it belongs in the current verse.
	digitPrefix := unicode.IsDigit([]rune(vClass)[0])
	switch digitPrefix {
	case true:
		// For books that begin with a number:
		targetClass = combine("text ", vClass)
		doc.Find(".text").Each(func(i int, s *goquery.Selection) {
			class, _ := s.Attr("class")
			switch class {
			case targetClass:
				// If text belongs in the current verse:
				vTemp = analyseTag(s, vTemp)
			}
		})
	default:
		// For all other books:
		targetClass = combine(".", vClass)
		// Find the current verse.
		doc.Find(targetClass).Each(func(i int, s *goquery.Selection) {
			vTemp = analyseTag(s, vTemp)
		})
	}
	// Clean up indentation for poetry verses.
	if containsPoetry {
		vTemp = cleanIndentTags(vTemp)
	}
	fullVerseText = str(vTemp)
	return
}

func analyseTag(s *goquery.Selection, vTemp []string) (vFmtd []string) {
	var classParents string
	classParents, _ = s.Parents().Attr("class")
	if strings.Contains(classParents, "chapter") {
		classParents, _ = s.Parents().Parents().Attr("class")
	}
	switch classParents {
	case "line":
		// If first-line poetry:
		containsPoetry = true
		vFmtd = parsePoetryLine(s, vTemp)
		return
	case "indent-1":
		// If second-line poetry:
		vFmtd = parsePoetryIndent(s, vTemp)
		return
	default:
		// If normal text:
		vFmtd = parseProse(s, vTemp)
		return
	}
}

func parseProse(s *goquery.Selection, vTemp []string) (vFmtd []string) {
	id, _ := s.Attr("id")
	switch id {
	case "title":
		// If this html id tag was labeled "title" earlier:
		vFmtd = fmtTitle(s, vTemp)
		return
	case "paragraph":
		// If this html id tag was labeled "paragraph" earlier:
		vTemp = tagParagraph(vTemp)
		fallthrough
	default:
		// If normal content:
		vFmtd = parseText(s.Contents(), vTemp)
		return
	}
}

func parsePoetryLine(s *goquery.Selection, vTemp []string) (vFmtd []string) {
	id, _ := s.Attr("id")
	switch id {
	case "title":
		// If this html id tag was labeled "title" earlier:
		vFmtd = fmtTitle(s, vTemp)
		return
	case "paragraph":
		// If this html id tag was labeled "paragraph" earlier:
		vTemp = tagParagraph(vTemp)
		fallthrough
	default:
		// If normal content:
		vTemp = append(vTemp, "<PI1>")         // Tag to single-indent
		vTemp = parseText(s.Contents(), vTemp) //     Indented text
		vFmtd = append(vTemp, "<CI>")          // Tag to close indent
		return
	}
}

func parsePoetryIndent(s *goquery.Selection, vTemp []string) (vFmtd []string) {
	id, _ := s.Contents().Attr("id")
	switch id {
	case "title":
		// If this html id tag was labeled "title" earlier:
		vFmtd = fmtTitle(s, vTemp)
	case "paragraph":
		// If this html id tag was labeled "paragraph" earlier:
		vTemp = tagParagraph(vTemp)
		fallthrough
	default:
		// If normal content:
		vTemp = append(vTemp, "<PI3>")         // Tag to triple-indent
		vTemp = parseText(s.Contents(), vTemp) //     Indented text
		vFmtd = append(vTemp, "<CI>")          // Tag to close indent

	}
	return
}

func tagParagraph(vTemp []string) (vFmtd []string) {
	vLength := len(vTemp)
	switch vLength {
	case 0:
		// If the new paragraph starts at the begininning of the verse:
		tagParagraphPreviousVerse()
		vFmtd = vTemp
		return
	default:
		// If the new paragraph starts in the middle of the verse:
		vFmtd = tagParagraphCurrentVerse(vTemp)
		return
	}
}

func tagParagraphPreviousVerse() {
	var cLength int
	cLength = len(chapterText)
	if cLength > 0 {
		// If the new paragraph does NOT start at the beginning of the chapter,
		//    put the paragraph tag at the very end of the last verse.
		cLast := cLength - 1
		prevVerse := chapterText[cLast]
		prevVerseNew := combine(prevVerse, "<CM>")
		cTemp := append(chapterText[:cLast], prevVerseNew)
		chapterText = nil
		chapterText = cTemp
	}
}

func tagParagraphCurrentVerse(vTemp []string) (vFmtd []string) {
	var vLast int
	vLength := len(vTemp)
	vLast = vLength - 1
	switch vTemp[vLast] {
	case "<Ts>":
		// If the new paragraph DOES begin after a title,
		//    the paragraph tag is NOT unneeded.
		vFmtd = vTemp
		return
	default:
		// If the new paragraph does NOT begin after a title,
		//    a paragraph tag IS needed.
		vFmtd = append(vTemp[:vLast], "<CM>", vTemp[vLast])
		return
	}
}

func fmtTitle(s *goquery.Selection, vTemp []string) (vFmtd []string) {
	vTemp = append(vTemp, "<TS>") // Opening title tag
	vTemp = parseText(s, vTemp)   //    Title text
	vFmtd = append(vTemp, "<Ts>") // Closing title tag
	return
}

func parseText(sel *goquery.Selection, vTemp []string) (vFmtd []string) {
	var textTemp []string
	for i := range sel.Nodes {
		// For each part of the current verse's text in the html document:
		s := sel.Eq(i)
		textTemp = analyseClass(s, textTemp)
	}
	textFmtd := str(textTemp)
	vFmtd = append(vTemp, textFmtd)
	return
}

func analyseClass(s *goquery.Selection, textTemp []string) (textFmtd []string) {
	class, _ := s.Attr("class")
	woj := strings.Contains(class, "woj")              // Words of Jesus
	footnote := strings.Contains(class, "footnote")    // Footnote
	smallCaps := strings.Contains(class, "small-caps") // Small-caps (e.g., "LORD")
	switch {
	case woj:
		textFmtd = fmtWoJ(s.Contents(), textTemp)
		return
	case footnote:
		textFmtd = fmtFootnote(s, textTemp)
		return
	case smallCaps:
		textFmtd = fmtSmallCaps(s, textTemp)
		return
	default:
		textFmtd = append(textTemp, s.Text())
		return
	}
}

func fmtWoJ(sel *goquery.Selection, textTemp []string) (textFmtd []string) {
	var wojTemp []string
	wojTemp = append(wojTemp, "<FR>") // Opening tag for words of Jesus
	for i := range sel.Nodes {
		// For each html segment of Jesus's words:
		s := sel.Eq(i)
		wojTemp = analyseClass(s, wojTemp) //    Words of Jesus
	}
	wojTemp = append(wojTemp, "<Fr>") // Closing tag for words of Jesus
	wojFmtd := str(wojTemp)
	textFmtd = append(textTemp, wojFmtd)
	return
}

func fmtFootnote(s *goquery.Selection, textTemp []string) (textFmtd []string) {
	var fTemp []string
	fnLetter := s.Find("a").Text()
	fnText := footnoteMap[fnLetter] // Footnote text that was marked earlier.
	fTemp = append(fTemp, "<RF>", fnText, "<Rf>")
	footnoteFmtd := str(fTemp)
	textFmtd = append(textTemp, footnoteFmtd)
	return
}

func fmtSmallCaps(s *goquery.Selection, textTemp []string) (textFmtd []string) {
	textFmtd = append(textTemp, strings.ToUpper(s.Text()))
	return
}

func cleanIndentTags(vTemp []string) (vFmtd []string) {
	var iTemp []string
	var foundFirst bool
	for _, s := range vTemp {
		if strings.Contains(s, "<PI1>") {
			switch foundFirst {
			case false:
				//If line IS the first one in the verse,
				//    keep it single-indented.
				foundFirst = true
				iTemp = append(iTemp, s)
			case true:
				//If line is NOT the first one in the verse,
				//    double its indentation.
				new := strings.Replace(s, "<PI1>", "<PI2>", -1)
				iTemp = append(iTemp, new)
			}
		} else {
			iTemp = append(iTemp, s)
		}
	}
	vFmtd = iTemp
	return
}
