package main

import (
	"io/ioutil"
	"text/template"

	"github.com/Sovianum/cooling-course-project/postprocessing/dataframes"
)

const (
	filePath = "core/postprocessing/templates/cycle_calc_template.tex"
)

func main() {
	var f, err = ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var templ, tErr = template.New("cycle").Delims("<-<", ">->").Parse(string(f))
	if tErr != nil {
		panic(tErr)
	}

	var data = new(dataframes.ThreeShaftsDF)
	err = templ.Execute(ioutil.Discard, data)
	if err != nil {
		panic(err)
	}
}
