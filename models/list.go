package models

import "github.com/goccy/go-yaml"

type List struct {
	Manifest
	Items []any `yaml:"items"`
}

func (l List) GetManifestBytes() ([]ManifestByte, error) {

	var res []ManifestByte

	for _, item := range l.Items {
		var b ManifestByte

		b, err := yaml.Marshal(item)
		if err != nil {
			return res, err
		}

		valid, err := b.IsValidManifest()
		if err != nil {
			return res, err
		}

		if valid {
			res = append(res, b)
		}

	}

	return res, nil

}

type ListByte []byte

func (lb ListByte) UnmarshalManifest() (List, error) {
	var m List
	err := yaml.Unmarshal(lb, &m)
	if err != nil {
		return m, err
	}
	return m, nil

}
