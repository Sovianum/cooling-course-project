package p3n

import (
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/library/parametric/free3n"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/cooling-course-project/scripts/article/cycle/p3n"
)

func NewBuilder(
	source schemes.ThreeShaftsCoolerScheme,
	power, t0, p0,
	lpcRpm0, hpcRpm0,
	lambdaIn0 float64,
	lptInletMeanDiameter, lptLambdaU0, lptStageNum,
	hptInletMeanDiameter, hptLambdaU0, hptStageNum,
	ftInletMeanDiameter, ftLambdaU0, ftStageNum,
	payloadRpm0,
	lpEtaM, hpEtaM,
	precision, relaxCoef float64, iterLimit int,
) *Builder {
	return &Builder{
		Builder: p3n.NewBuilder(
			source,
			power, t0, p0,
			lpcRpm0, hpcRpm0,
			lambdaIn0,
			lptInletMeanDiameter, lptLambdaU0, lptStageNum,
			hptInletMeanDiameter, hptLambdaU0, hptStageNum,
			ftInletMeanDiameter, ftLambdaU0, ftStageNum,
			payloadRpm0,
			lpEtaM, hpEtaM,
			precision, relaxCoef,iterLimit,
		),
	}
}

type Builder struct {
	*p3n.Builder
	TCoolOut float64
	CoolerSigma float64
}

func (b *Builder) Build() free3n.ThreeShaftCoolFreeScheme {
	return free3n.NewThreeShaftCoolFreeScheme(
		b.Source.GasSource().GasOutput().GetState().Value().(gases.Gas),
		b.T0, b.P0, b.Source.MainBurner().TStagOut(),

		b.BuildLPC(), b.BuildLPCPipe(),
		b.BuildLPT(), b.BuildLPTPipe(),
		b.LPEtaM,

		b.BuildHPC(), b.BuildHPCPipe(),
		b.BuildHPT(), b.BuildHPTPipe(),
		b.HPEtaM,

		b.BuildFT(), b.BuildFTPipe(),
		b.BuildBurner(), b.BuildPayload(),
		b.buildCooler(),
	)
}

func (b *Builder) buildCooler() constructive.CoolerNode {
	c := b.Source.(schemes.ThreeShaftsCoolerScheme).Cooler()
	return constructive.NewCoolerNode(c.TemperatureOutput().GetState().Value().(float64), c.Sigma())
}
