package common

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/methodics"
	"github.com/Sovianum/turbocycle/library/schemes"
)

func BuildCompressor(
	c constructive.CompressorNode, ccGen methodics.CompressorCharGen,
	rpm0, massRate0, precision float64,
) constructive.ParametricCompressorNode {
	//ccGen := methodics.NewCompressorCharGen(
	//	c.PiStag(), c.Eta(), massRate0, precision, relaxCoef, iterLimit,
	//)
	p0 := c.PStagIn()
	t0 := c.TStagIn()

	p := constructive.NewParametricCompressorNode(
		massRate0, c.PiStag(),
		rpm0, c.Eta(), t0, p0, precision,
		ccGen.GetNormEtaChar(),
		ccGen.GetNormRPMChar(),
	)

	CopyAll(
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

func BuildTurbine(
	t constructive.StaticTurbineNode, tChar methodics.TurbineCharacteristic,
	massRate0, inletDiameter, precision float64,
) constructive.ParametricTurbineNode {
	p0 := t.PStagIn()
	t0 := t.TStagIn()

	pt := constructive.NewParametricTurbineNode(
		massRate0,
		t.PiTStag(), t.Eta(), t0, p0, inletDiameter, precision,
		func(node constructive.TurbineNode) float64 {
			return t.LeakMassRateRel()
		},
		func(node constructive.TurbineNode) float64 {
			return t.CoolMassRateRel()
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		tChar.GetNormMassRateChar(),
		tChar.GetNormEtaChar(),
	)

	CopyAll(
		[]graph.Port{
			t.GasInput(), t.TemperatureInput(), t.PressureInput(),
			t.GasOutput(), t.TemperatureOutput(), t.PressureOutput(), t.MassRateOutput(),
		},
		[]graph.Port{
			pt.GasInput(), pt.TemperatureInput(), pt.PressureInput(),
			pt.GasOutput(), pt.TemperatureOutput(), pt.PressureOutput(), pt.MassRateOutput(),
		},
	)
	return pt
}

func BuildBurner(b constructive.BurnerNode, lambdaIn0, massRate0, precision, relaxCoef float64, iterLimit int) constructive.ParametricBurnerNode {
	pBurn := constructive.NewParametricBurnerNode(
		b.Fuel(), b.TFuel(), b.T0(), b.Eta(),
		lambdaIn0, b.PStagIn(), b.TStagIn(),
		massRate0, b.FuelRateRel(), precision, relaxCoef, iterLimit,
		func(lambda float64) float64 {
			return b.Sigma() // todo make something more precise
		},
	)

	CopyAll(
		[]graph.Port{
			b.GasInput(), b.TemperatureInput(), b.PressureInput(), b.MassRateInput(),
			b.GasOutput(), b.TemperatureOutput(), b.PressureOutput(), b.MassRateOutput(),
		},
		[]graph.Port{
			pBurn.GasInput(), pBurn.TemperatureInput(), pBurn.PressureInput(), pBurn.MassRateInput(),
			pBurn.GasOutput(), pBurn.TemperatureOutput(), pBurn.PressureOutput(), pBurn.MassRateOutput(),
		},
	)
	return pBurn
}

func GetMassRate(power float64, scheme schemes.Scheme, mrs nodes.MassRateSink) float64 {
	return schemes.GetMassRate(power, scheme) * mrs.MassRateInput().GetState().Value().(float64)
}
