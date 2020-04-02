package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func BenchmarkGetTextsAndTitlesFromDir(b *testing.B) {
	// functions for creating and deleting directory
	createDir := func(path string) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			b.Fatal(fmt.Sprintf("Cannot create dir '%s': %s", path, err))
		}
		b.Log(fmt.Sprintf("Created dir '%s'", path))
	}
	deleteDir := func(path string) {
		if err := os.RemoveAll(path); err != nil {
			b.Fatal(fmt.Sprintf("Cannot delete dir '%s': %s", path, err))
		}
		b.Log(fmt.Sprintf("Deleted dir '%s'", path))
	}

	// create temp dirs
	smallDirName := filepath.Join(os.TempDir(), fmt.Sprintf("%s-%d", "polisgo", time.Now().UnixNano()))
	createDir(smallDirName)
	defer deleteDir(smallDirName)

	largeDirName := filepath.Join(os.TempDir(), fmt.Sprintf("%s-%d", "polisgo", time.Now().UnixNano()))
	createDir(largeDirName)
	defer deleteDir(largeDirName)

	// create 1000 small files
	for i := 0; i < 5000; i++ {
		filename := filepath.Join(smallDirName, fmt.Sprintf("small-%d.txt", i))
		f, err := os.Create(filename)
		if err != nil {
			b.Fatal(fmt.Sprintf("Cannot create file '%s': %s", filename, err))
		}
		_, err = f.WriteString(strings.Repeat("a", 100))
		if err != nil {
			b.Fatal(fmt.Sprintf("Cannot write to file '%s': %s", filename, err))
		}
	}

	// create 10 large files
	for i := 0; i < 100; i++ {
		filename := filepath.Join(largeDirName, fmt.Sprintf("large-%d.txt", i))
		f, err := os.Create(filename)
		if err != nil {
			b.Fatal(fmt.Sprintf("Cannot create file '%s': %s", filename, err))
		}
		_, err = f.WriteString(strings.Repeat("a", 5000))
		if err != nil {
			b.Fatal(fmt.Sprintf("Cannot write to file '%s': %s", filename, err))
		}
	}

	// start benchmarks
	b.Run("read a lot of small files", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, _, err := getTextsAndTitlesFromDir(smallDirName); err != nil {
				b.Fatal("Error:", err)
			}
		}
	})

	b.Run("read a few of large files", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, _, err := getTextsAndTitlesFromDir(largeDirName); err != nil {
				b.Fatal("Error:", err)
			}
		}
	})
}
