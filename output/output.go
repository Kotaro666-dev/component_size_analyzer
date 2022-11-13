package output

import (
	"fmt"
	Analyzer "github.com/Kotaro666-dev/component_size_analyzer/analyzer"
	"os"
)

const outputFileName = "result.csv"

func calculateTotalStatements(results *[]Analyzer.Component) uint32 {
	var totalStatements uint32 = 0

	for _, result := range *results {
		totalStatements += result.Statements
	}

	return totalStatements
}

func calculatePercentage(results *[]Analyzer.Component) {
	totalStatements := calculateTotalStatements(results)

	for i := 0; i < len(*results); i++ {
		statements := (*results)[i].Statements
		(*results)[i].Percentage = float32(statements) / float32(totalStatements) * 100
	}
}

func ResultsToFile(results *[]Analyzer.Component) error {
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
		_, err := file.WriteString(fmt.Sprintf("%s, %s, %.2f, %d, %d\n", result.ComponentName, result.NameSpace, result.Percentage, result.Statements, result.Files))
		if err != nil {
			return err
		}
	}
	return nil
}
