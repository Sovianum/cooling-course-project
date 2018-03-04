package article

import (
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/helper"
	"github.com/Sovianum/turbocycle/library/parametric/free2n"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/turbocycle/material/gases"
)

func NewParametric2NRegBuilder(
	source schemes.TwoShaftsRegeneratorScheme,
	power, t0, p0,
	cRpm0, lambdaIn0,
	ctInletMeanDiameter, ctLambdaU0, ctStageNum,
	ftInletMeanDiameter, ftLambdaU0, ftStageNum,
	payloadRpm0, etaM,
	velocityHotIn0, velocityColdIn0,
	hydraulicDiameterHot, hydraulicDiameterCold float64,
	nuColdFunc, nuHotFunc constructive.NuFunc,
	tDropFunc constructive.TemperatureDropFunc,
	precision, relaxCoef float64, iterLimit int,
) *Parametric2NRegBuilder {
	return &Parametric2NRegBuilder{
		velocityHotIn0:        velocityHotIn0,
		velocityColdIn0:       velocityColdIn0,
		hydraulicDiameterHot:  hydraulicDiameterHot,
		hydraulicDiameterCold: hydraulicDiameterCold,
		nuColdFunc:            nuColdFunc, nuHotFunc: nuHotFunc,
		tDropFunc: tDropFunc,
		Parametric2NBuilder: &Parametric2NBuilder{
			source: source,
			power:  power,
			t0:     t0, p0: p0,
			cRpm0: cRpm0, lambdaIn0: lambdaIn0,
			ctInletMeanDiameter: ctInletMeanDiameter, ctLambdaU0: ctLambdaU0, ctStageNum: ctStageNum,
			ftInletMeanDiameter: ftInletMeanDiameter, ftLambdaU0: ftLambdaU0, ftStageNum: ftStageNum,
			payloadRpm0: payloadRpm0, etaM: etaM,
			precision: precision, relaxCoef: relaxCoef, iterLimit: iterLimit,
		},
	}
}

type Parametric2NRegBuilder struct {
	*Parametric2NBuilder
	velocityHotIn0        float64
	velocityColdIn0       float64
	hydraulicDiameterHot  float64
	hydraulicDiameterCold float64
	nuColdFunc            constructive.NuFunc
	nuHotFunc             constructive.NuFunc
	tDropFunc             constructive.TemperatureDropFunc
}

func (b *Parametric2NRegBuilder) Build() free2n.DoubleShaftRegFreeScheme {
	return free2n.NewDoubleShaftRegFreeScheme(
		b.source.GasSource().GasOutput().GetState().Value().(gases.Gas),
		b.t0, b.p0, b.source.Burner().TStagOut(),
		b.etaM, b.buildCompressor(), b.buildCompressorPipe(),
		b.buildRegenerator(), b.buildCycleBreaker(),
		b.buildBurner(), b.buildCompressorTurbine(), b.buildFreeTurbinePipe(),
		b.buildFreeTurbine(), b.buildFreeTurbinePipe(), b.buildPayload(),
	)
}

func (b *Parametric2NRegBuilder) buildRegenerator() constructive.RegeneratorNode {
	s := b.source.(schemes.TwoShaftsRegeneratorScheme)
	r := s.Regenerator()
	hi := r.HotInput()
	ci := r.ColdInput()
	return constructive.NewParametricRegeneratorNode(
		hi.GasInput().GetState().Value().(gases.Gas),
		ci.GasInput().GetState().Value().(gases.Gas),
		hi.MassRateInput().GetState().Value().(float64),
		ci.MassRateInput().GetState().Value().(float64),
		hi.TemperatureInput().GetState().Value().(float64),
		ci.TemperatureInput().GetState().Value().(float64),
		hi.PressureInput().GetState().Value().(float64),
		ci.PressureInput().GetState().Value().(float64),
		b.velocityHotIn0, b.velocityColdIn0,
		r.Sigma(),
		b.hydraulicDiameterHot,
		b.hydraulicDiameterCold,
		b.precision,
		b.tDropFunc,
		b.nuHotFunc, b.nuColdFunc,
	)
}

func (b *Parametric2NRegBuilder) buildCycleBreaker() helper.ComplexCycleBreakNode {
	ft := b.source.FreeTurbineBlock().FreeTurbine()
	return helper.NewComplexCycleBreakNode(
		ft.InputGas(),
		ft.TStagOut(),
		ft.PStagOut(),
		schemes.GetMassRate(b.power, b.source),
	)
}
