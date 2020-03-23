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
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Not enough arguments. Run ", args[0], " help")
	}
	switch args[1] {
	case "build":
		// get files in dir
		if len(args) < 3 {
			log.Fatal("Specify path to folder with texts")
		}
		files, err := ioutil.ReadDir(args[2])
		if err != nil {
			log.Fatal(fmt.Sprintf("Error while reading dir '%s': %s", args[1], err))
		}
		// get texts and titles
		texts, titles := getTextsAndTitlesFromFiles(args[2], files)
		// build index
		index, err := revindex.Build(texts, titles)
		if err != nil {
			log.Fatal("Error on building index:", err)
		}
		// open file for index
		f, err := os.Create("index.txt")
		if err != nil {
			log.Fatal("cannot open file for index:", err)
		}
		// save index
		err = index.Save(f)
		if err != nil {
			log.Fatal("Error on writing index to file:", err)
		}
		// close file
		if err = f.Close(); err != nil {
			log.Fatal("cannot close file with index:", err)
		}
	case "find":
		if len(args) < 4 {
			log.Fatal("Not enough arguments")
		}
		f, err := os.Open(args[2])
		if err != nil {
			log.Fatalf("Cannot open file '%s': %s\n", args[2], err)
		}
		// read index
		index, err := revindex.Read(f)
		if err != nil {
			log.Fatalf("Cannot read index from file '%s': %s\n", args[2], err)
		}
		// find words from phrase
		res := index.Find(args[3])
		if len(res) == 0 {
			log.Println("No entries")
			return
		}
		log.Println("Entries:")
		for title, amount := range res {
			fmt.Printf("%s; entries: %d\n", title, amount)
		}
	case "help":
		log.Println("Tool for creating index by texts and find phrases in it")
		log.Println("Usage:")
		log.Println(args[0], "build <dir>			builds index by files in dir")
		log.Println(args[0], "find <index> \"<prhase>\"		find phrase in specified index")
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
