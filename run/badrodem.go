package run

import (
	"badrodem/platform"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/manifoldco/promptui"
	"github.com/spkg/bom"
)

const (
	YES             = "Yes"
	NO              = "No"
	YES_ALL         = "Yes, all"
	TEXT_FILES_ONLY = "Text files only"
	REMOVE_BOM      = "Remove BOM"
)

var (
	bomBytes = []byte{0xEF, 0xBB, 0xBF}
)

func exit(code int) {
	if runtime.GOOS == "windows" && platform.IsDoubleClickRun() {
		// keep console open on exit
		fmt.Print("Press any key to continue . . .")
		os.Stdin.Read(make([]byte, 1))
	}
	os.Exit(code)
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

	if bytes.Equal(b[0:3], bomBytes) {
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

	_, err = f.Write(bomBytes)
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

func removeBom(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		err := fmt.Errorf("can not open file: %w", err)
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	br := bom.NewReader(r)

	b, err := io.ReadAll(br)
	if err != nil {
		err := fmt.Errorf("can not read file: %w", err)
		return err
	}

	f, err = os.Create(filePath)
	if err != nil {
		err := fmt.Errorf("can not open file: %w", err)
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		err := fmt.Errorf("can not write file: %w", err)
		return err
	}

	return nil
}

func Run() error {
	if len(os.Args) <= 1 {
		err := fmt.Errorf("one or more arguments are required")
		return err
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
			if err == io.EOF {
				log.Printf("empty file: %s %v", filePath, err)
			} else {
				log.Printf("can not read file: %s %v", filePath, err)
			}
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
		return err
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
			items = []string{NO, TEXT_FILES_ONLY, YES_ALL, REMOVE_BOM}
		}
		prompt = promptui.Select{
			Label: label,
			Items: items,
		}
	} else {
		label = "Add a BOM to ALL files?"
		prompt = promptui.Select{
			Label: label,
			Items: []string{YES, REMOVE_BOM, NO},
		}
	}

	_, result, err := prompt.Run()
	if err != nil {
		return err
	}

	if result == NO {
		return nil
	}

	if result == TEXT_FILES_ONLY {
		filePaths = textFilePaths
	}

	for _, a := range filePaths {
		filePath := a
		var err error
		if result == REMOVE_BOM {
			err = removeBom(filePath)
		} else {
			err = addBom(filePath)
		}
		if err != nil {
			log.Printf("skipped: %s %v", filePath, err)
			continue
		}
		if result == REMOVE_BOM {
			fmt.Printf("BOM removed: %s\n", filePath)
		} else {
			fmt.Printf("BOM added: %s\n", filePath)
		}
	}

	return nil
}
