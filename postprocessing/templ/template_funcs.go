package templ

import (
	"github.com/Sovianum/cooling-course-project/postprocessing/dataframes"
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
		"Round":      dataframes.Round,
		"Round1":     dataframes.Round1,
		"Round2":     dataframes.Round2,
		"Round3":     dataframes.Round3,
		"DivideE3":   dataframes.DivideE3,
		"MultiplyE3": dataframes.MultiplyE3,
		"DivideE5":   dataframes.DivideE5,
		"MultiplyE5": dataframes.MultiplyE5,
		"DivideE6":   dataframes.DivideE6,
		"MultiplyE6": dataframes.MultiplyE6,
		"Abs":        math.Abs,
		"Degree":     common.ToDegrees,
		"Radian":     common.ToRadians,
	}
}
