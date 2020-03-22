package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"polisgo/revindex"
)

func main() {
	// get files in dir
	args := os.Args
	if len(os.Args) < 2 {
		log.Fatal("Specify path to folder with texts")
	}
	files, err := ioutil.ReadDir(args[1])
	if err != nil {
		log.Fatal(fmt.Sprintf("Error while reading dir '%s': %s", args[1], err))
	}
	// get texts and titles
	texts, titles := getTextsAndTitlesFromFiles(args[1], files)
	// build index
	index := revindex.Build(texts)
	// open file for index
	f, err := os.Create("index.txt")
	if err != nil {
		log.Fatal("cannot open file for index:", err)
	}
	// save index
	err = revindex.Save(index, titles, f)
	if err != nil {
		log.Fatal("Error on writing index to file:", err)
	}
	// close file
	if err = f.Close(); err != nil {
		log.Fatal("cannot close file with index:", err)
	}
}

func getTextsAndTitlesFromFiles(path string, files []os.FileInfo) ([]string, []string) {
	texts := make([]string, 0, len(files))
	titles := make([]string, 0, len(files))
	for _, fileInfo := range files {
		filename := filepath.Join(path, fileInfo.Name())
		// read bytes from file
		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("Error on reading '%s': %s\n", filename, err)
		}
		texts = append(texts, string(bytes))
		titles = append(titles, fileInfo.Name())
	}
	return texts, titles
}
