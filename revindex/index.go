package revindex

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

func Build(texts []string) map[string]Set {
	// todo check if we can specify length
	index := make(map[string]Set)
	for i, text := range texts {
		// add all words to index
		for _, word := range strings.Fields(text) {
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
	return index
}

func Save(index map[string]Set, titles []string, writer io.Writer) error {
	res := make([]byte, 0)
	// save matching of title to index
	for i, title := range titles {
		res = append(res, []byte(fmt.Sprintf("%d:%s\n", i, title))...)
	}
	// save delimiter
	res = append(res, []byte("-\n")...)
	// save index
	for word, keySet := range index {
		keys := keySet.Keys()
		res = append(res, []byte(fmt.Sprintf("%s:%d\n", word, keys))...)
	}
	if _, err := writer.Write(res); err != nil {
		return fmt.Errorf("cannot write index: %s", err)
	}
	return nil
}
