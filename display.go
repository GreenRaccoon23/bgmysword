package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
	//"github.com/PuerkitoBio/goquery"
	//"github.com/fatih/color"
)

var (
	Red      = color.New(color.FgRed)
	Blue     = color.New(color.FgBlue)
	Green    = color.New(color.FgGreen)
	Magenta  = color.New(color.FgMagenta)
	White    = color.New(color.FgWhite)
	Black    = color.New(color.FgBlack)
	BRed     = color.New(color.FgRed, color.Bold)
	BBlue    = color.New(color.FgBlue, color.Bold)
	BGreen   = color.New(color.FgGreen, color.Bold)
	BMagenta = color.New(color.FgMagenta, color.Bold)
	BWhite   = color.New(color.Bold, color.FgWhite)
	BBlack   = color.New(color.Bold, color.FgBlack)
)

func imgINRI() {
	image := []string{
		"              .======.              ",
		"              | INRI |              ",
		"              |      |              ",
		"              |      |              ",
		"     .========'      '========.     ",
		"     |   _      xxxx      _   |     ",
		"     |  /_;-.__ / _\\  _.-;_\\  |     ",
		"     |     `-._`'`_/'`.-'     |     ",
		"     '========.`\\   /`========'     ",
		"              | |  / |              ",
		"              |/-.(  |              ",
		"              |\\_._\\ |              ",
		"              | \\ \\`;|              ",
		"              |  > |/|              ",
		"              | / // |              ",
		"              | |//  |              ",
		"              | \\(\\  |              ",
		"              |  ``  |              ",
		"              |      |              ",
		"              |      |              ",
		"              |      |              ",
		"              |      |              ",
		"  \\\\jgs _  _\\\\| \\//  |//_   _ \\// _ ",
		" ^ `^`^ ^`` `^ ^` ``^^`  `^^` `^ `^ ",
	}
	PrintArt(Red, image)
}

func Break(c *color.Color, s string) {
	c.Println(strings.Repeat(s, 79))
}

func Line(c *color.Color) {
	Break(c, "-")
}

func BLine(c *color.Color) {
	Break(c, "=")
}

func PrintCenter(c *color.Color, t string) {
	w := 39 - len(t)/2
	s := strings.Repeat(" ", w)
	c.Printf("%v%v%v\n", s, t, s)
}

func PrintCenterLines(c *color.Color, st string) {
	sl := strings.Split(st, "\n")
	for _, s := range sl {
		PrintCenter(c, s)
	}
}

func PrintCenterUnknown(c *color.Color, st string) {
	switch {
	case len(st) > 78:
		sl := strings.SplitAfter(st, " ")
		half := len(sl) / 2
		stBeg := strings.TrimSpace(str(sl[:half]))
		stEnd := strings.TrimSpace(str(sl[half:]))
		PrintCenter(c, stBeg)
		PrintCenter(c, stEnd)
	default:
		PrintCenter(c, st)
	}
}

func PrintArt(c *color.Color, sl []string) {
	for _, st := range sl {
		PrintCenterLines(c, st)
	}
}

func input(c *color.Color) string {
	c.Set()
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	return answer
}

func inputTranslation() string {
	Red.Println(
		"No translation specified.",
	)
	Black.Println(
		"(You can specify the translation by typing 'bgmysword XXX')",
	)
	Blue.Println(
		"Enter the abbreviation of the translation to download",
	)
	Blue.Println(
		"  using the same abbreviation that biblegateway.com uses.",
	)
	Black.Println(
		"(They're listed at biblegateway.com/versions)",
	)
	BWhite.Printf(
		"==> ",
	)
	return input(BGreen)
}

func copyrightAccept(copyrightInfo, publisherInfo string) {
	Line(BWhite)
	PrintCenter(BWhite,
		"License Agreement")
	Line(BWhite)
	PrintCenter(Blue, fmt.Sprintf(
		`This program will generate a MySword bible module for the %v Bible`,
		translation,
	))
	PrintCenter(Blue, fmt.Sprintf(
		`by looking up the %v Bible on https://www.biblegateway.org/`,
		translation,
	))
	fmt.Printf("\n")
	PrintCenter(BMagenta, fmt.Sprintf(
		`BibleGateway License for %v`,
		translation,
	))
	PrintCenter(Magenta, transName)
	PrintCenterUnknown(Magenta, copyrightInfo)
	fmt.Printf("\n")
	PrintCenter(BRed,
		"Agreement",
	)
	PrintCenterLines(Red, fmt.Sprintf(
		`  By using this program, I declare that my local copyright law does not forbid
the use of this website for a purpose other than webbrowsing. I will not
distribute any module generated by this program unless explicitely permitted by
%v
I accept full responsibility for all copyright violations that
are the result of using this program and absolve the author of any
responsibility for them.
I understand that I may distribute and/or modify this program itself,
but not any modules generated by it.`,
		publisherInfo,
	))
	PrintCenter(Black,
		"In other words, don't abuse this program.",
	)
	PrintCenter(Red,
		`If you wish to accept the statements above enter the exact word "accept."`,
	)
	PrintCenter(Red,
		"Then press [ENTER]. ",
	)
	BWhite.Printf(
		"==> ",
	)
	ansAgree := input(BGreen)
	accepted := strings.EqualFold(ansAgree, "accept\n")
	switch accepted {
	case false:
		Red.Println("License not accepted.")
		color.Unset()
		os.Exit(0)
	}
}

func progressTranslation() {
	BLine(BWhite)
	BWhite.Println(transName)
}

func progressBook(title, titleSpacing, numSpacing string, currNum, totalNum int) {
	BGreen.Printf("\r=>%v%v%v%d/%d\n",
		title, titleSpacing, numSpacing, currNum, totalNum)
}

func progressChapter(current, total int) {
	Blue.Printf("\r ->%d/%d", current, total)
}