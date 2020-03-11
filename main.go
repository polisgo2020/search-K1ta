package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

// Empty struct needs zero bytes
type Void struct{}

// Set to store file indexes
type Set map[int]Void

func (s *Set) Put(val int) {
	(*s)[val] = Void{}
}

func (s *Set) Keys() []int {
	keys := make([]int, 0, len(*s))
	for key := range *s {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	return keys
}

func main() {
	// get files in dir
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Specify path to folder with texts")
	}
	files, err := ioutil.ReadDir(args[1])
	if err != nil {
		log.Fatal(fmt.Sprintf("Error while reading dir '%s': %s", args[1], err))
	}
	// allocate map for index
	index := make(map[string]Set)
	// allocate array for filenames
	filenames := make([]string, len(files))
	// read all files
	for i, fileInfo := range files {
		filename := filepath.Join(args[1], fileInfo.Name())
		// read bytes from file
		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("Error on reading '%s': %s\n", filename, err)
		}
		// add file to names
		filenames[i] = fileInfo.Name()
		// add all words to index
		for _, word := range strings.Fields(string(bytes)) {
			word = strings.TrimFunc(word, func(r rune) bool {
				return !unicode.IsDigit(r) && !unicode.IsLetter(r)
			})
			word = strings.ToLower(word)
			// add word to
			if set, ok := index[word]; ok {
				set.Put(i)
			} else {
				index[word] = Set{i: Void{}}
			}
		}
	}
	// open file for index
	f, err := os.Create("index.txt")
	if err != nil {
		log.Fatal("cannot open file for index:", err)
	}
	// first save filenames
	for i, v := range filenames {
		if _, err = f.Write([]byte(fmt.Sprintf("%d:%s\n", i, v))); err != nil {
			log.Fatal("cannot write filenames to file:", err)
		}
	}
	// add delimiter
	if _, err = f.WriteString("-\n"); err != nil {
		log.Fatal("cannot write filenames to file:", err)
	}
	// save index
	var res []byte
	for k, v := range index {
		res = append(res, []byte(fmt.Sprintf("%s:%d\n", k, v.Keys()))...)
	}
	if _, err = f.Write(res); err != nil {
		log.Fatal("cannot write index to file:", err)
	}
	if err = f.Close(); err != nil {
		log.Fatal("cannot close file with index:", err)
	}
}
