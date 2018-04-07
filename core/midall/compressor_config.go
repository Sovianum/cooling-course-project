package midall

import (
	"fmt"
	"github.com/Sovianum/turbocycle/core/math"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/stage/common"
	"github.com/Sovianum/turbocycle/impl/stage/compressor"
	"github.com/Sovianum/turbocycle/impl/stage/ditributions"
	"gonum.org/v1/gonum/mat"
)

type CompressorConfig struct {
	// global parameters
	StageNum int     // количество ступеней
	RPM      float64 // частота вращения ротора
	MassRate float64 // расход воздуха через компрессор

	// geometric parameters
	DRelIn              float64   // относительный диаметр втулки на входе в компрессор
	RotorElongationArr  []float64 // удлинение лопаток ротора
	DeltaRotorRelArr    []float64 // относительный зазор за лопатками ротора
	StatorElongationArr []float64 // удлинение лопаток статора
	DeltaStatorRelArr   []float64 // относительный зазор за лопатками статора
	GammaInArr          []float64 // внутренний угол раскрытия проточной части
	GammaOutArr         []float64 // внешний угол раскрытия проточной части

	// thermodynamic parameters
	// принимается бипараболическая форма изменения коэффициента нагрузки
	HtLossStart float64 // уменьшение коэффициента напора ко входу компрессора
	HtLossEnd   float64 // уменьшение коэффициента напора к выходу компрессора
	HtMax       float64 // максимальное значение коэффициента нагрузки
	HtMaxCoord  float64 // координата максимального значения коэффициента нагрузки
	HtLimit     float64 // предельное значение, до которого оптимизатор может поднимать коэффициент нагрузки ступеней

	// принимается бипараболическая форма изменения КПД
	EtaLossStart float64 // уменьшение КПД по входу компрессора
	EtaLossEnd   float64 // уменьшение КПД к выходу компрессора
	EtaMax       float64 // максимальное значение КПД компрессора
	EtaMaxCoord  float64 // координата максимального значения КПД компрессора
	EtaLimit     float64 // предельное значение, до которого оптимизатор может поднимать КПД ступени

	// принимается линейная форма изменения степени реактивности ступеней
	ReactivityStart float64
	ReactivityEnd   float64

	// принимается линейная форма изменения коэффициента расхода по ступеням
	CaStart float64
	CaEnd   float64

	// принимается постоянное значение коэффициента работы
	LabourCoef float64

	// numerical parameters
	Precision  float64
	RelaxCoef  float64
	InitLambda float64
	IterLimit  int
}

func (conf *CompressorConfig) GetFittedStagedCompressor(
	cycleNode constructive.CompressorNode,
	solverGen math.SolverGenerator,
) (compressor.StagedCompressorNode, error) {
	if err := cycleNode.Process(); err != nil {
		return nil, fmt.Errorf("failed to process cycleNode %s", err.Error())
	}
	staged, err := conf.GetStagedCompressor()
	if err != nil {
		return nil, err
	}
	staged.MassRateInput().SetState(states.NewMassRatePortState(conf.MassRate))

	eqSys := compressor.GetCycleFitEqSys(
		staged, cycleNode,
		common.Scaler(conf.getHtNormFunc()),
		common.Scaler(conf.getEtaFunc()),
		conf.HtLimit, conf.EtaLimit,
	)

	solver, err := solverGen(eqSys)
	if err != nil {
		return nil, fmt.Errorf("failed to create solver %s", err.Error())
	}

	_, err = solver.Solve(
		mat.NewVecDense(2, []float64{1, 1}),
		conf.Precision,
		1, 1000,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fit to cycle: %s", err.Error())
	}

	return staged, nil
}

func (conf *CompressorConfig) GetStagedCompressor() (compressor.StagedCompressorNode, error) {
	if err := conf.validate(); err != nil {
		return nil, err
	}
	geomList := make([]compressor.IncompleteStageGeometryGenerator, conf.StageNum)
	for i := 0; i != conf.StageNum; i++ {
		geomList[i] = compressor.NewIncompleteStageGeomGen(
			compressor.NewIncompleteGenerator(
				conf.RotorElongationArr[i], conf.DeltaRotorRelArr[i], conf.GammaInArr[i], conf.GammaOutArr[i],
			),
			compressor.NewIncompleteGenerator(
				conf.StatorElongationArr[i], conf.DeltaStatorRelArr[i], conf.GammaInArr[i], conf.GammaOutArr[i],
			),
		)
	}
	xEnd := float64(conf.StageNum - 1)
	htLaw := conf.getHtNormFunc()
	reactivityLaw := ditributions.GetLinear(0, conf.ReactivityStart, xEnd, conf.ReactivityEnd)
	labourCoefLaw := ditributions.GetConstant(conf.LabourCoef)
	etaAdLaw := conf.getEtaFunc()
	caCoefLaw := ditributions.GetLinear(0, conf.CaStart, xEnd, conf.CaEnd)

	return compressor.NewStagedCompressorNode(
		conf.RPM, conf.DRelIn, geomList,
		common.FromDistribution(htLaw),
		common.FromDistribution(reactivityLaw),
		common.FromDistribution(labourCoefLaw),
		common.FromDistribution(etaAdLaw),
		common.FromDistribution(caCoefLaw),
		conf.Precision, conf.RelaxCoef, conf.InitLambda, conf.IterLimit,
	), nil
}

func (conf *CompressorConfig) validate() error {
	if len(conf.RotorElongationArr) < conf.StageNum {
		return fmt.Errorf("invalid number of rotorElongationArr")
	}
	if len(conf.DeltaRotorRelArr) < conf.StageNum {
		return fmt.Errorf("invalid number of deltaRotorRelArr")
	}
	if len(conf.StatorElongationArr) < conf.StageNum {
		return fmt.Errorf("invalid number of statorElongationArr")
	}
	if len(conf.DeltaStatorRelArr) < conf.StageNum {
		return fmt.Errorf("invalid number of deltaStatorRelArr")
	}
	if len(conf.GammaInArr) < conf.StageNum {
		return fmt.Errorf("invalid number of gammaInArr")
	}
	if len(conf.GammaOutArr) < conf.StageNum {
		return fmt.Errorf("invalid number of gammaOutArr")
	}
	return nil
}

func (conf *CompressorConfig) getEtaFunc() common.Func1D {
	xEnd := float64(conf.StageNum - 1)
	return ditributions.GetUnitBiParabolic(0, xEnd, conf.EtaMaxCoord, conf.EtaLossStart, conf.EtaLossEnd).
		Scale(conf.EtaMax)
}

func (conf *CompressorConfig) getHtNormFunc() common.Func1D {
	xEnd := float64(conf.StageNum - 1)
	return ditributions.GetUnitBiParabolic(
		0, xEnd, conf.HtMaxCoord, conf.HtLossStart, conf.HtLossEnd,
	).Scale(conf.HtMax)
}
