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
func (mb ManifestByte) IsValidManifest() (bool, error) {
	var empty Manifest

	mani, err := mb.UnmarshalManifest()
	if err != nil {
		return false, err
	}

	return mani != empty, nil
}

func (mb ManifestByte) UnmarshalManifest() (Manifest, error) {
	var m Manifest
	err := yaml.Unmarshal(mb, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

type Cmd int

const (
	CmdGet Cmd = iota
	CmdDiff
	CmdApply
)

func (mb ManifestByte) GetCmd(cmdType Cmd) (string, error) {

	var cmd string
	cmdString := map[Cmd]string{
		CmdGet:   "kubectl get -f %v -oyaml",
		CmdDiff:  "kubectl diff -f %v",
		CmdApply: "kubectl apply -f %v",
	}

	cmdStr, valid := cmdString[cmdType]
	if !valid {
		return cmd, fmt.Errorf("Unknown cmd: %v", cmdType)
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
