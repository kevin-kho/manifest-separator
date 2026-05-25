package export

import (
	"fmt"
	"manifest-seperator/models"
	"os"
	"strings"
)

func RemoveAllKindDir() error {

	fmt.Println("Clearing /out directory")
	err := os.RemoveAll("out/")
	if err != nil {
		return err
	}

	return nil
}

func CreateKindDir(kinds map[string]bool) error {
	for kind := range kinds {
		path := fmt.Sprintf("out/%v", kind)
		fmt.Printf("Creating directory: %v\n", path)
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteManifestToFile(mb models.ManifestByte) error {

	m, err := mb.UnmarshalManifest()
	if err != nil {
		return err
	}

	fileName := m.GetFileName()
	filePath := fmt.Sprintf("out/%v/%v", m.Kind, fileName)

	fmt.Printf("Writing file: %v\n", filePath)
	err = os.WriteFile(filePath, mb, 0644)
	if err != nil {
		return err
	}
	return nil

}

func WriteCmdFile(cmds []string, cmdType string) error {
	concat := strings.Join(cmds, "\n")
	filePaths := map[string]string{
		"diff": "out/diff-cmds.txt",
		"get":  "out/get-cmds.txt",
	}
	filePath := filePaths[cmdType]
	if filePath == "" {
		return fmt.Errorf("Unknown cmdType: %v", cmdType)
	}

	fmt.Printf("Writing kubectl %v file\n", cmdType)
	err := os.WriteFile(filePath, []byte(concat), 0644)
	if err != nil {
		return err
	}

	return nil
}
