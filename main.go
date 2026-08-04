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
	var err error
	if *fileFlag != "" {
		fmt.Printf("Reading file: %v\n", *fileFlag)
		data, err = os.ReadFile(*fileFlag)
	} else {
		data, err = helper.ReadStdIn()
	}

	if err != nil {
		log.Fatal(err)
	}

	if len(data) == 0 {
		log.Fatalf("Length of data is 0")
	}

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
