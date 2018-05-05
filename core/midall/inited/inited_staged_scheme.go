package inited

import (
	"fmt"
	"github.com/Sovianum/cooling-course-project/core/midall"
	"github.com/Sovianum/cooling-course-project/core/schemes/s3n"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/library/schemes"
	"math"
)

const (
	RPMHigh = 12e3
	RPMLow  = 9e3
	RPMFree = 7.8e3

	precision  = 1e-3
	relaxCoef  = 0.1
	initLambda = 0.3
	iterLimit  = 1000

	power = 16e6 / 0.98
)

func GetInitedStagedNodes() (*midall.StagedScheme3n, error) {
	source := s3n.GetDiplomaInitedThreeShaftsScheme()
	network, err := source.GetNetwork()
	if err != nil {
		return nil, err
	}
	if err := network.Solve(relaxCoef, 2, 100, precision); err != nil {
		return nil, err
	}
	fmt.Printf(
		"pi_LPC = %.3f, pi_HPC = %.3f, pi_HPT = %.3f, pi_LPT = %.3f, pi_FT = %.3f\n",
		source.LPC().PiStag(), source.HPC().PiStag(),
		source.HPT().PiTStag(), source.LPT().PiTStag(),
		source.FT().PiTStag(),
	)

	massRate := schemes.GetMassRate(power, source)
	lpcMassRate := massRate * source.LPC().MassRateInput().GetState().Value().(float64)
	hpcMassRate := massRate * source.HPC().MassRateInput().GetState().Value().(float64)
	hptMassRate := massRate * source.HPT().MassRateInput().GetState().Value().(float64)
	lptMassRate := massRate * source.LPT().MassRateInput().GetState().Value().(float64)
	ftMassRate := massRate * source.FT().MassRateInput().GetState().Value().(float64)

	lpcConfig := getLPCConfig()
	lpcConfig.MassRate = lpcMassRate

	hpcConfig := getHPCConfig()
	hpcConfig.MassRate = hpcMassRate

	hptConfig := getHPTConfig()
	hptConfig.MassRate = hptMassRate

	lptConfig := getLPTConfig()
	lptConfig.MassRate = lptMassRate

	ftConfig := getFTConfig()
	ftConfig.MassRate = ftMassRate

	return midall.NewStagedScheme3n(source, lpcConfig, hpcConfig, hptConfig, lptConfig, ftConfig)
}

func getLPCConfig() midall.CompressorConfig {
	return midall.CompressorConfig{
		StageNum: 7,
		RPM:      RPMLow,

		DRelIn: 0.43,

		RotorElongationArr: []float64{3.65, 3.56, 3.14, 3.17, 3.1, 3.16, 3.04},
		DeltaRotorRelArr:   []float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1},

		StatorElongationArr: []float64{4.36, 3.75, 3.0, 3.00, 2.85, 3.13, 2.95},
		DeltaStatorRelArr:   []float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1},

		GammaInArr: []float64{
			common.ToRadians(8),
			common.ToRadians(8),
			common.ToRadians(8),
			common.ToRadians(5),
			common.ToRadians(5),
			common.ToRadians(5),
			common.ToRadians(5),
		},
		GammaOutArr: []float64{
			common.ToRadians(-8),
			common.ToRadians(-8),
			common.ToRadians(-8),
			common.ToRadians(-5),
			common.ToRadians(-5),
			common.ToRadians(-4),
			common.ToRadians(-4),
		},

		HtLossStart: 0.1,
		HtLossEnd:   0.01,
		HtMax:       0.2,
		HtMaxCoord:  3,
		HtLimit:     0.5,

		EtaLossStart: 0.02,
		EtaLossEnd:   0.02,
		EtaMax:       0.82,
		EtaMaxCoord:  1,
		EtaLimit:     0.9,

		ReactivityStart: 0.5,
		ReactivityEnd:   0.5,
		HasPreTwist:     true,

		CaStart: 0.5,
		CaEnd:   0.4,

		LabourCoef: 0.99,

		Precision:  precision,
		RelaxCoef:  relaxCoef,
		InitLambda: 1,
		IterLimit:  iterLimit,
	}
}

func getHPCConfig() midall.CompressorConfig {
	return midall.CompressorConfig{
		StageNum: 5,
		RPM:      RPMHigh,

		DRelIn: 0.85,

		RotorElongationArr: []float64{2.5, 2.5, 2.5, 2.25, 2.25, 2, 3, 3},
		DeltaRotorRelArr:   []float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1},

		StatorElongationArr: []float64{2.5, 2.5, 2.5, 2.25, 2.25, 2, 3, 3},
		DeltaStatorRelArr:   []float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1},

		GammaInArr: []float64{
			common.ToRadians(9),
			common.ToRadians(9),
			common.ToRadians(9),
			common.ToRadians(8),
			common.ToRadians(7),
			common.ToRadians(6),
			common.ToRadians(12),
			common.ToRadians(3),
		},
		GammaOutArr: []float64{0, 0, 0, 0, 0, 0, 0, 0},

		HtLossStart: 0.02,
		HtLossEnd:   0.02,
		HtMax:       0.32,
		HtMaxCoord:  2,
		HtLimit:     0.5,

		EtaLossStart: 0.02,
		EtaLossEnd:   0.02,
		EtaMax:       0.82,
		EtaMaxCoord:  2,
		EtaLimit:     0.9,

		ReactivityStart: 0.5,
		ReactivityEnd:   0.5,
		HasPreTwist:     false,

		CaStart: 0.4,
		CaEnd:   0.5,

		LabourCoef: 0.99,

		Precision:  precision,
		RelaxCoef:  relaxCoef,
		InitLambda: initLambda,
		IterLimit:  iterLimit,
	}
}

func getHPTConfig() midall.TurbineConfig {
	return midall.TurbineConfig{
		StageNum: 1,
		RPM:      RPMHigh,
		Alpha1:   common.ToRadians(12),

		TotalHeatDrop: math.NaN(), // heat drop will be set wile fitting

		LRelIn: 0.09,

		StatorElongationArr: []float64{1.3, 1.3},
		DeltaStatorRelArr:   []float64{0.1, 0.1},
		ApproxTRotorRel:     []float64{0.7, 0.7},

		RotorElongationArr: []float64{1.75, 1.75},
		DeltaRotorRelArr:   []float64{0.1, 0.1},
		ApproxTStatorRel:   []float64{0.7, 0.7},

		GammaInArr: []float64{
			common.ToRadians(8),
			common.ToRadians(-3),
		},
		GammaOutArr: []float64{
			common.ToRadians(20),
			common.ToRadians(3),
		},

		PhiStartLoss: 0, PhiEndLoss: 0, PhiMax: 0.97, PhiMaxCoord: 0,
		PsiStartLoss: 0, PsiEndLoss: 0, PsiMax: 0.97, PsiMaxCoord: 0,

		HtStartLoss: 0, HtEndLoss: 0.1, HtMaxCoord: 0,

		ReactivityStart: 0.3, ReactivityEnd: 0.3,
		AirGapRelStart: 0.001, AirGapRelEnd: 0.001,

		Precision: precision,
	}
}

func getLPTConfig() midall.TurbineConfig {
	return midall.TurbineConfig{
		StageNum: 1,
		RPM:      RPMLow,
		Alpha1:   common.ToRadians(16),

		TotalHeatDrop: math.NaN(), // heat drop will be set wile fitting

		LRelIn: 0.13,

		StatorElongationArr: []float64{2.5, 2.5},
		DeltaStatorRelArr:   []float64{0.1, 0.1},
		ApproxTRotorRel:     []float64{0.7, 0.7},

		RotorElongationArr: []float64{2.7, 2.7},
		DeltaRotorRelArr:   []float64{0.1, 0.1},
		ApproxTStatorRel:   []float64{0.7, 0.7},

		GammaInArr: []float64{
			common.ToRadians(8),
			common.ToRadians(0),
		},
		GammaOutArr: []float64{
			common.ToRadians(20),
			common.ToRadians(10),
		},

		PhiStartLoss: 0, PhiEndLoss: 0, PhiMax: 0.97, PhiMaxCoord: 0,
		PsiStartLoss: 0, PsiEndLoss: 0, PsiMax: 0.97, PsiMaxCoord: 0,

		HtStartLoss: 0, HtEndLoss: 0.07, HtMaxCoord: 0,

		ReactivityStart: 0.3, ReactivityEnd: 0.5,
		AirGapRelStart: 0.001, AirGapRelEnd: 0.001,

		Precision: precision,
	}
}

func getFTConfig() midall.TurbineConfig {
	return midall.TurbineConfig{
		StageNum: 2,
		RPM:      RPMFree,
		Alpha1:   common.ToRadians(13),

		TotalHeatDrop: math.NaN(), // heat drop will be set wile fitting

		LRelIn: 0.245,

		StatorElongationArr: []float64{3.5, 3.5, 3.5},
		DeltaStatorRelArr:   []float64{0.1, 0.1, 0.1},
		ApproxTRotorRel:     []float64{0.7, 0.7, 0.7},

		RotorElongationArr: []float64{4, 4, 4},
		DeltaRotorRelArr:   []float64{0.1, 0.1, 0.1},
		ApproxTStatorRel:   []float64{0.7, 0.7, 0.7},

		GammaInArr: []float64{
			common.ToRadians(0),
			common.ToRadians(0),
			common.ToRadians(0),
		},
		GammaOutArr: []float64{
			common.ToRadians(10),
			common.ToRadians(10),
			common.ToRadians(10),
		},

		PhiStartLoss: 0, PhiEndLoss: 0, PhiMax: 0.97, PhiMaxCoord: 0,
		PsiStartLoss: 0, PsiEndLoss: 0, PsiMax: 0.97, PsiMaxCoord: 0,

		HtStartLoss: 0.1, HtEndLoss: 0, HtMaxCoord: 1,

		ReactivityStart: 0.42, ReactivityEnd: 0.42,
		AirGapRelStart: 0.001, AirGapRelEnd: 0.001,

		Precision: precision,
	}
}
