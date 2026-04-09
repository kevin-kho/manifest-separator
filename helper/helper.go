package helper

import (
	"bytes"
	"manifest-seperator/models"
	"slices"
)

func SeparateManifests(data []byte) []models.ManifestByte {

	var res []models.ManifestByte

	var curr []byte

	for row := range bytes.SplitSeq(data, []byte{'\n'}) {
		if slices.Equal(row, []byte("---")) && len(curr) > 0 {
			res = append(res, curr)
			curr = []byte{}
			continue
		}
		curr = append(curr, row...)
		curr = append(curr, '\n')
	}

	curr = bytes.TrimSpace(curr)
	if len(curr) > 0 {
		res = append(res, curr)
	}

	return res
}

func GetKinds(mb []models.ManifestByte) (map[string]bool, error) {

	res := make(map[string]bool)
	for _, m := range mb {
		mani, err := m.MarshlToManifest()
		if err != nil {
			return res, err
		}
		res[mani.Kind] = true
	}

	return res, nil

}
