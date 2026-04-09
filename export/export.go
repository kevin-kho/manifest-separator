package export

import (
	"fmt"
	"manifest-seperator/models"
	"os"
)

func RemoveAllKindDir() error {

	err := os.RemoveAll("out/")
	if err != nil {
		return err
	}

	return nil
}

func CreateKindDir(kinds map[string]bool) error {
	for kind := range kinds {
		path := fmt.Sprintf("out/%v", kind)
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteManifestToFile(mb models.ManifestByte) error {

	m, err := mb.MarshlToManifest()
	if err != nil {
		return err
	}

	fileName := m.GetFileName()
	filePath := fmt.Sprintf("out/%v/%v", m.Kind, fileName)

	err = os.WriteFile(filePath, mb, 0644)
	if err != nil {
		return err
	}
	return nil

}

func GetKubectlDiffCmd(mb models.ManifestByte) (string, error) {

	var cmd string

	m, err := mb.MarshlToManifest()
	if err != nil {
		return cmd, err
	}

	fileName := m.GetFileName()
	filePath := fmt.Sprintf("out/%v/%v", m.Kind, fileName)

	cmd = fmt.Sprintf("kubectl diff -f %v", filePath)

	return cmd, nil

}
