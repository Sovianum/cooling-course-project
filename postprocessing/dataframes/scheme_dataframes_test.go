package dataframes

import (
	"testing"
	"github.com/Sovianum/cooling-course-project/core/schemes/three_shafts"
	"github.com/Sovianum/cooling-course-project/core"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"os"
	"io/ioutil"
	templ2 "github.com/Sovianum/cooling-course-project/postprocessing/templ"
)

const (
	power = 16000e3
	relaxCoef = 0.1
	iterNum = 100

	cycleTemplateFilePath = "../templates/cycle_calc_template.tex"
)

func TestNewThreeShaftsDF_Smoke(t *testing.T) {
	var scheme = three_shafts.GetInitedThreeShaftsScheme()
	var pi = 10.
	var piFactor = 0.5

	var generator = core.GetDoubleCompressorDataGenerator(scheme, power, relaxCoef, iterNum)
	_, err := generator(pi, piFactor)
	assert.Nil(t, err)

	var df = NewThreeShaftsDF(power, scheme)
	var b, _ = json.MarshalIndent(df, "", "    ")
	os.Stdout.Write(b)
}

func TestTemplateSmoke(t *testing.T) {
	var f, fileErr = ioutil.ReadFile(cycleTemplateFilePath)
	assert.Nil(t, fileErr)

	var funcMap = templ2.GetFuncMap()
	var templ, tErr = templ2.GetTemplate(
		"stage",
		string(f),
		funcMap,
	)
	if tErr != nil {
		panic(tErr)
	}

	var scheme = three_shafts.GetInitedThreeShaftsScheme()
	var pi = 10.
	var piFactor = 0.5
	var iterNum = 100

	var generator = core.GetDoubleCompressorDataGenerator(scheme, power, relaxCoef, iterNum)
	var _, err = generator(pi, piFactor)
	if err != nil {
		panic(err)
	}
	var df = NewThreeShaftsDF(power, scheme)

	fileErr = templ.Execute(os.Stdout, &df)
	if fileErr != nil {
		panic(fileErr)
	}
}
