package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

// getWordsForGenerator returns a lot of words in latin. Requires words package
func getWordsForGenerator() ([]string, error) {
	wordsBytes, err := ioutil.ReadFile("/usr/share/dict/words")
	if err != nil {
		return nil, fmt.Errorf("cannot read file /usr/share/dict/words: %s", err)
	}
	return strings.Split(string(wordsBytes), "\n"), nil
}

type generatorSettings struct {
	// number of files3 to generate
	fileNumber int

	// number of words per each file
	wordsNumber int

	// max amount of words to get from dictionary
	maxDictWords int
}

// generateFiles creates dir with dirPath, then generates files3 in this dir with generator settings
func generateFiles(settings generatorSettings, dict []string, dirPath string) error {
	// create dir for random files3
	if err := createDir(dirPath); err != nil {
		return err
	}
	// init random
	rand.Seed(time.Now().UnixNano())
	// set length of dictionary
	n := settings.maxDictWords
	if n == 0 || len(dict) < n {
		n = len(dict)
	}
	// generate files3
	for i := 0; i < settings.fileNumber; i++ {
		fileName := fmt.Sprintf("%s/test-file_%d.txt", dirPath, i)
		f, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("cannot create file '%s' for test: %s", fileName, err)
		}
		// generate words in file
		var res string
		for j := 0; j < settings.wordsNumber; j++ {
			res += dict[rand.Int()%n] + "\n"
		}
		if _, err = f.WriteString(res); err != nil {
			return fmt.Errorf("cannot write to file '%s': %s", fileName, err)
		}
		if err = f.Close(); err != nil {
			return fmt.Errorf("cannot close file '%s': %s", fileName, err)
		}
	}
	return nil
}

func createDir(dirPath string) error {
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("cannot create dir '%s': %s", dirPath, err)
	}
	// change mode of dir (first time it doesn't apply because of umask)
	if err := os.Chmod(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("cannot change mode dir '%s': %s", dirPath, err)
	}
	return nil
}

func BenchmarkMain(b *testing.B) {
	// get dictionary
	dict, err := getWordsForGenerator()
	if err != nil {
		b.Fatal("Error on getting dict:", err)
	}

	// create temp dir for tests
	allTestsDir := os.TempDir() + "/allTests"
	if err = createDir(allTestsDir); err != nil {
		b.Fatal("Error on creating dir for all tests:", err)
	}

	// start benchmarks
	b.Run("simple test", func(b *testing.B) {
		testDir := allTestsDir + "/test1"
		err := generateFiles(generatorSettings{fileNumber: 100, wordsNumber: 200}, dict, testDir)
		if err != nil {
			b.Fatal("Error on generating files3:", err)
		}
		startTest(testDir, b)
	})
	b.Run("a lot of small files3", func(b *testing.B) {
		testDir := allTestsDir + "/test2"
		err := generateFiles(generatorSettings{10000, 100, 50}, dict, testDir)
		if err != nil {
			b.Fatal("Error on generating files3:", err)
		}
		startTest(testDir, b)
	})
	b.Run("a couple of large files3", func(b *testing.B) {
		testDir := allTestsDir + "/test3"
		err := generateFiles(generatorSettings{100, 10000, 1000}, dict, testDir)
		if err != nil {
			b.Fatal("Error on generating files3:", err)
		}
		startTest(testDir, b)
	})
	b.Run("a lot of repeated words", func(b *testing.B) {
		testDir := allTestsDir + "/test4"
		err := generateFiles(generatorSettings{50, 10000, 10}, dict, testDir)
		if err != nil {
			b.Fatal("Error on generating files3:", err)
		}
		startTest(testDir, b)
	})
	b.Run("a lot of unique words", func(b *testing.B) {
		testDir := allTestsDir + "/test4"
		err := generateFiles(generatorSettings{fileNumber: 50, wordsNumber: 10000}, dict, testDir)
		if err != nil {
			b.Fatal("Error on generating files3:", err)
		}
		startTest(testDir, b)
	})

	// clear dir with all test files3
	if err = os.RemoveAll(allTestsDir); err != nil {
		b.Fatal("Error on clearing dir with all tests:", err)
	}
}

// startTest changes command line arguments for a test time, and then
// restores it.
func startTest(dirPath string, b *testing.B) {
	// change arguments for test
	oldArgs := os.Args
	os.Args = []string{os.Args[0], dirPath}
	defer func() { os.Args = oldArgs }()
	// start benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		main()
	}
}
