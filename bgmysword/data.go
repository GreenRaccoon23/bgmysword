package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"time"
	//"github.com/PuerkitoBio/goquery"
	//"github.com/fatih/color"
)

type BibleArchive struct {
	Index        int
	Book         string
	ChapterRange int
}

var (
	db              *sql.DB
	tx              *sql.Tx
	sqlStmtInsBible string         = `insert into Bible values(?,?,?,?)`
	Bible           []BibleArchive = []BibleArchive{
		{1, "Genesis", 50},
		{2, "Exodus", 40},
		{3, "Leviticus", 27},
		{4, "Numbers", 36},
		{5, "Deuteronomy", 34},
		{6, "Joshua", 24},
		{7, "Judges", 21},
		{8, "Ruth", 4},
		{9, "1+Samuel", 31},
		{10, "2+Samuel", 24},
		{11, "1+Kings", 22},
		{12, "2+Kings", 25},
		{13, "1+Chronicles", 29},
		{14, "2+Chronicles", 36},
		{15, "Ezra", 10},
		{16, "Nehemiah", 13},
		{17, "Esther", 10},
		{18, "Job", 42},
		{19, "Psalms", 150},
		{20, "Proverbs", 31},
		{21, "Ecclesiastes", 12},
		{22, "Song+of+Solomon", 8},
		{23, "Isaiah", 66},
		{24, "Jeremiah", 52},
		{25, "Lamentations", 5},
		{26, "Ezekiel", 48},
		{27, "Daniel", 12},
		{28, "Hosea", 14},
		{29, "Joel", 3},
		{30, "Amos", 9},
		{31, "Obadiah", 1},
		{32, "Jonah", 4},
		{33, "Micah", 7},
		{34, "Nahum", 3},
		{35, "Habakkuk", 3},
		{36, "Zephaniah", 3},
		{37, "Haggai", 2},
		{38, "Zechariah", 14},
		{39, "Malachi", 4},
		{40, "Matthew", 28},
		{41, "Mark", 16},
		{42, "Luke", 24},
		{43, "John", 21},
		{44, "Acts", 28},
		{45, "Romans", 16},
		{46, "1+Corinthians", 16},
		{47, "2+Corinthians", 13},
		{48, "Galatians", 6},
		{49, "Ephesians", 6},
		{50, "Philippians", 4},
		{51, "Colossians", 4},
		{52, "1+Thessalonians", 5},
		{53, "2+Thessalonians", 3},
		{54, "1+Timothy", 6},
		{55, "2+Timothy", 4},
		{56, "Titus", 3},
		{57, "Philemon", 1},
		{58, "Hebrews", 13},
		{59, "James", 5},
		{60, "1+Peter", 5},
		{61, "2+Peter", 3},
		{62, "1+John", 5},
		{63, "2+John", 1},
		{64, "3+John", 1},
		{65, "Jude", 1},
		{66, "Revelation", 22},
	}
)

func genBible() {
	ext := ".bbl.mybible"
	fileName := str([]string{translation, ext})
	if _, err := os.Stat(fileName); err == nil {
		os.Remove(fileName)
		Black.Println("Removed", fileName)
	}
	var err error
	db, err = sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}
	crtDetails := `
	create table Details
	(Description NVARCHAR(255), Abbreviation NVARCHAR(50),
	Comments TEXT, Version TEXT, VersionDate DATETIME,
	PublishDate DATETIME, RightToLeft BOOL,
	OT BOOL, NT BOOL, Strong BOOL)
	`
	sqlCrt(crtDetails)
	crtBible := `
	create table Bible
    (Book INT, Chapter INT,
    Verse INT, Scripture TEXT)
	`
	sqlCrt(crtBible)
	tx, err = db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	sqlInsDetails()
}

func sqlCrt(sqlStmt string) {
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	return
}

func sqlInsDetails() {
	insDetails := `
	insert into Details(
		Description, Abbreviation, Comments,
	 	Version, VersionDate, PublishDate, RightToLeft,
	 	OT, NT, Strong)
	values(?,?,?,?,?,?,?,?,?,?)
	`
	stmt, err := tx.Prepare(insDetails)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(
		transName, translation, "None", "1.0",
		time.Now().String(), "Unknown",
		"false", "true", "true", "false",
	)
	return
}

func sqlInsBible(book, chapter string, verse int, text string) {
	stmt, err := tx.Prepare(sqlStmtInsBible)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(
		book, chapter, verse, text,
	)
}
