package templ

import (
	"text/template"
	"fmt"
	"math"
)

func GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"Round": func(value float64) string {
			return fmt.Sprintf("%.0f", value)
		},
		"Round1": func(value float64) string {
			return fmt.Sprintf("%.1f", value)
		},
		"Round2": func(value float64) string {
			return fmt.Sprintf("%.2f", value)
		},
		"Round3": func(value float64) string {
			return fmt.Sprintf("%.3f", value)
		},
		"DivideE3": func(value float64) float64 {
			return value / 1e3
		},
		"MultiplyE3": func(value float64) float64 {
			return value * 1e3
		},
		"DivideE5": func(value float64) float64 {
			return value / 1e5
		},
		"MultiplyE5": func(value float64) float64 {
			return value * 1e5
		},
		"DivideE6": func(value float64) float64 {
			return value / 1e6
		},
		"MultiplyE6": func(value float64) float64 {
			return value * 1e6
		},
		"Abs": func(value float64) float64 {
			return math.Abs(value)
		},
	}
}