package main

import (
	"flag"
	"fmt"
	"log"
	"manifest-seperator/helper"
	"manifest-seperator/models"
)

func main() {

	// Parse Flags
	modeFlag := flag.String("mode", "dash", "specify the method used to parse manifests. Possible values: dash, list, appset")
	fileFlag := flag.String("f", "", "path to manifest file or directory")
	flag.Parse()

	var err error
	var config models.Config

	err = config.SetMode(modeFlag)
	if err != nil {
		log.Fatal(err)
	}

	err = config.HandleFileFlag(fileFlag)
	if err != nil {
		log.Fatal(err)
	}

	if len(config.FileMap) > 0 {
		err = config.CombineFiles()
	}
	if err != nil {
		log.Fatal(err)
	}

	if len(config.Data) == 0 {
		log.Fatalf("Length of data is 0")
	}

	switch config.Mode {
	case models.ModeDash:
		err = helper.HandleDashes(config)
	case models.ModeList:
		err = helper.HandleList(config)
	case models.ModeAppSet:
		err = helper.HandleAppSet(config)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Separated Successfully!")
}
