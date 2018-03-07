package templ

import (
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"math"
	"text/template"
)

func GetTemplate(name, content string, funcMap template.FuncMap) (*template.Template, error) {
	return template.
		New(name).
		Delims("<-<", ">->").
		Funcs(funcMap).
		Parse(content)
}

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
		"Degree": func(value float64) float64 {
			return common.ToDegrees(value)
		},
		"Radian": func(value float64) float64 {
			return common.ToRadians(value)
		},
	}
}
