package main

import (
	"log"
	"manifest-seperator/export"
	"manifest-seperator/helper"
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

		diffCmd, err := export.GetCmd(mb, "diff")
		if err != nil {
			log.Fatal(err)
		}

		diffCmds = append(diffCmds, diffCmd)

		getCmd, err := export.GetCmd(mb, "get")
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
