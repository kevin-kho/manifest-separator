package helper

import (
	"bytes"
	"fmt"
	"manifest-seperator/models"
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
