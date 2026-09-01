package eval

import (
	"fmt"
	"manifest-seperator/export"
	"manifest-seperator/helper"
	"manifest-seperator/models"
)

func EvaluateConfig(config models.Config) error {
	var err error
	switch config.Mode {
	case models.ModeDash:
		err = HandleDashes(config)
	case models.ModeList:
		err = HandleList(config)
	case models.ModeAppSet:
		err = HandleAppSet(config)
	}
	return err
}

func HandleDashes(config models.Config) error {

	fmt.Println("Handling manifests separated by triple dashes")
	data := config.Data

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

	// Check
	var manifests []models.Manifest
	for _, mb := range manifestBytes {
		mani, err := mb.UnmarshalManifest()
		if err != nil {
			return err
		}
		manifests = append(manifests, mani)
	}
	if helper.ContainsDupes(manifests) {
		return fmt.Errorf("Duplicate Manifest found")
	}

	kinds, err := helper.GetKinds(manifestBytes)
	if err != nil {
		return err
	}

	// Write
	err = export.HandleWrite(kinds, manifestBytes)
	if err != nil {
		return err
	}

	return nil

}

func HandleList(config models.Config) error {

	fmt.Println("Handling manifest as a Kind: List")
	data := config.Data

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

	// Check
	var manifests []models.Manifest
	for _, mb := range manifestBytes {
		mani, err := mb.UnmarshalManifest()
		if err != nil {
			return err
		}
		manifests = append(manifests, mani)
	}
	if helper.ContainsDupes(manifests) {
		return fmt.Errorf("Duplicate Manifest found")
	}

	kinds, err := helper.GetKinds(manifestBytes)
	if err != nil {
		return err
	}

	// Write
	err = export.HandleWrite(kinds, manifestBytes)
	if err != nil {
		return err
	}

	return nil
}

func HandleAppSet(config models.Config) error {
	data := config.Data
	fmt.Println("Handling ArgoCD AppSet")

	// Clear
	err := export.RemoveAllKindDir()
	if err != nil {
		return err
	}

	// Parse
	// Generated AppSet comes as array of App
	manifestBytes, err := helper.SeparateAppSet(data)
	if err != nil {
		return err
	}

	// Check
	var manifests []models.Manifest
	for _, mb := range manifestBytes {
		mani, err := mb.UnmarshalManifest()
		if err != nil {
			return err
		}
		manifests = append(manifests, mani)
	}
	if helper.ContainsDupes(manifests) {
		return fmt.Errorf("Duplicate Manifest found")
	}

	kinds, err := helper.GetKinds(manifestBytes)
	if err != nil {
		return err
	}

	// Write
	err = export.HandleWrite(kinds, manifestBytes)
	if err != nil {
		return err
	}

	return nil
}
