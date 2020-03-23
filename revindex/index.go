package revindex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	// todo check if we can specify length
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
	for i, title := range index.Titles {
		res = append(res, []byte(fmt.Sprintf("%d:%s\n", i, title))...)
	}
	// save delimiter
	res = append(res, []byte("-\n")...)
	// save index
	for word, keySet := range index.Data {
		keys := keySet.Keys()
		// marshal keys to json to simplify reading
		marshaledKeys, _ := json.Marshal(keys)
		res = append(res, []byte(fmt.Sprintf("%s:%s\n", word, marshaledKeys))...)
	}
	if _, err := writer.Write(res); err != nil {
		return fmt.Errorf("cannot write index: %s", err)
	}
	return nil
}

func Read(reader io.Reader) (Index, error) {
	// todo check if read by lines is better
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
		lineInfo := strings.SplitN(line, ":", 2)
		// todo remove title index from result file 'index.txt'
		index.Titles = append(index.Titles, lineInfo[1])
	}
	// get index itself
	for _, line := range strings.Split(strings.Trim(tokens[1], "\n"), "\n") {
		// get word and indices of texts with it
		lineInfo := strings.Split(line, ":")
		if len(lineInfo) != 2 {
			return Index{}, fmt.Errorf("invalid format of words map in index")
		}
		// unmarshal indices
		var titleIndices []int
		err = json.Unmarshal([]byte(lineInfo[1]), &titleIndices)
		if err != nil {
			return Index{}, fmt.Errorf("cannot unmarshal list with indices: %s", err)
		}
		// put indices to set
		set := Set{}
		set.PutAll(titleIndices)
		index.Data[lineInfo[0]] = set
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
		// for each text title add one entry
		for titleIndex := range titleIndices {
			title := index.Titles[titleIndex]
			entriesMap[title]++
		}
	}
	return entriesMap
}
