package revindex

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"testing"
)

type Dictionary []string

// getDictionary returns a lot of words in latin. Requires words package
func getDictionary() (Dictionary, error) {
	wordsBytes, err := ioutil.ReadFile("/usr/share/dict/words")
	if err != nil {
		return nil, fmt.Errorf("cannot read file /usr/share/dict/words: %s", err)
	}
	return strings.Split(string(wordsBytes), "\n"), nil
}

type generatorSettings struct {
	textsNumber  int
	wordsPerText int
	maxDictSize  int
}

func (d Dictionary) generateTextsAndTitles(s generatorSettings) ([]string, []string) {
	texts := make([]string, s.textsNumber)
	titles := make([]string, s.textsNumber)
	// init random with static number so we can get same number sequence each time
	r := rand.New(rand.NewSource(100))
	// set length of dictionary
	n := s.maxDictSize
	if n == 0 || len(d) < n {
		n = len(d)
	}
	rands := make([]int, n)
	for i := 0; i < n; i++ {
		rands[i] = r.Intn(len(d))
	}
	// generate texts
	for i := 0; i < s.textsNumber; i++ {
		titles[i] = fmt.Sprintf("test-file_%d", i)
		startIndex := i * s.wordsPerText
		for j := 0; j < s.wordsPerText; j++ {
			texts[i] += d[rands[(startIndex+j)%n]] + "\n"
		}
	}
	return texts, titles
}

func BenchmarkBuild_SimpleTest(b *testing.B) {
	// get words
	dict, err := getDictionary()
	if err != nil {
		b.Fatal("Error on getting words for generator:", err)
	}

	// func for start testing
	goTest := func(b *testing.B, texts []string, titles []string) {
		for i := 0; i < b.N; i++ {
			if _, err := Build(texts, titles); err != nil {
				b.Fatal(err)
			}
		}
	}

	texts, titles := dict.generateTextsAndTitles(generatorSettings{
		textsNumber: 10, wordsPerText: 50000, maxDictSize: 50,
	})
	b.Run("few texts, lot identical words", func(b *testing.B) { goTest(b, texts, titles) })

	texts, titles = dict.generateTextsAndTitles(generatorSettings{
		textsNumber: 10, wordsPerText: 50000,
	})
	b.Run("few texts, lot different words", func(b *testing.B) { goTest(b, texts, titles) })

	texts, titles = dict.generateTextsAndTitles(generatorSettings{
		textsNumber: 1000, wordsPerText: 1000, maxDictSize: 50,
	})
	b.Run("lot texts, lot identical words", func(b *testing.B) { goTest(b, texts, titles) })

	texts, titles = dict.generateTextsAndTitles(generatorSettings{
		textsNumber: 1000, wordsPerText: 1000,
	})
	b.Run("lot texts, lot different words", func(b *testing.B) { goTest(b, texts, titles) })

	texts, titles = dict.generateTextsAndTitles(generatorSettings{
		textsNumber: 50000, wordsPerText: 10, maxDictSize: 50,
	})
	b.Run("lot texts, few identical words", func(b *testing.B) { goTest(b, texts, titles) })

	texts, titles = dict.generateTextsAndTitles(generatorSettings{
		textsNumber: 50000, wordsPerText: 10,
	})
	b.Run("lot texts, few different words", func(b *testing.B) { goTest(b, texts, titles) })
}

// Generates Index with this format:
//title-with-number-1
//...
//<textsNumber>
//-
//w1:[0,1,...,<textsNumber>]
//w2:[0,1,...,<textsNumber>]
func generateIndex(textsNumber int, wordsNumber int) *Index {
	titles := make([]string, textsNumber)
	entries := make(map[string]Set)
	for i := 0; i < textsNumber; i++ {
		titles[i] = fmt.Sprintf("title-with-number-%d", i)
	}
	for i := 0; i < wordsNumber; i++ {
		set := Set{}
		for j := 0; j < textsNumber; j++ {
			set.Put(j)
		}
		entries[fmt.Sprintf("w%d", i)] = set
	}
	return &Index{
		Titles: titles,
		Data:   entries,
	}
}

func BenchmarkIndex_Find(b *testing.B) {
	b.Run("lot texts, few words, 100-words phrase", func(b *testing.B) {
		index := generateIndex(1000, 10000)
		phrase := "w0"
		for i := 1; i < 100; i++ {
			phrase = phrase + " " + fmt.Sprintf("w%d", i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = index.Find(phrase)
		}
	})

	b.Run("lot texts, few words, 500-words phrase", func(b *testing.B) {
		index := generateIndex(1000, 10000)
		phrase := "w0"
		for i := 1; i < 500; i++ {
			phrase = phrase + " " + fmt.Sprintf("w%d", i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = index.Find(phrase)
		}
	})

	b.Run("few texts, lot words, 100-words phrase", func(b *testing.B) {
		index := generateIndex(10000, 1000)
		phrase := "w0"
		for i := 1; i < 100; i++ {
			phrase = phrase + " " + fmt.Sprintf("w%d", i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = index.Find(phrase)
		}
	})

	b.Run("few texts, lot words, 500-words phrase", func(b *testing.B) {
		index := generateIndex(10000, 1000)
		phrase := "w0"
		for i := 1; i < 500; i++ {
			phrase = phrase + " " + fmt.Sprintf("w%d", i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = index.Find(phrase)
		}
	})
}
