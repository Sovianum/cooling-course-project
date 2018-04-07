package inited

import (
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/stage/compressor"
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
	fmt.Println(getCompressorAngleMessage(lpc))

	hpc := data.HPC
	fmt.Println("HPC")
	fmt.Println(getCompressorAngleMessage(hpc))
}

func getCompressorAngleMessage(compressor compressor.StagedCompressorNode) string {
	result := ""
	for _, stage := range compressor.Stages() {
		pack := stage.GetDataPack()
		inletTriangle := pack.InletTriangle
		midTriangle := pack.MidTriangle
		outletTriangle := pack.OutletTriangle

		result += fmt.Sprintf(
			"alpha1: %.3f, beta1: %.3f, alpha2: %.3f, beta2: %.3f, alpha3: %.3f, pi: %.3f, u: %.1f, ca1: %.3f, ca3: %.3f, ht: %.3f, dOut1: %.3f\n",
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
		)
	}
	return result
}
