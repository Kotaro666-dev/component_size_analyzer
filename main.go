package main

import (
	"encoding/json"
	"fmt"
	Analyzer "github.com/Kotaro666-dev/component_size_analyzer/analyzer"
	Output "github.com/Kotaro666-dev/component_size_analyzer/output"
	"os"
)

type ConfigFile struct {
	RootDir            string `json:"rootDir"`
	StatementCharacter string `json:"statement_character"`
}

const configFileName = "config.json"

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

	results, err := Analyzer.AnalyzeComponentsSize(rootDir, cfg.StatementCharacter)
	if err != nil {
		fmt.Printf("error in analyzeComponentsSize: %v\n", err)
		return
	}

	err = Output.ResultsToFile(&results)
	if err != nil {
		fmt.Printf("error in outputResultsToFile: %v\n", err)
	}
}
