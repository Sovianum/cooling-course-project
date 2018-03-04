package p2n

import (
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/methodics"
	"github.com/Sovianum/turbocycle/library/parametric/free2n"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/core/graph"
)

func NewBuilder(
	source schemes.TwoShaftsScheme,
	power, t0, p0,
	cRpm0, lambdaIn0,
	ctInletMeanDiameter, ctLambdaU0, ctStageNum,
	ftInletMeanDiameter, ftLambdaU0, ftStageNum,
	payloadRpm0, etaM,
	precision, relaxCoef float64, iterLimit int,
) *Builder {
	return &Builder{
		Source:              source,
		Power:               power,
		T0:                  t0,
		P0:                  p0,
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
	}
}

type Builder struct {
	Source schemes.TwoShaftsScheme
	Power  float64
	T0     float64
	P0     float64

	CRpm0 float64

	LambdaIn0 float64

	CtInletMeanDiameter float64
	CtLambdaU0          float64
	CtStageNum          float64

	FtInletMeanDiameter float64
	FtLambdaU0          float64
	FtStageNum          float64

	PayloadRpm0 float64

	EtaM float64

	Precision float64
	RelaxCoef float64
	IterLimit int
}

func (b *Builder) Build() free2n.DoubleShaftFreeScheme {
	return free2n.NewDoubleShaftFreeScheme(
		b.Source.GasSource().GasOutput().GetState().Value().(gases.Gas),
		b.T0, b.P0, b.Source.Burner().TStagOut(),
		b.EtaM, b.BuildCompressor(), b.BuildCompressorPipe(),
		b.BuildBurner(), b.BuildCompressorTurbine(), b.BuildFreeTurbinePipe(),
		b.BuildFreeTurbine(), b.BuildFreeTurbinePipe(), b.BuildPayload(),
	)
}

func (b *Builder) BuildCompressor() constructive.ParametricCompressorNode {
	c := b.Source.Compressor()
	ccGen := methodics.NewCompressorCharGen(
		c.PiStag(), c.Eta(), b.massRate(), b.Precision, b.RelaxCoef, b.IterLimit,
	)
	p0 := c.PStagIn()
	t0 := c.TStagIn()

	p := constructive.NewParametricCompressorNode(
		b.massRate(), c.PiStag(),
		b.CRpm0, c.Eta(), t0, p0, b.Precision,
		ccGen.GetNormEtaChar(),
		ccGen.GetNormRPMChar(),
	)

	copyAll(
		[]graph.Port{
			c.GasInput(), c.TemperatureInput(), c.PressureInput(), c.MassRateInput(),
			c.GasOutput(), c.TemperatureOutput(), c.PressureOutput(), c.MassRateOutput(),
		},
		[]graph.Port{
			p.GasInput(), p.TemperatureInput(), p.PressureInput(), p.MassRateInput(),
			p.GasOutput(), p.TemperatureOutput(), p.PressureOutput(), p.MassRateOutput(),
		},
	)
	return p
}

func (b *Builder) BuildCompressorPipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(1) // todo extract data about pipe from somewhere
}

func (b *Builder) BuildBurner() constructive.ParametricBurnerNode {
	burn := b.Source.Burner()
	pBurn := constructive.NewParametricBurnerNode(
		burn.Fuel(), burn.TFuel(), burn.T0(), burn.Eta(),
		b.LambdaIn0, burn.PStagIn(), burn.TStagIn(),
		b.massRate()*burn.MassRateInput().GetState().Value().(float64),
		burn.FuelRateRel(), b.Precision, func(lambda float64) float64 {
			return burn.Sigma() // todo make something more precise
		},
	)

	copyAll(
		[]graph.Port{
			burn.GasInput(), burn.TemperatureInput(), burn.PressureInput(), burn.MassRateInput(),
			burn.GasOutput(), burn.TemperatureOutput(), burn.PressureOutput(), burn.MassRateOutput(),
		},
		[]graph.Port{
			pBurn.GasInput(), pBurn.TemperatureInput(), pBurn.PressureInput(), pBurn.MassRateInput(),
			pBurn.GasOutput(), pBurn.TemperatureOutput(), pBurn.PressureOutput(), pBurn.MassRateOutput(),
		},
	)
	return pBurn
}

func (b *Builder) BuildCompressorTurbine() constructive.ParametricTurbineNode {
	ct := b.Source.TurboCascade().Turbine()
	tcGen := methodics.NewKazandjanTurbineCharacteristic()
	p0 := ct.PStagIn()
	t0 := ct.TStagIn()

	pt := constructive.NewParametricTurbineNode(
		b.massRate()*ct.MassRateInput().GetState().Value().(float64),
		ct.PiTStag(), ct.Eta(), t0, p0, b.CtInletMeanDiameter, b.Precision,
		func(node constructive.TurbineNode) float64 {
			return ct.LeakMassRateRel()
		},
		func(node constructive.TurbineNode) float64 {
			return ct.CoolMassRateRel()
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		tcGen.GetNormMassRateChar(),
		tcGen.GetNormEtaChar(),
	)

	copyAll(
		[]graph.Port{
			ct.GasInput(), ct.TemperatureInput(), ct.PressureInput(),
			ct.GasOutput(), ct.TemperatureOutput(), ct.PressureOutput(), ct.MassRateOutput(),
		},
		[]graph.Port{
			pt.GasInput(), pt.TemperatureInput(), pt.PressureInput(),
			pt.GasOutput(), pt.TemperatureOutput(), pt.PressureOutput(), pt.MassRateOutput(),
		},
	)
	return pt
}

func (b *Builder) BuildCTPipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(b.Source.CompressorTurbinePipe().Sigma())
}

func (b *Builder) BuildFreeTurbine() constructive.ParametricTurbineNode {
	ft := b.Source.FreeTurbineBlock().FreeTurbine()
	tcGen := methodics.NewKazandjanTurbineCharacteristic()
	p0 := ft.PStagIn()
	t0 := ft.TStagIn()

	pt := constructive.NewParametricTurbineNode(
		b.massRate()*ft.MassRateInput().GetState().Value().(float64),
		ft.PiTStag(), ft.Eta(), t0, p0, b.FtInletMeanDiameter, b.Precision,
		func(node constructive.TurbineNode) float64 {
			return ft.LeakMassRateRel()
		},
		func(node constructive.TurbineNode) float64 {
			return ft.CoolMassRateRel()
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		tcGen.GetNormMassRateChar(),
		tcGen.GetNormEtaChar(),
	)
	copyAll(
		[]graph.Port{
			ft.GasInput(), ft.TemperatureInput(), ft.PressureInput(),
			ft.GasOutput(), ft.TemperatureOutput(), ft.PressureOutput(), ft.MassRateOutput(),
		},
		[]graph.Port{
			pt.GasInput(), pt.TemperatureInput(), pt.PressureInput(),
			pt.GasOutput(), pt.TemperatureOutput(), pt.PressureOutput(), pt.MassRateOutput(),
		},
	)
	return pt
}

func (b *Builder) BuildFreeTurbinePipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(b.Source.FreeTurbineBlock().OutletPressureLoss().Sigma())
}

func (b *Builder) BuildPayload() constructive.Payload {
	return constructive.NewPayload(
		b.PayloadRpm0, b.Power, func(normRpm float64) float64 {
			//delta := normRpm - 1
			//return normRpm - delta * delta
			return normRpm * normRpm * normRpm // todo add smth more precise
		},
	)
}

func (b *Builder) massRate() float64 {
	return schemes.GetMassRate(b.Power, b.Source)
}

func copyAll(p1s, p2s []graph.Port) {
	for i, p1 := range p1s {
		copyState(p1, p2s[i])
	}
}

func copyState(p1, p2 graph.Port) {
	p2.SetState(p1.GetState())
}
