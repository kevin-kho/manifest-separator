package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"manifest-seperator/export"
	"manifest-seperator/helper"
	"manifest-seperator/models"
	"os"
)

func handleDashes(data []byte) error {

	fmt.Println("Handling manifests separated by triple dashes")

	// Clear
	err := export.RemoveAllKindDir()
	if err != nil {
		return err
	}

	// Parse
	manifestBytes, err := helper.SeparateManifests(data)
	if err != nil {
		return err
	}

	kinds, err := helper.GetKinds(manifestBytes)
	if err != nil {
		return err
	}

	// Write
	err = handleWrite(kinds, manifestBytes)
	if err != nil {
		return err
	}

	return nil

}

func handleList(data []byte) error {

	fmt.Println("Handling manifest as a Kind: List")

	// Clear
	err := export.RemoveAllKindDir()
	if err != nil {
		return err
	}

	// Parse
	var lb models.ListByte = data
	lst, err := lb.UnmarshalManifest()
	if err != nil {
		return err
	}

	manifestBytes, err := lst.GetManifestBytes()
	if err != nil {
		return err
	}

	kinds, err := helper.GetKinds(manifestBytes)
	if err != nil {
		return err
	}

	// Write
	err = handleWrite(kinds, manifestBytes)
	if err != nil {
		return err
	}

	return nil
}

func handleWrite(kinds map[string]bool, manifestBytes []models.ManifestByte) error {
	err := export.CreateKindDir(kinds)
	if err != nil {
		return err
	}

	var diffCmds []string
	var getCmds []string

	// TODO: split into two loops?
	for _, mb := range manifestBytes {

		err := export.WriteManifestToFile(mb)
		if err != nil {
			return err
		}

		diffCmd, err := mb.GetCmd("diff")
		if err != nil {
			return err
		}

		diffCmds = append(diffCmds, diffCmd)

		getCmd, err := mb.GetCmd("get")
		if err != nil {
			return err
		}

		getCmds = append(getCmds, getCmd)

	}

	err = export.WriteCmdFile(diffCmds, "diff")
	if err != nil {
		return err
	}

	err = export.WriteCmdFile(getCmds, "get")
	if err != nil {
		return err
	}

	return nil
}

func readStdIn() ([]byte, error) {
	var res []byte

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Reading StdIn")
	fmt.Println("Press Ctrl-D to exit out.")

	for scanner.Scan() {
		res = append(res, scanner.Bytes()...)
		res = append(res, '\n')
	}

	if err := scanner.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func main() {

	listFlag := flag.Bool("list", false, "parse the manifest as if it's Kind: List. Same as -l flag")
	lFlag := flag.Bool("l", false, "parse the manifest as if it's Kind: List. Same as -list flag")
	fileFlag := flag.String("f", "", "path to manifest file")
	flag.Parse()

	var data []byte
	var err error
	if *fileFlag != "" {
		fmt.Printf("Reading file: %v\n", *fileFlag)
		data, err = os.ReadFile(*fileFlag)
	} else {
		data, err = readStdIn()
	}

	if err != nil {
		log.Fatal(err)
	}

	if len(data) == 0 {
		log.Fatalf("Length of data is 0")
	}

	if *listFlag || *lFlag {
		err = handleList(data)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		err = handleDashes(data)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Separated Successfully!")
}
