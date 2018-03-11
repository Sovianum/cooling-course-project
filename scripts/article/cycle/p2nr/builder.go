package p2nr

import (
	"github.com/Sovianum/cooling-course-project/scripts/article/cycle/p2n"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/helper"
	"github.com/Sovianum/turbocycle/library/parametric/free2n"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/cooling-course-project/scripts/article/cycle/common"
)

func NewBuilder(
	source schemes.TwoShaftsRegeneratorScheme,
	power,
	cRpm0, lambdaIn0,
	ctInletMeanDiameter, ctLambdaU0, ctStageNum,
	ftInletMeanDiameter, ftLambdaU0, ftStageNum,
	payloadRpm0, etaM,
	velocityHotIn0, velocityColdIn0,
	hydraulicDiameterHot, hydraulicDiameterCold float64,
	nuColdFunc, nuHotFunc constructive.NuFunc,
	tDropFunc constructive.TemperatureDropFunc,
	precision, relaxCoef float64, iterLimit int,
) *Builder {
	return &Builder{
		VelocityHotIn0:        velocityHotIn0,
		VelocityColdIn0:       velocityColdIn0,
		HydraulicDiameterHot:  hydraulicDiameterHot,
		HydraulicDiameterCold: hydraulicDiameterCold,
		NuColdFunc:            nuColdFunc, NuHotFunc: nuHotFunc,
		TDropFunc: tDropFunc,
		Builder: &p2n.Builder{
			Source:              source,
			Power:               power,
			CRpm0:               cRpm0,
			LambdaIn0:           lambdaIn0,
			CtInletMeanDiameter: ctInletMeanDiameter,
			CtLambdaU0:          ctLambdaU0,
			CtStageNum:          ctStageNum,
			FtInletMeanDiameter: ftInletMeanDiameter,
			FtLambdaU0:          ftLambdaU0,
			FtStageNum:          ftStageNum,
			PayloadRpm0:         payloadRpm0,
			EtaM:                etaM,
			Precision:           precision,
			RelaxCoef:           relaxCoef,
			IterLimit:           iterLimit,
		},
	}
}

type Builder struct {
	*p2n.Builder
	VelocityHotIn0        float64
	VelocityColdIn0       float64
	HydraulicDiameterHot  float64
	HydraulicDiameterCold float64
	NuColdFunc            constructive.NuFunc
	NuHotFunc             constructive.NuFunc
	TDropFunc             constructive.TemperatureDropFunc
}

func (b *Builder) Build() free2n.DoubleShaftRegFreeScheme {
	return free2n.NewDoubleShaftRegFreeScheme(
		b.Source.GasSource().GasOutput().GetState().Value().(gases.Gas),
		b.Source.Compressor().TStagIn(),
		b.Source.Compressor().PStagIn(),
		b.Source.Burner().TStagOut(),
		b.EtaM, b.BuildCompressor(), b.BuildCompressorPipe(),
		b.buildRegenerator(), b.buildCycleBreaker(),
		b.BuildBurner(), b.BuildCompressorTurbine(), b.BuildCTPipe(),
		b.BuildFreeTurbine(), b.BuildFreeTurbinePipe(), b.BuildPayload(),
	)
}

func (b *Builder) buildRegenerator() constructive.RegeneratorNode {
	s := b.Source.(schemes.TwoShaftsRegeneratorScheme)
	r := s.Regenerator()
	hi := r.HotInput()
	ci := r.ColdInput()

	hotMassRate := common.GetMassRate(b.Power, b.Source, b.Source.Compressor())
	coldMassRate := common.GetMassRate(b.Power, b.Source, b.Source.FreeTurbineBlock().FreeTurbine())

	return constructive.NewParametricRegeneratorNode(
		hi.GasInput().GetState().Value().(gases.Gas),
		ci.GasInput().GetState().Value().(gases.Gas),
		hotMassRate, coldMassRate,
		hi.TemperatureInput().GetState().Value().(float64),
		ci.TemperatureInput().GetState().Value().(float64),
		hi.PressureInput().GetState().Value().(float64),
		ci.PressureInput().GetState().Value().(float64),
		b.VelocityHotIn0, b.VelocityColdIn0,
		r.Sigma(),
		b.HydraulicDiameterHot,
		b.HydraulicDiameterCold,
		b.Precision, 1, nodes.DefaultN,
		b.TDropFunc,
		b.NuHotFunc, b.NuColdFunc,
	)
}

func (b *Builder) buildCycleBreaker() helper.ComplexCycleBreakNode {
	reg := b.Source.(schemes.TwoShaftsRegeneratorScheme).Regenerator()
	return helper.NewComplexCycleBreakNode(
		reg.HotInput().GasInput().GetState().Value().(gases.Gas),
		reg.HotInput().TemperatureInput().GetState().Value().(float64),
		reg.HotInput().PressureInput().GetState().Value().(float64),
		schemes.GetMassRate(b.Power, b.Source),
	)
}
