package run

import (
	"badrodem/localize"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/manifoldco/promptui"
	"github.com/spkg/bom"
)

var (
	bomBytes = []byte{0xEF, 0xBB, 0xBF}

	itemYes           string
	itemNo            string
	itemYesAll        string
	itemTextFilesOnly string
	itemRemoveBOM     string
)

func isProbablyText(b []byte) bool {
	s := strings.TrimRight(string(b), string(rune(0)))
	bytes := []byte(s)
	mimeType := http.DetectContentType(bytes)
	return strings.HasPrefix(mimeType, "text/")
}

func addBom(filePath string) error {
	b, err := os.ReadFile(filePath)
	if err != nil {
		err := fmt.Errorf("%s: %w", localize.T("can_not_open_file"), err)
		return err
	}

	if bytes.Equal(b[0:3], bomBytes) {
		err := fmt.Errorf(localize.T("already_has_a_bom"))
		return err
	}

	if !utf8.Valid(b) {
		err := fmt.Errorf(localize.T("not_valid_utf-8_encoded"))
		return err
	}

	f, err := os.Create(filePath)
	if err != nil {
		err := fmt.Errorf("%s: %w", localize.T("can_not_open_file"), err)
		return err
	}
	defer f.Close()

	// TODO: make these atomic
	_, err = f.Write(bomBytes)
	if err != nil {
		err := fmt.Errorf("%s: %w", localize.T("can_not_write_file"), err)
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		err := fmt.Errorf("%s: %w", localize.T("can_not_write_file"), err)
		return err
	}

	return nil
}

func removeBom(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		err := fmt.Errorf("%s: %w", localize.T("can_not_open_file"), err)
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	br := bom.NewReader(r)

	b, err := io.ReadAll(br)
	if err != nil {
		err := fmt.Errorf("%s: %w", localize.T("can_not_read_file"), err)
		return err
	}

	f, err = os.Create(filePath)
	if err != nil {
		err := fmt.Errorf("%s: %w", localize.T("can_not_open_file"), err)
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		err := fmt.Errorf("%s: %w", localize.T("can_not_write_file"), err)
		return err
	}

	return nil
}

func Run() error {
	if len(os.Args) <= 1 {
		err := fmt.Errorf(localize.T("one_or_more_arguments_are_required"))
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
			log.Printf("%s: %s %v", localize.T("can_not_open_file"), filePath, err)
			continue
		}
		defer f.Close()
		first512bytes := make([]byte, 512)
		_, err = f.Read(first512bytes)
		if err != nil {
			if err == io.EOF {
				log.Printf("%s: %s %v", localize.T("file_is_empty"), filePath, err)
			} else {
				log.Printf("%s: %s %v", localize.T("can_not_read_file"), filePath, err)
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
		err := fmt.Errorf(localize.T("no_files_picked"))
		return err
	}

	fmt.Printf(`%s:
%s: %d
%s: %d
`,
		localize.T("your_pick"),
		localize.T("text_files"), len(textFilePaths),
		localize.T("non_text_files"), len(nonTextFilePaths),
	)

	itemYes = localize.T("item_yes")
	itemNo = localize.T("item_no")
	itemYesAll = localize.T("item_yes_all")
	itemTextFilesOnly = localize.T("item_text_files_only")
	itemRemoveBOM = localize.T("item_remove_bom")

	var label string
	var prompt promptui.Select
	if len(nonTextFilePaths) > 0 {
		label = localize.T("are_you_sure_want_to_add_a_bom_to_all_files_anyway")
		items := []string{
			itemNo,
			itemYes,
		}
		if len(textFilePaths) > 0 {
			items = []string{
				itemNo,
				itemTextFilesOnly,
				itemYesAll,
				itemRemoveBOM,
			}
		}
		prompt = promptui.Select{
			Label: label,
			Items: items,
		}
	} else {
		label = localize.T("add_a_bom_to_all_files")
		prompt = promptui.Select{
			Label: label,
			Items: []string{
				itemYes,
				itemRemoveBOM,
				itemNo,
			},
		}
	}

	_, result, err := prompt.Run()
	if err != nil {
		return err
	}

	if result == itemNo {
		return nil
	}

	if result == itemTextFilesOnly {
		filePaths = textFilePaths
	}

	for _, a := range filePaths {
		filePath := a
		var err error
		if result == itemRemoveBOM {
			err = removeBom(filePath)
		} else {
			err = addBom(filePath)
		}
		if err != nil {
			log.Printf("%s: %v", localize.T("skipped"), err)
			continue
		}
		if result == itemRemoveBOM {
			fmt.Printf("%s: %s\n", localize.T("bom_removed"), filePath)
		} else {
			fmt.Printf("%s: %s\n", localize.T("bom_added"), filePath)
		}
	}

	return nil
}
