package analyzer

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Component struct {
	DirName    string
	Percentage float32
	Statements uint32
	Files      uint32
}

func AnalyzeComponentsSize(rootDirName string, statementChar string) ([]Component, error) {
	var results []Component

	currentDir := ""
	component := Component{}
	err := filepath.Walk(rootDirName, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if path == rootDirName {
			return nil
		}

		if info.IsDir() {
			if len(component.DirName) != 0 {
				results = append(results, component)
			}

			component = Component{}
			dirName := path
			component.DirName = dirName
			currentDir = dirName
		} else {
			if len(component.DirName) == 0 {
				return nil
			}
			component.Files++
			statements, err := countOccurences(path, statementChar)
			if err != nil {
				fmt.Printf("error in countOccureneces: %v\n", err)
				return err
			}
			component.Statements += statements
		}

		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", rootDirName, err)
		return []Component{}, err
	}
	return results, nil
}

func countOccurences(pathName string, targetChar string) (uint32, error) {
	file, err := os.Open(pathName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	buf := make([]byte, 1024)
	var occurences uint32 = 0
	for {
		_, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if err != nil && err == io.EOF {
			break
		}
		occurences += uint32(strings.Count(string(buf), targetChar))
	}
	return occurences, nil
}
