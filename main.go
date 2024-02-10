package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/manifoldco/promptui"
)

const (
	YES             = "Yes"
	NO              = "No"
	YES_ALL         = "Yes, all"
	TEXT_FILES_ONLY = "Text files only"
)

var (
	bom = []byte{0xEF, 0xBB, 0xBF}
)

func panicOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func isProbablyText(b []byte) bool {
	s := strings.TrimRight(string(b), string(0))
	bytes := []byte(s)
	mimeType := http.DetectContentType(bytes)
	return strings.HasPrefix(mimeType, "text/")
}

func addBom(filePath string) error {
	b, err := os.ReadFile(filePath)
	if err != nil {
		err := fmt.Errorf("can not open file: %w", err)
		return err
	}

	if bytes.Equal(b[0:3], bom) {
		err := fmt.Errorf("already has a BOM")
		return err
	}

	if !utf8.Valid(b) {
		err := fmt.Errorf("not valid UTF-8 encoded")
		return err
	}

	f, err := os.Create(filePath)
	if err != nil {
		err := fmt.Errorf("can not open file: %w", err)
		return err
	}
	defer f.Close()

	_, err = f.Write(bom)
	if err != nil {
		err := fmt.Errorf("can not write file: %w", err)
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		err := fmt.Errorf("can not write file: %w", err)
		return err
	}

	return nil
}

func main() {
	if len(os.Args) <= 1 {
		err := fmt.Errorf("one or more arguments are required")
		panicOnError(err)
	}

	filePathArgs := os.Args[1:]
	var filePathsClean []string
	for _, filePathArg := range filePathArgs {
		filePathClean := filepath.Clean(filePathArg)
		filePathsClean = append(filePathsClean, filePathClean)
	}

	slices.Sort(filePathsClean)
	filePathsUnique := slices.Compact(filePathsClean)

	var textFilePaths []string
	var nonTextFilePaths []string

	for _, filePath := range filePathsUnique {
		f, err := os.Open(filePath)
		if err != nil {
			log.Printf("can not open file: %s %v", filePath, err)
			continue
		}
		defer f.Close()
		first512bytes := make([]byte, 512)
		_, err = f.Read(first512bytes)
		if err != nil {
			log.Printf("can not read file: %s %v", filePath, err)
			continue
		}
		if !isProbablyText(first512bytes) {
			nonTextFilePaths = append(nonTextFilePaths, filePath)
			continue
		}
		textFilePaths = append(textFilePaths, filePath)
	}

	var filePaths []string
	if len(nonTextFilePaths) > 0 {
		filePaths = append(textFilePaths, nonTextFilePaths...)
	} else {
		filePaths = textFilePaths
	}

	if len(filePaths) == 0 {
		err := fmt.Errorf("no files picked")
		panicOnError(err)
	}

	fmt.Printf(`Your pick:
      text files: %d
  non-text files: %d
`, len(textFilePaths), len(nonTextFilePaths))

	var label string
	var prompt promptui.Select
	if len(nonTextFilePaths) > 0 {
		label = "Are you sure you want to add a BOM to ALL files anyway?"
		items := []string{NO, YES}
		if len(textFilePaths) > 0 {
			items = []string{NO, TEXT_FILES_ONLY, YES_ALL}
		}
		prompt = promptui.Select{
			Label: label,
			Items: items,
		}
	} else {
		label = "Are you sure you want to add a BOM to ALL files?"
		prompt = promptui.Select{
			Label: label,
			Items: []string{NO, YES},
		}
	}

	_, result, err := prompt.Run()
	panicOnError(err)

	if result == NO {
		os.Exit(0)
	}

	if result == TEXT_FILES_ONLY {
		filePaths = textFilePaths
	}

	for _, a := range filePaths {
		filePath := a
		err := addBom(filePath)
		if err != nil {
			log.Printf("skipped: %s %v", filePath, err)
			continue
		}
		fmt.Printf("BOM added: %s\n", filePath)
	}
}
