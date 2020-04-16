package main

import (
	"fmt"
	"github.com/polisgo2020/search-K1ta/revindex"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// logger for console
var console = log.New(os.Stdout, "", 0)

func main() {
	app := &cli.App{
		Usage: "Tool for creating an index on texts and searching phrases in it",
		Commands: []*cli.Command{
			{
				Name:      "build",
				Aliases:   []string{"b"},
				Usage:     "Build index by files in dir",
				ArgsUsage: "<dir>",
				Action: func(ctx *cli.Context) error {
					dir := ctx.Args().Get(0)
					if dir == "" {
						console.Fatal("Specify dir with files")
					}
					build(dir)
					return nil
				},
			},
			{
				Name:      "find",
				Aliases:   []string{"f"},
				Usage:     "Find phrase in specified index",
				ArgsUsage: "<index_file> \"<phrase>\"",
				Action: func(ctx *cli.Context) error {
					index := ctx.Args().Get(0)
					if index == "" {
						console.Fatal("specify index file")
					}
					phrase := ctx.Args().Get(1)
					find(index, phrase)
					return nil
				},
			},
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "Start server for searching phrases in the specified index on port",
				Description: "Server API:\n" +
					"GET /find?phrase=<phrase> - find phrase in index. Response is json from find function\n" +
					"GET / - returns main page",
				ArgsUsage: "<index_file> <port>",
				Action: func(ctx *cli.Context) error {
					index := ctx.Args().Get(0)
					if index == "" {
						console.Fatal("specify index file")
					}
					port := ctx.Args().Get(1)
					if port == "" {
						console.Fatal("specify server port")
					}
					start(index, port)
					return nil
				},
			},
		},
		HideVersion: true,
	}
	err := app.Run(os.Args)
	if err != nil {
		console.Fatal(err)
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
