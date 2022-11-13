package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type ConfigFile struct {
	RootDir            string `json:"rootDir"`
	StatementCharacter string `json:"statement_character"`
}

const configFileName = "config.json"
const outputFileName = "result.csv"

func getRootDirName(cfg *ConfigFile) (string, error) {
	file, err := os.Open(configFileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return "", err
	}

	return cfg.RootDir, nil
}

type Component struct {
	dirName    string
	percentage float32
	statements uint32
	files      uint32
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

func analyzeComponentsSize(rootDirName string, statementChar string) ([]Component, error) {
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
			if len(component.dirName) != 0 {
				results = append(results, component)
			}

			component = Component{}
			dirName := path
			component.dirName = dirName
			currentDir = dirName
		} else {
			if len(component.dirName) == 0 {
				return nil
			}
			component.files++
			statements, err := countOccurences(path, statementChar)
			if err != nil {
				fmt.Printf("error in countOccureneces: %v\n", err)
				return err
			}
			component.statements += statements
		}

		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", rootDirName, err)
		return []Component{}, err
	}
	return results, nil
}

func calculateTotalStatements(results *[]Component) uint32 {
	var totalStatements uint32 = 0

	for _, result := range *results {
		totalStatements += result.statements
	}

	return totalStatements
}

func calculatePercentage(results *[]Component) {
	totalStatements := calculateTotalStatements(results)

	for i := 0; i < len(*results); i++ {
		statements := (*results)[i].statements
		(*results)[i].percentage = float32(statements) / float32(totalStatements) * 100
	}
}

func outputResultsToFile(results *[]Component) error {
	file, err := os.Create(outputFileName) // 同階層にファイルを作成する
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("コンポーネント名, 名前空間, パーセント, 総ステートメント数, ファイル数\n")
	if err != nil {
		return err
	}

	calculatePercentage(results)
	for _, result := range *results {
		_, err := file.WriteString(fmt.Sprintf("%s, %.2f, %d, %d\n", result.dirName, result.percentage, result.statements, result.files))
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var cfg ConfigFile
	rootDir, err := getRootDirName(&cfg)
	if err != nil {
		fmt.Printf("error in getRootDirName: %v\n", err)
		return
	}
	if len(rootDir) == 0 {
		fmt.Printf("RootDir is empty string in config.json.")
		return
	}

	results, err := analyzeComponentsSize(rootDir, cfg.StatementCharacter)
	if err != nil {
		fmt.Printf("error in analyzeComponentsSize: %v\n", err)
		return
	}

	err = outputResultsToFile(&results)
	if err != nil {
		fmt.Printf("error in outputResultsToFile: %v\n", err)
	}
}
