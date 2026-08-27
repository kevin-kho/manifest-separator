package helper

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"manifest-seperator/export"
	"manifest-seperator/models"
	"os"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
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
			} else {
				slog.Warn("Skipping invalid manifest")
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

	if !valid {
		slog.Warn("Skipping invalid manifest")
		return res, nil
	}

	if len(curr) > 0 {
		res = append(res, curr)
	}

	return res, nil
}

func SeparateAppSet(data []byte) ([]models.ManifestByte, error) {

	var appSet models.AppSet
	var res []models.ManifestByte

	err := yaml.Unmarshal(data, &appSet)
	if err != nil {
		return res, err
	}

	for _, app := range appSet {
		b, err := yaml.Marshal(app)
		if err != nil {
			return res, err
		}
		res = append(res, b)
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

func HandleAppSet(data []byte) error {
	fmt.Println("Handling ArgoCD AppSet")

	// Clear
	err := export.RemoveAllKindDir()
	if err != nil {
		return err
	}

	// Parse
	// Generated AppSet comes as array of App
	manifestBytes, err := SeparateAppSet(data)
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
	var applyCmds []string

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

		applyCmd, err := mb.GetCmd("apply")
		if err != nil {
			return err
		}

		applyCmds = append(applyCmds, applyCmd)

	}

	err = export.WriteCmdFile(diffCmds, "diff")
	if err != nil {
		return err
	}

	err = export.WriteCmdFile(getCmds, "get")
	if err != nil {
		return err
	}

	err = export.WriteCmdFile(applyCmds, "apply")
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

func ReadDir(path string) ([][]byte, error) {
	var res [][]byte
	entries, err := os.ReadDir(path)
	if err != nil {
		return res, err
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".yaml") || strings.HasSuffix(entry.Name(), ".yml") {

			data, err := os.ReadFile(fmt.Sprintf("%v/%v", path, entry.Name()))
			if err != nil {
				return res, err
			}
			res = append(res, data)
		}
	}

	return res, nil

}

func CombineFiles(files [][]byte, mode models.Mode) ([]byte, error) {
	var res []byte
	var err error

	switch mode {
	case models.ModeList:
		fmt.Println("CombineFiles ModeList")
		var items []any
		for _, f := range files {
			var lb models.ListByte = f
			lst, err := lb.UnmarshalManifest()
			if err != nil {
				return res, err
			}
			items = append(items, lst.Items...)
		}

		l := models.List{
			ApiVersion: "v1",
			Kind:       "List",
			Metadata:   models.Metadata{},
			Items:      items,
		}

		res, err = yaml.Marshal(l)

	case models.ModeAppSet:
		fmt.Println("CombineFiles ModeAppSet")
		for _, f := range files {
			res = append(res, f...)
			res = append(res, []byte("\n")...)
		}
	case models.ModeDash:
		fmt.Println("CombineFiles ModeDash")
		for _, f := range files {
			res = append(res, f...)
			res = append(res, []byte("\n---\n")...)
		}
	}

	if err != nil {
		return res, err
	}

	return res, nil

}
