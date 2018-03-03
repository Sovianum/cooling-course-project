package article

import (
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/methodics"
	"github.com/Sovianum/turbocycle/library/parametric/free2n"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/core/graph"
)

func NewParametric2NBuilder(
	source *schemes.TwoShaftsSchemeImpl,
	power, t0, p0,
	cRpm0, lambdaIn0,
	ctInletMeanDiameter, ctLambdaU0, ctStageNum,
	ftInletMeanDiameter, ftLambdaU0, ftStageNum,
	payloadRpm0, etaM,
	precision, relaxCoef float64, iterLimit int,
) *Parametric2NBuilder {
	return &Parametric2NBuilder{
		source: source,
		power:  power,
		t0:     t0, p0: p0,
		cRpm0: cRpm0, lambdaIn0: lambdaIn0,
		ctInletMeanDiameter: ctInletMeanDiameter, ctLambdaU0: ctLambdaU0, ctStageNum: ctStageNum,
		ftInletMeanDiameter: ftInletMeanDiameter, ftLambdaU0: ftLambdaU0, ftStageNum: ftStageNum,
		payloadRpm0: payloadRpm0, etaM: etaM,
		precision: precision, relaxCoef: relaxCoef, iterLimit: iterLimit,
	}
}

type Parametric2NBuilder struct {
	source *schemes.TwoShaftsSchemeImpl
	power  float64
	t0     float64
	p0     float64

	cRpm0 float64

	lambdaIn0 float64

	ctInletMeanDiameter float64
	ctLambdaU0          float64
	ctStageNum          float64

	ftInletMeanDiameter float64
	ftLambdaU0          float64
	ftStageNum          float64

	payloadRpm0 float64

	etaM float64

	precision float64
	relaxCoef float64
	iterLimit int
}

func (b *Parametric2NBuilder) Build() free2n.DoubleShaftFreeScheme {
	return free2n.NewDoubleShaftFreeScheme(
		b.source.GasSource().GasOutput().GetState().Value().(gases.Gas),
		b.t0, b.p0, b.source.GasGenerator().Burner().TStagOut(),
		b.etaM, b.buildCompressor(), b.buildCompressorPipe(),
		b.buildBurner(), b.buildCompressorTurbine(), b.buildFreeTurbinePipe(),
		b.buildFreeTurbine(), b.buildFreeTurbinePipe(), b.buildPayload(),
	)
}

func (b *Parametric2NBuilder) buildCompressor() constructive.ParametricCompressorNode {
	c := b.source.Compressor()
	ccGen := methodics.NewCompressorCharGen(
		c.PiStag(), c.Eta(), b.massRate(), b.precision, b.relaxCoef, b.iterLimit,
	)
	p0 := c.PStagIn()
	t0 := c.TStagIn()

	p := constructive.NewParametricCompressorNode(
		b.massRate(), c.PiStag(),
		b.cRpm0, c.Eta(), t0, p0, b.precision,
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

func (b *Parametric2NBuilder) buildCompressorPipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(1) // todo extract data about pipe from somewhere
}

func (b *Parametric2NBuilder) buildBurner() constructive.ParametricBurnerNode {
	burn := b.source.GasGenerator().Burner()
	pBurn := constructive.NewParametricBurnerNode(
		burn.Fuel(), burn.TFuel(), burn.T0(), burn.Eta(),
		b.lambdaIn0, burn.PStagIn(), burn.TStagIn(),
		b.massRate()*burn.MassRateInput().GetState().Value().(float64),
		burn.FuelRateRel(), b.precision, func(lambda float64) float64 {
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

func (b *Parametric2NBuilder) buildCompressorTurbine() constructive.ParametricTurbineNode {
	ct := b.source.GasGenerator().TurboCascade().Turbine()
	tcGen := methodics.NewKazandjanTurbineCharacteristic()
	p0 := ct.PStagIn()
	t0 := ct.TStagIn()

	pt := constructive.NewParametricTurbineNode(
		b.massRate()*ct.MassRateInput().GetState().Value().(float64),
		ct.PiTStag(), ct.Eta(), t0, p0, b.ctInletMeanDiameter, b.precision,
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

func (b *Parametric2NBuilder) buildCTPipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(b.source.CompressorTurbinePipe().Sigma())
}

func (b *Parametric2NBuilder) buildFreeTurbine() constructive.ParametricTurbineNode {
	ft := b.source.FreeTurbineBlock().FreeTurbine()
	tcGen := methodics.NewKazandjanTurbineCharacteristic()
	p0 := ft.PStagIn()
	t0 := ft.TStagIn()

	pt := constructive.NewParametricTurbineNode(
		b.massRate()*ft.MassRateInput().GetState().Value().(float64),
		ft.PiTStag(), ft.Eta(), t0, p0, b.ftInletMeanDiameter, b.precision,
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

func (b *Parametric2NBuilder) buildFreeTurbinePipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(b.source.FreeTurbineBlock().OutletPressureLoss().Sigma())
}

func (b *Parametric2NBuilder) buildPayload() constructive.Payload {
	return constructive.NewPayload(
		b.payloadRpm0, b.power, func(normRpm float64) float64 {
			//delta := normRpm - 1
			//return normRpm - delta * delta
			return normRpm * normRpm * normRpm // todo add smth more precise
		},
	)
}

func (b *Parametric2NBuilder) massRate() float64 {
	return schemes.GetMassRate(b.power, b.source)
}

func copyAll(p1s, p2s []graph.Port) {
	for i, p1 := range p1s {
		copyState(p1, p2s[i])
	}
}

func copyState(p1, p2 graph.Port) {
	p2.SetState(p1.GetState())
}
