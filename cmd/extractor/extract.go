package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %v <input file>", os.Args[0])
		os.Exit(1)
	}

	err := extract(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

func extract(ifile string) error {
	f, err := os.Open(ifile)
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	for len(content) > 0 {
		data, err := decompress(content)
		if err == nil {
			for i := range data {
				fmt.Println(data[i])
			}
		}
		content = content[1:]
	}

	return nil
}

func decompress(content []byte) ([]string, error) {
	reader, err := zlib.NewReader(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	raw, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	data := make([]string, 0)
	for i, j := 0, -1; i < len(raw); i += 1 {
		if raw[i] == '\n' && j+5 < i {
			data = append(data, string(raw[5+j:i]))
			j = i
		}
	}

	return data, nil
}
