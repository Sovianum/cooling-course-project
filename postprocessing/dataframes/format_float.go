package dataframes

import (
	"fmt"
	"strings"
)

func changeDecimal(string string) string {
	return strings.Replace(string, ".", ",", -1)
}

func Round(value float64) string {
	return changeDecimal(fmt.Sprintf("%.0f", value))
}

func Round1(value float64) string {
	return changeDecimal(fmt.Sprintf("%.1f", value))
}

func Round2(value float64) string {
	return changeDecimal(fmt.Sprintf("%.2f", value))
}

func Round3(value float64) string {
	return changeDecimal(fmt.Sprintf("%.3f", value))
}

func DivideE3(value float64) float64 {
	return value / 1e3
}

func MultiplyE3(value float64) float64 {
	return value * 1e3
}

func DivideE5(value float64) float64 {
	return value / 1e5
}

func MultiplyE5(value float64) float64 {
	return value * 1e5
}

func DivideE6(value float64) float64 {
	return value / 1e6
}

func MultiplyE6(value float64) float64 {
	return value * 1e6
}

type FormatFloat float64

func (f FormatFloat) FormatFloat(formatter func(f float64) float64) FormatFloat {
	return FormatFloat(formatter(float64(f)))
}

func (f FormatFloat) FormatString(formatter func(f float64) string) string {
	return formatter(float64(f))
}
