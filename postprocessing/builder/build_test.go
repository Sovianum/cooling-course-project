package builder

import (
	"github.com/Sovianum/cooling-course-project/core"
	"github.com/Sovianum/cooling-course-project/core/schemes/three_shafts"
	"github.com/Sovianum/cooling-course-project/postprocessing/dataframes"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/turbine/geometry"
	"github.com/Sovianum/turbocycle/impl/turbine/nodes"
	states2 "github.com/Sovianum/turbocycle/impl/turbine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/suite"
	"math"
	"testing"
)

const (
	n             = 1e4
	stageHeatDrop = 3e5
	reactivity    = 0.5
	phi           = 0.98
	psi           = 0.98
	airGapRel     = 0.001
	precision     = 0.05

	c0       = 50.
	tg       = 1200.
	pg       = 1e6
	massRate = 100.

	gammaIn  = -0.09
	gammaOut = 0.09
	baRel    = 4
	lRelOut  = 0.15
	deltaRel = 0.1

	statorApproxTRel = 0.7
	rotorApproxTRel  = 0.7

	alpha = 14

	meanLineTemplateFilePath = "../templates/mean_line_calc_template.tex"
)

type insertionInfo struct {
	outputPath string
	data       interface{}
}

type BuildTestSuite struct {
	suite.Suite
	insertionMap map[string]insertionInfo
}

func (suite *BuildTestSuite) SetupTest() {
	suite.insertionMap = getInsertionMap()
}

func TestStageNodeTestSuite(t *testing.T) {
	suite.Run(t, new(BuildTestSuite))
}

func getInsertionMap() map[string]insertionInfo {
	var templateDir = "/home/artem/gowork/src/github.com/Sovianum/cooling-course-project/postprocessing/templates"
	var outputDir = "/home/artem/gowork/src/github.com/Sovianum/cooling-course-project/build"

	var getTemplateFile = func(name string) string { return templateDir + "/" + name }
	var getOutputFile = func(name string) string { return outputDir + "/" + name }

	var result = make(map[string]insertionInfo)

	result[getTemplateFile("root.tex")] = insertionInfo{
		outputPath: getOutputFile("root.tex"),
	}
	result[getTemplateFile("var_list_template.tex")] = insertionInfo{
		outputPath: getOutputFile("var_list.tex"),
	}
	result[getTemplateFile("cycle_calc_template.tex")] = insertionInfo{
		outputPath: getOutputFile("cycle_calc.tex"),
		data:       getCycleDf(),
	}
	result[getTemplateFile("mean_line_calc_template.tex")] = insertionInfo{
		outputPath: getOutputFile("mean_line_calc.tex"),
		data:       getStageDf(),
	}

	return result
}

func getCycleDf() dataframes.ThreeShaftsDF {
	var scheme = three_shafts.GetInitedThreeShaftsScheme()
	var pi = 10.
	var piFactor = 0.5
	var iterNum = 100
	var nE = 16e6
	var etaR = 0.98

	var generator = core.GetDoubleCompressorDataGenerator(scheme, nE, 0.1, iterNum)
	var _, err = generator(pi, piFactor)
	if err != nil {
		panic(err)
	}
	return dataframes.NewThreeShaftsDF(nE, etaR, scheme)
}

func getStageDf() dataframes.StageDF {
	var gen = geometry.NewStageGeometryGenerator(
		lRelOut,
		geometry.NewIncompleteGeneratorFromProfileAngles(baRel, deltaRel, gammaIn, gammaOut, statorApproxTRel),
		geometry.NewIncompleteGeneratorFromProfileAngles(baRel, deltaRel, gammaIn, gammaOut, rotorApproxTRel),
	)

	var stage = nodes.NewTurbineStageNode(
		n, stageHeatDrop, reactivity, phi, psi, airGapRel, precision, gen,
	)

	stage.GasInput().SetState(states.NewGasPortState(gases.GetAir()))
	stage.VelocityInput().SetState(states2.NewVelocityPortState(
		states2.NewInletTriangle(0, c0, math.Pi/2),
		states2.InletTriangleType,
	))

	stage.TemperatureInput().SetState(states.NewTemperaturePortState(tg))
	stage.PressureInput().SetState(states.NewPressurePortState(pg))
	stage.MassRateInput().SetState(states.NewMassRatePortState(massRate))

	stage.SetAlpha1FirstStage(common.ToRadians(alpha))

	var err = stage.Process()
	if err != nil {
		panic(err)
	}
	var df, _ = dataframes.NewStageDF(stage)
	return df
}
