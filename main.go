package main

import (
	"fmt"
	"log"
	"manifest-seperator/eval"
	"manifest-seperator/models"

	"github.com/spf13/pflag"
)

func main() {

	// Parse Flags
	modeFlag := pflag.StringP("mode", "m", "dash", "specify the method used to parse manifests. Possible values: dash, list, appset")
	fileFlag := pflag.StringP("file", "f", "", "path to manifest file or directory. All .yaml/.yml files in directory must be of the same type (dash, list, appset)")
	pflag.Parse()

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

	err = eval.EvaluateConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Separated Successfully!")
}
