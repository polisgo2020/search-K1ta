package revindex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strings"
	"unicode"
)

type Index struct {
	Titles []string
	Data   map[string]Set
}

func unifyWord(word string) string {
	res := strings.TrimFunc(word, func(r rune) bool {
		return !unicode.IsDigit(r) && !unicode.IsLetter(r)
	})
	res = strings.ToLower(res)
	return res
}

func Build(texts []string, titles []string) (Index, error) {
	if len(texts) != len(titles) {
		return Index{}, errors.New("length of texts is not equal to length of titles")
	}
	index := make(map[string]Set)
	for i, text := range texts {
		// add all words to index
		for _, word := range strings.Fields(text) {
			word = unifyWord(word)
			// add word to
			if set, ok := index[word]; ok {
				set.Put(i)
			} else {
				index[word] = Set{i: Void{}}
			}
		}
	}
	return Index{
		Titles: titles,
		Data:   index,
	}, nil
}

func (index *Index) Save(writer io.Writer) error {
	res := make([]byte, 0)
	// save matching of title to index
	for _, title := range index.Titles {
		res = append(res, []byte(fmt.Sprintf("%s\n", title))...)
	}
	// save delimiter
	res = append(res, []byte("-\n")...)
	// save sorted index
	// sort keys first
	keys := make([]string, 0, len(index.Data))
	for k := range index.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// iterate with sorted keys
	for _, word := range keys {
		indices := index.Data[word]
		sortedIndices := indices.SortedKeys()
		// marshal keys to json to simplify reading
		marshaledIndices, _ := json.Marshal(sortedIndices)
		res = append(res, []byte(fmt.Sprintf("%s:%s\n", word, marshaledIndices))...)
	}
	if _, err := writer.Write(res); err != nil {
		return fmt.Errorf("cannot write index: %s", err)
	}
	return nil
}

func Read(reader io.Reader) (Index, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return Index{}, fmt.Errorf("cannot read index: %s", err)
	}
	// split bytes into titles declarations and index
	tokens := strings.Split(string(bytes), "-\n")
	if len(tokens) != 2 {
		return Index{}, fmt.Errorf("invalid format of index")
	}
	// declare index
	index := Index{Data: make(map[string]Set)}
	// get titles declarations
	for _, line := range strings.Split(strings.Trim(tokens[0], "\n"), "\n") {
		index.Titles = append(index.Titles, line)
	}
	// get index itself
	for _, line := range strings.Split(strings.Trim(tokens[1], "\n"), "\n") {
		// get word and indices of texts with it
		lastColon := strings.LastIndex(line, ":")
		if lastColon == -1 || lastColon == len(line)-1 {
			return Index{}, fmt.Errorf("invalid format of words map in index. Line: %s", line)
		}
		title := line[:lastColon]
		indices := line[lastColon+1:]
		// unmarshal indices
		var titleIndices []int
		err = json.Unmarshal([]byte(indices), &titleIndices)
		if err != nil {
			return Index{}, fmt.Errorf("cannot unmarshal list with indices: %s", err)
		}
		// put indices to set
		set := Set{}
		set.PutAll(titleIndices)
		index.Data[title] = set
	}
	return index, nil
}

func (index *Index) Find(phrase string) map[string]int {
	entriesMap := make(map[string]int)
	for _, word := range strings.Fields(phrase) {
		word = unifyWord(word)
		// get indexes of texts with this word
		titleIndices, exist := index.Data[word]
		if !exist {
			continue
		}
		// for each text titlesToIn add one entry
		for titleIndex := range titleIndices {
			title := index.Titles[titleIndex]
			entriesMap[title]++
		}
	}
	return entriesMap
}
