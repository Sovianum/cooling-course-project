package main

import (
	"io/ioutil"

	"github.com/Sovianum/cooling-course-project/postprocessing/dataframes"
	"os"
	"github.com/Sovianum/cooling-course-project/core/schemes/three_shafts"
	"github.com/Sovianum/cooling-course-project/core"
	templ2 "github.com/Sovianum/cooling-course-project/postprocessing/templ"
)

const (
	filePath = "postprocessing/templates/cycle_calc_template.tex"
)

func main() {
	testMain()
}

func testMain() {
	var f, err = ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var funcMap = templ2.GetFuncMap()
	var templ, tErr = templ2.GetTemplate("mean", string(f), funcMap)

	if tErr != nil {
		panic(tErr)
	}

	var data = getCycleDf()
	//err = templ.Execute(ioutil.Discard, data)
	err = templ.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

func getCycleDf() *dataframes.ThreeShaftsDF {
	var scheme = three_shafts.GetInitedThreeShaftsScheme()
	var pi = 10.
	var piFactor = 0.5
	var power = 16e6
	var relaxCoef = 0.1
	var iterNum = 100

	var generator = core.GetDoubleCompressorDataGenerator(scheme, power, relaxCoef, iterNum)
	_, err := generator(pi, piFactor)
	if err != nil {
		panic(err)
	}
	var result = dataframes.NewThreeShaftsDF(power, scheme)
	return &result
}
