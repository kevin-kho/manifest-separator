package models

import (
	"errors"
	"fmt"
	"manifest-seperator/filereader"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	FileMap map[string][]byte
	Data    []byte
	Mode    Mode
}

func (c *Config) CombineFiles() error {
	if c.FileMap == nil {
		return errors.New("FileMap is nil")
	}
	if len(c.FileMap) == 0 {
		return errors.New("FileMap is empty")
	}
	if len(c.Data) > 0 {
		return errors.New("Cannot CombineFiles, config.Data is non-empty")
	}
	var res []byte
	var err error

	switch c.Mode {
	case ModeList:
		fmt.Println("CombineFiles ModeList")
		var items []any
		for name, data := range c.FileMap {
			fmt.Println("Combining " + name)
			var lb ListByte = data
			lst, err := lb.UnmarshalManifest()
			if err != nil {
				return err
			}
			items = append(items, lst.Items...)
		}

		l := List{
			ApiVersion: "v1",
			Kind:       "List",
			Metadata:   Metadata{},
			Items:      items,
		}

		res, err = yaml.Marshal(l)

	case ModeAppSet:
		fmt.Println("CombineFiles ModeAppSet")
		for name, data := range c.FileMap {
			fmt.Println("Combining " + name)
			res = append(res, data...)
			res = append(res, []byte("\n")...)
		}
	case ModeDash:
		fmt.Println("CombineFiles ModeDash")
		for name, data := range c.FileMap {
			fmt.Println("Combining " + name)
			res = append(res, data...)
			res = append(res, []byte("\n---\n")...)
		}
	}

	if err != nil {
		return err
	}

	c.Data = res

	return nil

}

func (c *Config) SetMode(modeFlag *string) error {
	switch *modeFlag {
	case "dash":
		c.Mode = ModeDash
	case "list":
		c.Mode = ModeList
	case "appset":
		c.Mode = ModeAppSet
	default:
		return errors.New("unknown mode flag. Must be one of the following: dash, list, appset")
	}
	return nil
}

func (c *Config) HandleFileFlag(fileFlag *string) error {

	var err error
	if *fileFlag == "" {
		c.Data, err = filereader.ReadStdIn()
		if err != nil {
			return err
		}
		return nil

	}

	fileInfo, err := os.Stat(*fileFlag)
	if err != nil {
		return err
	}
	switch {
	case fileInfo.Mode().IsDir():
		fmt.Printf("Reading dir: %v\n", *fileFlag)
		c.FileMap, err = filereader.ReadDir(*fileFlag)
	case fileInfo.Mode().IsRegular():
		fmt.Printf("Reading file: %v\n", *fileFlag)
		c.Data, err = os.ReadFile(*fileFlag)
	}
	if err != nil {
		return err
	}
	return nil

}
