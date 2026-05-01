package export

import (
	"fmt"
	"manifest-seperator/models"
	"os"
	"strings"
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

// TODO: make into receiver function?
func WriteManifestToFile(mb models.ManifestByte) error {

	m, err := mb.UnmarshalManifest()
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

// TODO: make into receiver function?
func GetCmd(mb models.ManifestByte, cmdType string) (string, error) {
	// TODO: make cmdType an enum?

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
	filePath := fmt.Sprintf(cmdStr, m.Kind, fileName)

	cmd = fmt.Sprintf(cmdStr, filePath)

	return cmd, nil

}

// TODO: make into receiver function?
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

	err := os.WriteFile(filePath, []byte(concat), 0644)
	if err != nil {
		return err
	}

	return nil
}
