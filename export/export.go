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

func WriteCmdFile(cmds []string, cmdType models.Cmd) error {
	concat := strings.Join(cmds, "\n")
	filePaths := map[models.Cmd]string{
		models.CmdDiff:  "out/diff-cmds.txt",
		models.CmdGet:   "out/get-cmds.txt",
		models.CmdApply: "out/apply-cmds.txt",
	}
	filePath, valid := filePaths[cmdType]
	if !valid {
		return fmt.Errorf("Unknown cmdType: %v", cmdType)
	}

	fmt.Printf("Writing file: %v\n", filePath)
	err := os.WriteFile(filePath, []byte(concat), 0644)
	if err != nil {
		return err
	}

	return nil
}

func HandleWrite(kinds map[string]bool, manifestBytes []models.ManifestByte) error {
	err := CreateKindDir(kinds)
	if err != nil {
		return err
	}

	var diffCmds []string
	var getCmds []string
	var applyCmds []string

	// TODO: split into two loops?
	for _, mb := range manifestBytes {

		err := WriteManifestToFile(mb)
		if err != nil {
			return err
		}

		diffCmd, err := mb.GetCmd(models.CmdDiff)
		if err != nil {
			return err
		}

		diffCmds = append(diffCmds, diffCmd)

		getCmd, err := mb.GetCmd(models.CmdGet)
		if err != nil {
			return err
		}

		getCmds = append(getCmds, getCmd)

		applyCmd, err := mb.GetCmd(models.CmdApply)
		if err != nil {
			return err
		}

		applyCmds = append(applyCmds, applyCmd)

	}

	err = WriteCmdFile(diffCmds, models.CmdDiff)
	if err != nil {
		return err
	}

	err = WriteCmdFile(getCmds, models.CmdGet)
	if err != nil {
		return err
	}

	err = WriteCmdFile(applyCmds, models.CmdApply)
	if err != nil {
		return err
	}

	return nil
}
