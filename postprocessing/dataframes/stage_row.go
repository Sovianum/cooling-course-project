package dataframes

import "fmt"

func NewTurbineStageRow(name, dimension string, dfs []TurbineStageDF, extractor func(df TurbineStageDF) float64) StageRow {
	result := StageRow{Name: name, Dimension: dimension, Values: make([]float64, len(dfs))}
	for i, df := range dfs {
		result.Values[i] = extractor(df)
	}
	return result
}

func NewCompressorStageRow(name, dimension string, dfs []CompressorStageDF, extractor func(df CompressorStageDF) float64) StageRow {
	result := StageRow{Name: name, Dimension: dimension, Values: make([]float64, len(dfs))}
	for i, df := range dfs {
		result.Values[i] = extractor(df)
	}
	return result
}

type StageRow struct {
	ID        int
	Name      string
	Dimension string
	Values    []float64

	floatFormatters []func(float64) float64
	stringFormatter func(float64) string
}

func (row StageRow) FormatFloat(f func(float64) float64) StageRow {
	row.floatFormatters = append(row.floatFormatters, f)
	return row
}

func (row StageRow) FormatString(f func(float64) string) StageRow {
	if row.stringFormatter != nil {
		panic(fmt.Errorf("string formatter already set"))
	}
	row.stringFormatter = f
	return row
}

func (row StageRow) GetStr() string {
	if row.stringFormatter == nil {
		row.stringFormatter = func(f float64) string {
			return fmt.Sprintf("%f", f)
		}
	}

	result := fmt.Sprintf("%d & %s & %s &", row.ID, row.Name, row.Dimension)
	for _, v := range row.Values[:len(row.Values)-1] {
		result += fmt.Sprintf(" %s &", row.getFormatFloat(v))
	}
	last := fmt.Sprintf(" %s", row.getFormatFloat(row.Values[len(row.Values)-1]))
	result += last
	return result
}

func (row StageRow) getFormatFloat(val float64) string {
	for _, ff := range row.floatFormatters {
		val = ff(val)
	}
	return row.stringFormatter(val)
}
