package p3n

import (
	"github.com/Sovianum/cooling-course-project/scripts/article/cycle/common"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/library/parametric/free3n"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/methodics"
)

func NewBuilder(
	source schemes.ThreeShaftsScheme,
	power,
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
		Source: source,
		Power:  power,

		LPCRpm0: lpcRpm0,
		HPCRpm0: hpcRpm0,

		LambdaIn0: lambdaIn0,

		LPTInletMeanDiameter: lptInletMeanDiameter,
		LPTLambdaU0:          lptLambdaU0,
		LPTStageNum:          lptStageNum,

		HPTInletMeanDiameter: hptInletMeanDiameter,
		HPTLambdaU0:          hptLambdaU0,
		HPTStageNum:          hptStageNum,

		FTInletMeanDiameter: ftInletMeanDiameter,
		FTLambdaU0:          ftLambdaU0,
		FTStageNum:          ftStageNum,

		PayloadRpm0: payloadRpm0,

		LPEtaM: lpEtaM,
		HPEtaM: hpEtaM,

		Precision: precision,
		RelaxCoef: relaxCoef,
		IterLimit: iterLimit,
	}
}

type Builder struct {
	Source schemes.ThreeShaftsScheme
	Power  float64

	LPCRpm0 float64
	HPCRpm0 float64

	LambdaIn0 float64

	LPTInletMeanDiameter float64
	LPTLambdaU0          float64
	LPTStageNum          float64

	HPTInletMeanDiameter float64
	HPTLambdaU0          float64
	HPTStageNum          float64

	FTInletMeanDiameter float64
	FTLambdaU0          float64
	FTStageNum          float64

	PayloadRpm0 float64

	LPEtaM float64
	HPEtaM float64

	Precision float64
	RelaxCoef float64
	IterLimit int
}

func (b *Builder) Build() free3n.ThreeShaftFreeScheme {
	return free3n.NewThreeShaftFreeScheme(
		b.Source.GasSource().GasOutput().GetState().Value().(gases.Gas),
		b.Source.LPC().TemperatureInput().GetState().Value().(float64),
		b.Source.LPC().PressureInput().GetState().Value().(float64),
		b.Source.MainBurner().TStagOut(),

		b.BuildLPC(), b.BuildLPCPipe(),
		b.BuildLPT(), b.BuildLPTPipe(),
		b.LPEtaM,

		b.BuildHPC(), b.BuildHPCPipe(),
		b.BuildHPT(), b.BuildHPTPipe(),
		b.HPEtaM,

		b.BuildFT(), b.BuildFTPipe(),
		b.BuildBurner(), b.BuildPayload(),
	)
}

func (b *Builder) BuildLPC() constructive.ParametricCompressorNode {
	c := b.Source.LPC()
	massRate0 := common.GetMassRate(b.Power, b.Source, c)
	charGen := methodics.NewCompressorCharGen(
		c.PiStag(), c.Eta(), massRate0, precision, relaxCoef, b.IterLimit,
	)

	return constructive.NewParametricCompressorNodeFromProto(
		c,
		charGen.GetNormEtaChar(), charGen.GetNormRPMChar(),
		b.LPCRpm0, common.GetMassRate(b.Power, b.Source, b.Source.LPC()),
		b.Precision,
	)
}

func (b *Builder) BuildLPCPipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(b.Source.LPCPipe().Sigma())
}

func (b *Builder) BuildHPC() constructive.ParametricCompressorNode {
	c := b.Source.HPC()
	massRate0 := common.GetMassRate(b.Power, b.Source, c)
	charGen := methodics.NewCompressorCharGen(
		c.PiStag(), c.Eta(), massRate0, precision, relaxCoef, b.IterLimit,
	)

	return constructive.NewParametricCompressorNodeFromProto(
		c,
		charGen.GetNormEtaChar(), charGen.GetNormRPMChar(),
		b.LPCRpm0, common.GetMassRate(b.Power, b.Source, b.Source.LPC()),
		b.Precision,
	)
}

func (b *Builder) BuildHPCPipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(b.Source.HPCPipe().Sigma())
}

func (b *Builder) BuildBurner() constructive.ParametricBurnerNode {
	burn := b.Source.MainBurner()
	return constructive.NewParametricBurnerFromProto(
		burn, b.LambdaIn0,
		common.GetMassRate(b.Power, b.Source, burn),
		b.Precision, b.RelaxCoef, b.IterLimit,
	)
}

func (b *Builder) BuildHPT() constructive.ParametricTurbineNode {
	char := methodics.NewKazandjanTurbineCharacteristic()
	return constructive.NewParametricTurbineNodeFromProto(
		b.Source.HPT(),
		char.GetNormMassRateChar(), char.GetNormEtaChar(),
		common.GetMassRate(b.Power, b.Source, b.Source.HPT()),
		b.HPTInletMeanDiameter, b.Precision,
	)
}

func (b *Builder) BuildHPTPipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(b.Source.HPTPipe().Sigma())
}

func (b *Builder) BuildLPT() constructive.ParametricTurbineNode {
	char := methodics.NewKazandjanTurbineCharacteristic()
	return constructive.NewParametricTurbineNodeFromProto(
		b.Source.LPT(),
		char.GetNormMassRateChar(), char.GetNormEtaChar(),
		common.GetMassRate(b.Power, b.Source, b.Source.LPT()),
		b.LPTInletMeanDiameter, b.Precision,
	)
}

func (b *Builder) BuildLPTPipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(b.Source.LPTPipe().Sigma())
}

func (b *Builder) BuildFT() constructive.ParametricTurbineNode {
	char := methodics.NewKazandjanTurbineCharacteristic()
	return constructive.NewParametricTurbineNodeFromProto(
		b.Source.FT(),
		char.GetNormMassRateChar(), char.GetNormEtaChar(),
		common.GetMassRate(b.Power, b.Source, b.Source.FT()),
		b.FTInletMeanDiameter, b.Precision,
	)
}

func (b *Builder) BuildFTPipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(b.Source.FTBlock().OutletPressureLoss().Sigma())
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
