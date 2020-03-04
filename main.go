package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
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
	// read all files
	for i, fileInfo := range files {
		filename := args[1] + "/" + fileInfo.Name()
		// read bytes from file
		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("Error on reading '%s': %s\n", filename, err)
		}
		// add all words to index
		for _, word := range strings.Fields(string(bytes)) {
			if set, ok := index[word]; ok {
				set.Put(i)
			} else {
				index[word] = Set{i: Void{}}
			}
		}
	}
	// save index to file
	f, err := os.Create("index.txt")
	if err != nil {
		log.Fatal("cannot open file for index:", err)
	}
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
