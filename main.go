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

// logger for console
var console = log.New(os.Stdout, "", 0)

func main() {
	// get args
	args := os.Args
	// create func for length validation
	exitIfNotEnoughArgs := func(length int) {
		if len(args) < length {
			console.Fatal("Not enough arguments. Run ", args[0], " help")
		}
	}

	// exit if we don't have command
	exitIfNotEnoughArgs(2)

	switch args[1] {
	case "build":
		exitIfNotEnoughArgs(3)
		build(args[2])
	case "find":
		exitIfNotEnoughArgs(4)
		find(args[2], args[3])
	case "start":
		exitIfNotEnoughArgs(4)
		start(args[2], args[3])
	case "help":
		console.Println("Tool for creating index by texts and find phrases in it")
		console.Println("Usage:")
		console.Println(args[0], "build <dir>")
		console.Println("\tbuild index by files in dir")
		console.Println()
		console.Println(args[0], "find <index> \"<phrase>\"")
		console.Println("\tfind phrase in specified index")
		console.Println()
		console.Println(args[0], "start <index> <port>")
		console.Println("\tstart server for searching phrases in the specified index on port. console.Fatal by pressing 'q'")
		console.Println("\tServer API:")
		console.Println("\tGET /find?phrase=<phrase> - find phrase in index. Response is json from find function")
		console.Println("\tGET / - returns main page")
	}
}

// Build index from files in dir and save it to file "index.txt"
func build(dir string) {
	// get texts and titles
	texts, titles, err := getTextsAndTitlesFromDir(dir)
	if err != nil {
		console.Fatal("Error:", err)
	}
	// build index
	index, err := revindex.Build(texts, titles)
	if err != nil {
		console.Fatal("Error on building index:", err)
	}
	// open file for index
	f, err := os.Create("index.txt")
	if err != nil {
		console.Fatal("cannot open file for index:", err)
	}
	// save index
	err = index.Save(f)
	if err != nil {
		console.Fatal("Error on writing index to file:", err)
	}
	// close file
	if err = f.Close(); err != nil {
		console.Fatal("cannot close file with index:", err)
	}
}

// Get texts and titles from dir
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
				console.Printf("Error on reading '%s': %s\n", filename, err)
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

// Find phrase in index
func find(indexPath string, phrase string) {
	// read index
	f, err := os.Open(indexPath)
	if err != nil {
		console.Fatalf("Cannot open file '%s': %s\n", indexPath, err)
	}
	index, err := revindex.Read(f)
	if err != nil {
		console.Fatalf("Cannot read index from file '%s': %s\n", indexPath, err)
	}
	// find words from phrase
	res := index.Find(phrase)
	if len(res) == 0 {
		console.Println("No entries")
		return
	}
	console.Println("Entries:")
	for title, amount := range res {
		console.Printf("%s; entries: %d\n", title, amount)
	}
}
