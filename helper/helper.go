package helper

import (
	"bufio"
	"bytes"
	"fmt"
	"manifest-seperator/export"
	"manifest-seperator/models"
	"os"
	"slices"
)

func SeparateManifests(data []byte) ([]models.ManifestByte, error) {

	fmt.Println("Separating manifests...")

	var res []models.ManifestByte

	var curr models.ManifestByte

	for row := range bytes.SplitSeq(data, []byte{'\n'}) {
		if slices.Equal(row, []byte("---")) && len(curr) > 0 {
			valid, err := curr.IsValidManifest()
			if err != nil {
				return res, err
			}
			if valid {
				res = append(res, curr)
			}
			curr = []byte{}
			continue
		}
		curr = append(curr, row...)
		curr = append(curr, '\n')
	}

	curr = bytes.TrimSpace(curr)
	valid, err := curr.IsValidManifest()
	if err != nil {
		return res, err
	}
	if len(curr) > 0 && valid {
		res = append(res, curr)
	}

	return res, nil
}

func GetKinds(mb []models.ManifestByte) (map[string]bool, error) {

	fmt.Println("Getting set of Kinds")

	res := make(map[string]bool)
	for _, m := range mb {
		mani, err := m.UnmarshalManifest()
		if err != nil {
			return res, err
		}
		res[mani.Kind] = true
	}

	return res, nil

}

func ContainsDupes(manifests []models.Manifest) bool {

	seen := make(map[models.Manifest]bool)

	for _, m := range manifests {
		if seen[m] {
			return true
		}
		seen[m] = true
	}

	return false

}

func HandleDashes(data []byte) error {

	fmt.Println("Handling manifests separated by triple dashes")

	// Clear
	err := export.RemoveAllKindDir()
	if err != nil {
		return err
	}

	// Parse
	manifestBytes, err := SeparateManifests(data)
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
	if ContainsDupes(manifests) {
		return fmt.Errorf("Duplicate Manifest found")
	}

	kinds, err := GetKinds(manifestBytes)
	if err != nil {
		return err
	}

	// Write
	err = HandleWrite(kinds, manifestBytes)
	if err != nil {
		return err
	}

	return nil

}

func HandleList(data []byte) error {

	fmt.Println("Handling manifest as a Kind: List")

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
	if ContainsDupes(manifests) {
		return fmt.Errorf("Duplicate Manifest found")
	}

	kinds, err := GetKinds(manifestBytes)
	if err != nil {
		return err
	}

	// Write
	err = HandleWrite(kinds, manifestBytes)
	if err != nil {
		return err
	}

	return nil
}

func HandleWrite(kinds map[string]bool, manifestBytes []models.ManifestByte) error {
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

func ReadStdIn() ([]byte, error) {
	var res []byte

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Reading StdIn")
	fmt.Println("Press Ctrl-D to exit out.")

	for scanner.Scan() {
		res = append(res, scanner.Bytes()...)
		res = append(res, '\n')
	}

	if err := scanner.Err(); err != nil {
		return res, err
	}

	return res, nil
}
