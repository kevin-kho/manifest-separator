package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"manifest-seperator/helper"
	"manifest-seperator/models"
	"os"
)

func main() {

	// Parse Flags
	modeFlag := flag.String("mode", "dash", "specify the method used to parse manifests. Possible values: dash, list, appset")
	fileFlag := flag.String("f", "", "path to manifest file or directory")
	flag.Parse()

	var data []byte
	var fileMap map[string][]byte
	var err error
	var fileInfo os.FileInfo
	var mode models.Mode

	switch *modeFlag {
	case "dash":
		mode = models.ModeDash
	case "list":
		mode = models.ModeList
	case "appset":
		mode = models.ModeAppSet
	default:
		err = errors.New("unknown mode flag. Must be one of the following: dash, list, appset")
	}
	if err != nil {
		log.Fatal(err)
	}

	if *fileFlag != "" {
		fileInfo, err = os.Stat(*fileFlag)
		if err != nil {
			log.Fatal(err)
		}

		switch {
		case fileInfo.Mode().IsDir():
			fmt.Printf("Reading dir: %v\n", *fileFlag)
			fileMap, err = helper.ReadDir(*fileFlag)
		case fileInfo.Mode().IsRegular():
			fmt.Printf("Reading file: %v\n", *fileFlag)
			data, err = os.ReadFile(*fileFlag)
		}

	} else {
		data, err = helper.ReadStdIn()
	}

	if err != nil {
		log.Fatal(err)
	}

	if len(fileMap) > 0 {
		data, err = helper.CombineFiles(fileMap, mode)
	}

	if err != nil {
		log.Fatal(err)
	}

	if len(data) == 0 {
		log.Fatalf("Length of data is 0")
	}

	switch mode {
	case models.ModeDash:
		err = helper.HandleDashes(data)
	case models.ModeList:
		err = helper.HandleList(data)
	case models.ModeAppSet:
		err = helper.HandleAppSet(data)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Separated Successfully!")
}
