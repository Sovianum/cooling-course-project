package profiling

import (
	"encoding/csv"
	"fmt"
	"os"
)

func SaveString(path, data string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(data)
	return err
}

func SaveMatrix(path string, matrix [][]float64) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.WriteAll(ToRecords(matrix)); err != nil {
		return err
	}
	return nil
}

func ToRecords(matrix [][]float64) [][]string {
	var result = make([][]string, len(matrix))
	for i, record := range matrix {
		result[i] = make([]string, len(record))
		for j, num := range record {
			result[i][j] = fmt.Sprintf("%f", num)
		}
	}
	return result
}
