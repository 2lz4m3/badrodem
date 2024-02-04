package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
)

const (
	YES = "Yes"
	NO  = "No"
)

var (
	bom = []byte{0xEF, 0xBB, 0xBF}
)

func panicOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
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

	label := fmt.Sprintf("Are you sure you want to add a BOM to %d file(s)? [Yes/No]", len(os.Args[1:]))

	prompt := promptui.Select{
		Label: label,
		Items: []string{YES, NO},
	}
	_, result, err := prompt.Run()
	panicOnError(err)

	if result != YES {
		os.Exit(0)
	}

	for i, a := range os.Args[1:] {
		filePath := a
		err := addBom(filePath)
		if err != nil {
			log.Printf("error occurred at #%d: %v\n", i, err)
		}
	}
}
