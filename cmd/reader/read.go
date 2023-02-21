package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/happyslowly/dict/index"
	"github.com/happyslowly/dict/word"
)

const (
	BLACK     = "\033[30m"
	RED       = "\033[31m"
	GREEN     = "\033[32m"
	YELLOW    = "\033[33m"
	BLUE      = "\033[34m"
	MAGENTA   = "\033[35m"
	CYAN      = "\033[36m"
	WHITE     = "\033[37m"
	UNDERLINE = "\033[4m"
	BOLD      = "\033[1m"
	RESET     = "\033[0m"
)

func main() {
	if len(os.Args) < 2 {
		return
	}
	s := os.Args[1]

	indice, err := loadIndex()
	if err != nil {
		log.Fatal(err)
	}

	if i, ok := indice[s]; ok {
		w, err := query(i)
		if err != nil {
			log.Fatal(err)
		}
		render(w)
	} else {
		fmt.Println("Not found!")
	}
}

func loadIndex() (map[string]index.Record, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadFile(home + "/development/projects/dict/dict.idx")
	if err != nil {
		return nil, err
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(buf))

	idx := make(map[string]index.Record)

	if err := decoder.Decode(&idx); err != nil {
		return nil, err
	}

	return idx, nil
}

func query(i index.Record) (word.Word, error) {
	var w word.Word

	home, err := os.UserHomeDir()
	if err != nil {
		return w, err
	}

	f, err := os.Open(home + "/development/projects/dict/dict.bin")
	if err != nil {
		return w, err
	}

	f.Seek(int64(i.Pos), 0)

	buf := make([]byte, i.Offset)

	_, err = f.Read(buf)
	if err != nil {
		return w, err
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(buf))
	err = decoder.Decode(&w)
	if err != nil {
		return w, err
	}

	return w, nil
}

func render(w word.Word) {
	fmt.Printf("%v%v:%v\n", BOLD, w.Title, RESET)
	if w.ProAmE != "" {
		fmt.Printf("%vAmE:[%v]%v ", YELLOW, w.ProAmE, RESET)
	}
	if w.ProBrE != "" {
		fmt.Printf("%vBrE:[%v]%v\n", YELLOW, w.ProBrE, RESET)
	}
	for i := range w.Definitions {
		d := w.Definitions[i]
		fmt.Printf("%v%v.%v\n", BLUE+BOLD, d.Class, RESET)
		for j := range d.Translations {
			t := d.Translations[j]
			fmt.Printf("  %v%v%v\n", BLUE, t, RESET)
		}
	}
}
