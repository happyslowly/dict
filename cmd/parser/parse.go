package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/happyslowly/dict/index"
	"github.com/happyslowly/dict/word"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	dictFile, err := os.Open("dict.dat")
	if err != nil {
		log.Fatal(err)
	}
	defer dictFile.Close()

	outputFile, err := os.Create("dict.json")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	encodedFile, err := os.Create("dict.bin")
	if err != nil {
		log.Fatal(err)
	}
	defer encodedFile.Close()

	err = parse(dictFile, outputFile, encodedFile)
	if err != nil {
		log.Fatal(err)
	}
}

func parse(f io.Reader, jf *os.File, ef io.Writer) error {
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	idxRecords := make(map[string]index.Record, 0)

	start := 0
	for scanner.Scan() {
		w := parseWord(strings.NewReader(scanner.Text()))
		json, err := json.Marshal(w)

		if err != nil {
			return err
		}

		_, err = jf.WriteString(string(json) + "\n")
		if err != nil {
			return err
		}

		n, err := encodeWord(w, ef)
		if err != nil {
			return err
		}

		var ir = index.Record{Title: w.Title, Pos: start, Offset: start + n}
		idxRecords[w.Title] = ir

		start += n
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	err := encodeIndex(idxRecords)
	if err != nil {
		return err
	}

	return nil

}

func parseWord(reader io.Reader) word.Word {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	var wd word.Word
	wd.Definitions = make([]word.Definition, 0)
	wd.Classes = make([]string, 0)

	s := doc.Find(".entry").First()
	if title, ok := s.Attr("d:title"); ok {
		wd.Title = title
	}

	doc.Find(".ph").Each(func(i int, s *goquery.Selection) {
		if dialect, ok := s.Attr("dialect"); ok {
			if dialect == "BrE" {
				wd.ProBrE = strings.TrimSpace(s.Text())
			} else if dialect == "AmE" {
				wd.ProAmE = strings.TrimSpace(s.Text())
			}
		}
	})

	doc.Find(".gramb").Each(func(i int, s *goquery.Selection) {
		ps := strings.TrimSpace(s.Find(".ps").Text())

		var def word.Definition
		def.Class = ps
		def.Translations = make([]string, 0)

		s.Find(".trans").Each(func(i int, s *goquery.Selection) {
			if _, ok := s.Attr("d:def"); ok {
				def.Translations = append(def.Translations, strings.TrimSpace(s.Text()))
			}
		})

		if len(def.Translations) > 0 {
			wd.Definitions = append(wd.Definitions, def)
			wd.Classes = append(wd.Classes, ps)
		}
	})

	return wd
}

func encodeWord(wd word.Word, writer io.Writer) (int, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(wd); err != nil {
		return 0, err
	}

	n, err := writer.Write(buf.Bytes())

	if err != nil {
		return 0, err
	}

	return n, nil
}

func encodeIndex(idxRecords map[string]index.Record) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(idxRecords); err != nil {
		return err
	}

	f, err := os.Create("dict.idx")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}
