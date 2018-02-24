package dataframes

import (
	"encoding/json"
	templ2 "github.com/Sovianum/cooling-course-project/postprocessing/templ"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/turbine/geometry"
	"github.com/Sovianum/turbocycle/impl/turbine/nodes"
	states2 "github.com/Sovianum/turbocycle/impl/turbine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
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
	rotorApproxTRel = 0.7

	alpha = 14

	meanLineTemplateFilePath = "../templates/mean_line_calc_template.tex"
)

type StageDFTestSuite struct {
	suite.Suite
	node nodes.TurbineStageNode
	gen  geometry.StageGeometryGenerator
	df   StageDF
}

func (suite *StageDFTestSuite) SetupTest() {
	suite.gen = geometry.NewStageGeometryGenerator(
		lRelOut,
		geometry.NewIncompleteGeneratorFromProfileAngles(baRel, deltaRel, gammaIn, gammaOut, statorApproxTRel),
		geometry.NewIncompleteGeneratorFromProfileAngles(baRel, deltaRel, gammaIn, gammaOut, rotorApproxTRel),
	)

	suite.node = nodes.NewTurbineStageNode(
		n, stageHeatDrop, reactivity, phi, psi, airGapRel, precision, suite.gen,
	)

	suite.node.GasInput().SetState(states.NewGasPortState(gases.GetAir()))
	suite.node.VelocityInput().SetState(states2.NewVelocityPortState(
		states2.NewInletTriangle(0, c0, math.Pi/2),
		states2.InletTriangleType,
	))

	suite.node.TemperatureInput().SetState(states.NewTemperaturePortState(tg))
	suite.node.PressureInput().SetState(states.NewPressurePortState(pg))
	suite.node.MassRateInput().SetState(states.NewMassRatePortState(massRate))

	suite.node.SetAlpha1FirstStage(common.ToRadians(alpha))

	suite.node.Process()
	suite.df, _ = NewStageDF(suite.node)
}

func (suite *StageDFTestSuite) TestSmoke() {
	if _, err := json.MarshalIndent(suite.df, "", "    "); err != nil {
		panic(err)
	}
}

func (suite *StageDFTestSuite) TestTemplateSmoke() {
	var f, err = ioutil.ReadFile(meanLineTemplateFilePath)
	assert.Nil(suite.T(), err)

	var funcMap = templ2.GetFuncMap()
	var templ, tErr = templ2.GetTemplate(
		"stage",
		string(f),
		funcMap,
	)
	assert.Nil(suite.T(), tErr)

	err = templ.Execute(ioutil.Discard, &suite.df)
	assert.Nil(suite.T(), err)
}

func TestStageNodeTestSuite(t *testing.T) {
	suite.Run(t, new(StageDFTestSuite))
}
