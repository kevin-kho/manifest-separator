package models

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

type Manifest struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
}

func (m Manifest) GetFileName() string {
	if m.Metadata.Namespace == "" {
		return fmt.Sprintf("%v_%v.yaml", m.Kind, m.Metadata.Name)
	}
	return fmt.Sprintf("%v_%v_%v.yaml", m.Kind, m.Metadata.Name, m.Metadata.Namespace)
}

type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type ManifestByte []byte

// Covers the scenario where there's nothing inbetween --- ex:
// ---
// # Source: cni/templates/clusterrolebinding.yaml
// ---
func (mb ManifestByte) IsValidManifest() bool {
	var empty Manifest

	mani, err := mb.MarshalToManifest()
	if err != nil {
		return false // TODO: should it return error?
	}

	return mani != empty
}

func (mb ManifestByte) MarshalToManifest() (Manifest, error) {
	var m Manifest
	err := yaml.Unmarshal(mb, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}
