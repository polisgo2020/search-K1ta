package main

import (
	"fmt"
	"github.com/polisgo2020/search-K1ta/revindex"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
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
		// get texts and titles
		texts, titles, err := getTextsAndTitlesFromDir(args[2])
		if err != nil {
			log.Fatal("Error:", err)
		}
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

func getTextsAndTitlesFromDir(dirPath string) ([]string, []string, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error while reading dir '%s': %s", dirPath, err)
	}
	texts := make([]string, 0, len(files))
	titles := make([]string, 0, len(files))
	var mux sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, fileInfo := range files {
		go func(fileInfo os.FileInfo) {
			filename := filepath.Join(dirPath, fileInfo.Name())
			// read bytes from file
			bytes, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Printf("Error on reading '%s': %s\n", filename, err)
			}
			mux.Lock()
			texts = append(texts, string(bytes))
			titles = append(titles, fileInfo.Name())
			mux.Unlock()
			wg.Done()
		}(fileInfo)
	}
	wg.Wait()
	return texts, titles, nil
}
