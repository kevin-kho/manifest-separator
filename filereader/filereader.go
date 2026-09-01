package filereader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadStdIn() ([]byte, error) {
	var res []byte

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Reading StdIn")
	fmt.Println("Press Ctrl-D to exit out.")

	for scanner.Scan() {
		res = append(res, scanner.Bytes()...)
		res = append(res, '\n')
	}

	if err := scanner.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func ReadDir(path string) (map[string][]byte, error) {
	res := make(map[string][]byte)
	entries, err := os.ReadDir(path)
	if err != nil {
		return res, err
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".yaml") || strings.HasSuffix(entry.Name(), ".yml") {

			data, err := os.ReadFile(fmt.Sprintf("%v/%v", path, entry.Name()))
			if err != nil {
				return res, err
			}
			res[entry.Name()] = data
		}
	}

	return res, nil

}
