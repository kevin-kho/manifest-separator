package main

import (
	"flag"
	"fmt"
	"log"
	"manifest-seperator/helper"
	"os"
)

func main() {

	appSetFlag := flag.Bool("appset", false, "parse the manifest of a generated ArgoCD AppSet")
	listFlag := flag.Bool("list", false, "parse the manifest as if it's Kind: List. Same as -l flag")
	lFlag := flag.Bool("l", false, "parse the manifest as if it's Kind: List. Same as -list flag")
	fileFlag := flag.String("f", "", "path to manifest file")
	flag.Parse()

	var data []byte
	var files [][]byte
	var err error
	var fileInfo os.FileInfo

	if *fileFlag != "" {
		fileInfo, err = os.Stat(*fileFlag)
		if err != nil {
			log.Fatal(err)
		}
		if fileInfo.IsDir() {
			files, err = helper.ReadDir(*fileFlag)
		} else {
			fmt.Printf("Reading file: %v\n", *fileFlag)
			data, err = os.ReadFile(*fileFlag)
		}
	} else {
		data, err = helper.ReadStdIn()
	}

	if err != nil {
		log.Fatal(err)
	}

	if len(files) > 0 {
		var strategy string
		if *listFlag || *lFlag {
			strategy = "list"
		} else if *appSetFlag {
			strategy = "appset"
		} else {
			strategy = "tripleDash"
		}
		data, err = helper.CombineFiles(files, strategy)
	}

	if err != nil {
		log.Fatal(err)
	}

	if len(data) == 0 {
		log.Fatalf("Length of data is 0")
	}

	// TODO: refactor to use a single "mode" flag to determine how to process files
	// In future, have an "auto" mode which determines how the manifest should be handled
	if *listFlag || *lFlag {
		err = helper.HandleList(data)
		if err != nil {
			log.Fatal(err)
		}

	} else if *appSetFlag {
		err = helper.HandleAppSet(data)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = helper.HandleDashes(data)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Separated Successfully!")
}
