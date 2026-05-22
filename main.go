package main

import (
	"fmt"
	"log"
	"manifest-seperator/export"
	"manifest-seperator/helper"
	"manifest-seperator/models"
	"os"
)

func main() {

	args := os.Args
	if len(args) < 2 {
		log.Fatal("No path to manifest given")
	}

	path := args[1]
	data, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	lb := models.CreateListByte(data)
	fmt.Println(lb)

	fmt.Println(string(lb))

	return // TODO: remove when feature complete

	manifestBytes := helper.SeparateManifests(data)
	kinds, err := helper.GetKinds(manifestBytes)
	if err != nil {
		log.Fatal(err)
	}

	// Clear out all existing manifests
	err = export.RemoveAllKindDir()
	if err != nil {
		log.Fatal(err)
	}

	err = export.CreateKindDir(kinds)
	if err != nil {
		log.Fatal(err)
	}

	var diffCmds []string
	var getCmds []string

	// TODO: split into two loops?
	for _, mb := range manifestBytes {

		err := export.WriteManifestToFile(mb)
		if err != nil {
			log.Fatal(err)
		}

		diffCmd, err := mb.GetCmd("diff")
		if err != nil {
			log.Fatal(err)
		}

		diffCmds = append(diffCmds, diffCmd)

		getCmd, err := mb.GetCmd("get")
		if err != nil {
			log.Fatal(err)
		}

		getCmds = append(getCmds, getCmd)

	}

	err = export.WriteCmdFile(diffCmds, "diff")
	if err != nil {
		log.Fatal(err)
	}

	err = export.WriteCmdFile(getCmds, "get")
	if err != nil {
		log.Fatal(err)
	}

}
