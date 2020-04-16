package main

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/polisgo2020/search-K1ta/revindex"
	"log"
	"net/http"
	"os"
	"sync"
)

// Start server on port for searching phrases in index
func start(indexPath string, port string) {
	// read index
	f, err := os.Open(indexPath)
	if err != nil {
		console.Fatalf("Cannot open file '%s': %s\n", indexPath, err)
	}
	index, err := revindex.Read(f)
	if err != nil {
		console.Fatalf("Cannot read index from file '%s': %s\n", indexPath, err)
	}

	// create server with handler
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("static")))
	mux.Handle("/find", findPhrase(index))
	server := http.Server{Addr: ":" + port, Handler: mux}

	var wg sync.WaitGroup
	// start server
	wg.Add(1)
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			console.Println("Error on closing server:", err)
		} else {
			console.Println("Server is stopped")
		}
		wg.Done()
	}()
	// start console key listener
	wg.Add(1)
	go func() {
		in := bufio.NewReader(os.Stdin)
		for {
			r, _, err := in.ReadRune()
			if err != nil {
				console.Fatal("Cannot read rune from stdin:", err)
			}
			// shutdown server
			if r == 'q' {
				err := server.Shutdown(context.Background())
				if err != nil {
					console.Fatal("Error on server shutdown:", err)
				}
				break
			}
		}
		wg.Done()
	}()
	wg.Wait()
}

func findPhrase(index revindex.Index) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.New(os.Stdout, r.RemoteAddr+" ", log.Lmicroseconds)

		phrase := r.URL.Query().Get("phrase")
		logger.Println("Phrase:", phrase)

		result := index.Find(phrase)
		logger.Println("Result:", result)
		response, err := json.Marshal(result)
		if err != nil {
			logger.Println("Error on marshaling result to json:", err)
			http.Error(w, "Error on marshaling result to json", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(response); err != nil {
			logger.Println("Error on sending response:", err)
			http.Error(w, "Error on sending response", http.StatusInternalServerError)
		}
	}
}
