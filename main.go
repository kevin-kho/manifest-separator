package main

import (
	"log"
	"manifest-seperator/export"
	"manifest-seperator/helper"
	"manifest-seperator/models"
	"os"
)

func handleDashes(data []byte) error {
	manifestBytes := helper.SeparateManifests(data)
	kinds, err := helper.GetKinds(manifestBytes)
	if err != nil {
		return err
	}

	// Clear out all existing manifests
	err = export.RemoveAllKindDir()
	if err != nil {
		return err
	}

	err = export.CreateKindDir(kinds)
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

func handleList(data []byte) error {
	lb := models.CreateListByte(data)

	lst, err := lb.UnmarshalManifest()
	if err != nil {
		return err
	}

	itemsDetailed, err := lst.GetItemsDetailed()
	if err != nil {
		return err
	}

	var manifestBytes []models.ManifestByte
	for _, i := range itemsDetailed {
		manifestBytes = append(manifestBytes, i.ManifestByte)
	}

	kinds, err := helper.GetKinds(manifestBytes)
	if err != nil {
		return err
	}
	// Clear out all existing manifests
	err = export.RemoveAllKindDir()
	if err != nil {
		return err
	}

	err = export.CreateKindDir(kinds)
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

	args := os.Args
	if len(args) < 2 {
		log.Fatal("No path to manifest given")
	}

	path := args[1]
	data, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	err = handleList(data)
	if err != nil {
		log.Fatal(err)
	}

	return // TODO: remove when feature complete

	err = handleDashes(data)
	if err != nil {
		log.Fatal(err)
	}

}
