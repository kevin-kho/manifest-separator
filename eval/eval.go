package eval

import (
	"fmt"
	"manifest-seperator/export"
	"manifest-seperator/helper"
	"manifest-seperator/models"
	"maps"
	"slices"
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

	// Clear
	err := export.RemoveAllKindDir()
	if err != nil {
		return err
	}

	// Parse
	manifestBytes, err := helper.SeparateManifests(config.Data)
	if err != nil {
		return err
	}

	// Check
	// manifests, err := helper.UnmarshalManifestBytes(manifestBytes)
	mp, err := helper.UnmarshalManifestBytesToMap(manifestBytes)
	if err != nil {
		return err
	}
	manifests := slices.Collect(maps.Keys(mp))
	if helper.ContainsDupes(manifests) {
		return fmt.Errorf("Duplicate Manifest found")
	}

	kinds := helper.GetKinds(manifests)

	// Write
	err = export.HandleWrite(kinds, manifestBytes)
	if err != nil {
		return err
	}

	return nil

}

func HandleList(config models.Config) error {

	fmt.Println("Handling manifest as a Kind: List")

	// Clear
	err := export.RemoveAllKindDir()
	if err != nil {
		return err
	}

	// Parse
	var lb models.ListByte = config.Data
	lst, err := lb.UnmarshalManifest()
	if err != nil {
		return err
	}

	manifestBytes, err := lst.GetManifestBytes()
	if err != nil {
		return err
	}

	// Check
	mp, err := helper.UnmarshalManifestBytesToMap(manifestBytes)
	if err != nil {
		return err
	}
	manifests := slices.Collect(maps.Keys(mp))
	if helper.ContainsDupes(manifests) {
		return fmt.Errorf("Duplicate Manifest found")
	}

	kinds := helper.GetKinds(manifests)

	// Write
	err = export.HandleWrite(kinds, manifestBytes)
	if err != nil {
		return err
	}

	return nil
}

func HandleAppSet(config models.Config) error {
	fmt.Println("Handling ArgoCD AppSet")

	// Clear
	err := export.RemoveAllKindDir()
	if err != nil {
		return err
	}

	// Parse
	// Generated AppSet comes as array of App
	manifestBytes, err := helper.SeparateAppSet(config.Data)
	if err != nil {
		return err
	}

	// Check
	mp, err := helper.UnmarshalManifestBytesToMap(manifestBytes)
	if err != nil {
		return err
	}
	manifests := slices.Collect(maps.Keys(mp))
	if helper.ContainsDupes(manifests) {
		return fmt.Errorf("Duplicate Manifest found")
	}

	kinds := helper.GetKinds(manifests)

	// Write
	err = export.HandleWrite(kinds, manifestBytes)
	if err != nil {
		return err
	}

	return nil
}
