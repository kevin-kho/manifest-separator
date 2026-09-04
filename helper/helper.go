package helper

import (
	"bytes"
	"fmt"
	"log/slog"
	"manifest-seperator/models"
	"slices"

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
	if len(curr) == 0 {
		return res, nil
	}

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

// TODO: Refactor to GroupVersionKind
// Kind is not distinct enough
func GetKinds(manifests []models.Manifest) map[string]bool {
	fmt.Println("Getting set of Kinds")
	res := make(map[string]bool)
	for _, mani := range manifests {
		res[mani.Kind] = true
	}

	return res

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

func UnmarshalManifestBytesToMap(manifestBytes []models.ManifestByte) (map[models.Manifest]models.ManifestByte, error) {

	mp := make(map[models.Manifest]models.ManifestByte)
	for _, mb := range manifestBytes {
		mani, err := mb.UnmarshalManifest()
		if err != nil {
			return mp, err
		}
		mp[mani] = mb
	}
	return mp, nil

}
