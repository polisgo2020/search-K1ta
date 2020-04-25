package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	_ "github.com/lib/pq"
	"github.com/polisgo2020/search-K1ta/database"
	"github.com/polisgo2020/search-K1ta/revindex"
	"github.com/polisgo2020/search-K1ta/server"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// app config
type Config struct {
	Addr         string `env:"POLISGO_ADDR" envDefault:"localhost:8080"`
	Hostname     string `env:"DB_HOSTNAME" envDefault:"localhost"`
	Hostport     string `env:"DB_HOSTPORT" envDefault:"5432"`
	Username     string `env:"DB_USERNAME" envDefault:"postgres"`
	Password     string `env:"DB_PASSWORD" envDefault:"postgres"`
	DatabaseName string `env:"DB_NAME" envDefault:"postgres"`
}

// logger for console
var console = log.New(os.Stdout, "", 0)
var cfg Config

func main() {
	if err := env.Parse(&cfg); err != nil {
		logrus.Fatal("Error on parsing config:", err)
	}
	app := &cli.App{
		Usage: "Tool for creating an index on texts and searching phrases in it",
		Commands: []*cli.Command{
			{
				Name:    "build",
				Aliases: []string{"b"},
				Usage:   "Build index by files in dir",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "clear",
						Aliases: []string{"c"},
						Usage:   "clear database before saving index",
						Value:   false,
					},
				},
				ArgsUsage: "<dir>",
				Action: func(ctx *cli.Context) error {
					dir := ctx.Args().Get(0)
					if dir == "" {
						console.Fatal("Specify dir with files")
					}
					clearDb := ctx.Bool("clear")
					build(dir, clearDb)
					return nil
				},
			},
			{
				Name:      "find",
				Aliases:   []string{"f"},
				Usage:     "Find phrase in specified index",
				ArgsUsage: "\"<phrase>\"",
				Action: func(ctx *cli.Context) error {
					phrase := ctx.Args().Get(0)
					findInDb(phrase)
					return nil
				},
			},
			{
				Name:        "start",
				Aliases:     []string{"s"},
				Usage:       "Start server for searching phrases. Main page is on /",
				Description: "Env variable for server addr: POLISGO_ADDR=ADDR. Default is localhost:8080",
				Action: func(ctx *cli.Context) error {
					// connect to db
					db, err := database.Connect(cfg.Hostname, cfg.Hostport, cfg.Username, cfg.Password, cfg.DatabaseName)
					if err != nil {
						logrus.Fatal("Error on connecting to database:", err)
					}
					defer func() {
						err = db.Close()
						if err != nil {
							logrus.Fatal("Error on closing connection to database:", err)
						}
					}()
					return server.Start(cfg.Addr, db)
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
		return nil, nil, fmt.Errorf("error while reading dir '%s': %w", dirPath, err)
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
func build(dir string, clearDb bool) {
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
	// save to db
	db, err := database.Connect(cfg.Hostname, cfg.Hostport, cfg.Username, cfg.Password, cfg.DatabaseName)
	if err != nil {
		console.Fatal("Error on connecting to database:", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			console.Fatal("Error on closing connection to database:", err)
		}
	}()
	// clear db if we need
	if clearDb {
		console.Println("Clearing database")
		err = db.DropAll()
		if err != nil {
			console.Fatal("Error on clearing db:", err)
		}
	}
	err = db.Init()
	if err != nil {
		console.Fatal("Error on init db:", err)
	}
	err = index.SaveToDb(db)
	if err != nil {
		console.Fatal("Error on saving index to db:", err)
	}
}

func findInDb(phrase string) {
	db, err := database.Connect(cfg.Hostname, cfg.Hostport, cfg.Username, cfg.Password, cfg.DatabaseName)
	if err != nil {
		console.Fatal("Error on connecting to database:", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			console.Fatal("Error on closing connection to database:", err)
		}
	}()
	res, err := revindex.FindInDb(phrase, db)
	if err != nil {
		console.Fatal("Cannot find phrase in db:", err)
	}
	if len(res) == 0 {
		console.Println("No entries")
		return
	}
	console.Println("Entries:")
	for title, amount := range res {
		console.Printf("%s; entries: %d\n", title, amount)
	}
}
