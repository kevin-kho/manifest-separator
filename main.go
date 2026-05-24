package main

import (
	"flag"
	"log"
	"manifest-seperator/export"
	"manifest-seperator/helper"
	"manifest-seperator/models"
	"os"
)

func handleDashes(data []byte) error {

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

func main() {

	listFlag := flag.Bool("list", false, "parse the file as if it's Kind: List")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("No path to manifest given")
	}

	path := args[0]
	data, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	if *listFlag {
		err = handleList(data)
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	err = handleDashes(data)
	if err != nil {
		log.Fatal(err)
	}

}
