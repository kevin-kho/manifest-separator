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

		if b.IsValidManifest() {
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

	mani, err := mb.UnmarshalManifest()
	if err != nil {
		return false // TODO: should it return error?
	}

	return mani != empty
}

func (mb ManifestByte) UnmarshalManifest() (Manifest, error) {
	var m Manifest
	err := yaml.Unmarshal(mb, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

func (mb ManifestByte) GetCmd(cmdType string) (string, error) {

	var cmd string
	cmdString := map[string]string{
		"get":  "kubectl get -f %v -oyaml",
		"diff": "kubectl diff -f %v",
	}

	cmdStr := cmdString[cmdType]
	if cmdStr == "" {
		return cmd, fmt.Errorf("Unknown cmdType: %v", cmdType)
	}

	m, err := mb.UnmarshalManifest()
	if err != nil {
		return cmd, err
	}

	fileName := m.GetFileName()
	filePath := fmt.Sprintf("out/%v/%v", m.Kind, fileName)

	cmd = fmt.Sprintf(cmdStr, filePath)

	return cmd, nil
}
