package inited

import (
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/stage/compressor"
	"github.com/Sovianum/turbocycle/impl/stage/geometry"
	"github.com/Sovianum/turbocycle/impl/stage/turbine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetInitedStagedNodes(t *testing.T) {
	data, err := GetInitedStagedNodes()
	assert.NoError(t, err)
	if err != nil {
		return
	}

	lpc := data.LPC
	fmt.Println("LPC")
	fmt.Println(getCompressorMessage(lpc))

	hpc := data.HPC
	fmt.Println("HPC")
	fmt.Println(getCompressorMessage(hpc))

	hpt := data.HPT
	fmt.Println("HPT")
	fmt.Println(getTurbineMessage(hpt))

	lpt := data.LPT
	fmt.Println("LPT")
	fmt.Println(getTurbineMessage(lpt))

	ft := data.FT
	fmt.Println("FT")
	fmt.Println(getTurbineMessage(ft))
}

func getCompressorMessage(compressor compressor.StagedCompressorNode) string {
	result := ""
	for _, stage := range compressor.Stages() {
		pack := stage.GetDataPack()
		inletTriangle := pack.InletTriangle
		midTriangle := pack.MidTriangle
		outletTriangle := pack.OutletTriangle

		result += fmt.Sprintf(
			"alpha1: %.3f, beta1: %.3f, alpha2: %.3f, beta2: %.3f, alpha3: %.3f, pi: %.3f, u: %.1f, ca1: %.3f, ca3: %.3f, ht: %.3f, dOut1: %.3f, dIn1: %.3f, dx: %.3f, dRelIn: %.3f\n",
			common.ToDegrees(inletTriangle.Alpha()),
			common.ToDegrees(inletTriangle.Beta()),
			common.ToDegrees(midTriangle.Alpha()),
			common.ToDegrees(midTriangle.Beta()),
			common.ToDegrees(outletTriangle.Alpha()),
			pack.PiStag,
			pack.UOut,
			pack.InletTriangle.CA(),
			pack.OutletTriangle.CA(),
			pack.HTCoef,
			pack.StageGeometry.RotorGeometry().OuterProfile().Diameter(0),
			pack.StageGeometry.RotorGeometry().InnerProfile().Diameter(0),
			pack.StageGeometry.RotorGeometry().XGapOut()+pack.StageGeometry.StatorGeometry().XGapOut(),
			geometry.DRel(0, pack.StageGeometry.RotorGeometry()),
		)
	}
	return result
}

func getTurbineMessage(turbine turbine.StagedTurbineNode) string {
	result := ""
	for _, stage := range turbine.Stages() {
		pack := stage.GetDataPack()
		rotorInletTriangle := pack.RotorInletTriangle
		rotorOutletTriangle := pack.RotorOutletTriangle

		lRelOutFunc := func(bladingGeom geometry.BladingGeometry) float64 {
			dOut := bladingGeom.OuterProfile().Diameter(bladingGeom.XGapOut())
			dIn := bladingGeom.InnerProfile().Diameter(bladingGeom.XGapOut())
			l := (dOut - dIn) / 2
			dMean := (dOut + dIn) / 2
			return l / dMean
		}
		result += fmt.Sprintf(
			"alpha1: %.3f, beta1: %.3f, alpha2: %.3f, beta2: %.3f, pi: %.3f, eta: %.3f, u1: %.1f, ca1: %.3f, ca2: %.3f, dMean: %.3f, dOut: %.3f, dIn: %.3f, dx: %.3f, lRelOut: %.3f\n",
			common.ToDegrees(rotorInletTriangle.Alpha()),
			common.ToDegrees(rotorInletTriangle.Beta()),
			common.ToDegrees(rotorOutletTriangle.Alpha()),
			common.ToDegrees(rotorOutletTriangle.Beta()),
			pack.PiStag,
			pack.EtaTStag,
			pack.U1,
			rotorInletTriangle.CA(),
			rotorOutletTriangle.CA(),
			pack.StageGeometry.RotorGeometry().MeanProfile().Diameter(0),
			pack.StageGeometry.RotorGeometry().OuterProfile().Diameter(0),
			pack.StageGeometry.RotorGeometry().InnerProfile().Diameter(0),
			pack.StageGeometry.StatorGeometry().XGapOut()+pack.StageGeometry.RotorGeometry().XGapOut(),
			lRelOutFunc(pack.StageGeometry.RotorGeometry()),
		)
	}
	return result
}