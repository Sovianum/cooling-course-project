package midall

import (
	"fmt"
	"github.com/Sovianum/turbocycle/core/math"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/stage/common"
	"github.com/Sovianum/turbocycle/impl/stage/ditributions"
	"github.com/Sovianum/turbocycle/impl/stage/turbine"
	"gonum.org/v1/gonum/mat"
)

type TurbineConfig struct {
	// global parameters
	StageNum      int     // количество ступеней
	RPM           float64 // частота вращения ротора
	MassRate      float64 // расход воздуха через компрессор
	Alpha1        float64 // направление потока за статором первой ступени
	TotalHeatDrop float64 // полный теплоперепад на ступени

	// geometric parameters
	LRelIn float64 // относительная длина лопатки на входе в турбину

	StatorElongationArr []float64 // удлинение лопаток статора
	DeltaStatorRelArr   []float64 // относительный зазор за лопатками статора
	ApproxTStatorRel    []float64 // примерный относительный шаг лопаток статора

	RotorElongationArr []float64 // удлинение лопаток ротора
	DeltaRotorRelArr   []float64 // относительный зазор за лопатками ротора
	ApproxTRotorRel    []float64 // примерный относительный шаг лопаток ротора

	GammaInArr  []float64 // внутренний угол раскрытия проточной части
	GammaOutArr []float64 // внешний угол раскрытия проточной части

	// thermodynamic parameters
	// принимается бипараболическая форма распределения параметра phi
	PhiStartLoss float64 // уменьшение phi ко входу турбины
	PhiEndLoss   float64 // уменьшение phi к выходу из турбины
	PhiMax       float64 // наибольшее значение phi
	PhiMaxCoord  float64 // коордианат максимального значения phi

	// принимается бипараболическая форма распределения параметра psi
	PsiStartLoss float64 // уменьшение psi ко входу турбины
	PsiEndLoss   float64 // умешьнешие psi к выходу из турбины
	PsiMax       float64 // наибольшее значение psi
	PsiMaxCoord  float64 // коордианата максимального значения psi

	// принимается бипараболическая форма распределения теплоперепада по ступеням
	HtStartLoss float64 // уменьшение теплопереда ко входу турбины (обычно отрицательное)
	HtEndLoss   float64 // уменьшение теплопереаада к выходу из турбины
	HtMaxCoord  float64 // координата максимального значения теплоперепада

	// принимается линейная форма изменения степени реактивности ступеней
	ReactivityStart float64
	ReactivityEnd   float64

	// принимается линейная форма изменения относительного радиального зазора
	AirGapRelStart float64
	AirGapRelEnd   float64

	// numerical parameters
	Precision float64
}

func (conf *TurbineConfig) GetFittedStagedTurbine(
	cycleNode constructive.StaticTurbineNode,
	solverGen math.SolverGenerator,
) (turbine.StagedTurbineNode, error) {
	if err := cycleNode.Process(); err != nil {
		return nil, fmt.Errorf("failed to process cycleNode %s", err.Error())
	}
	staged, err := conf.GetStagedTurbine()
	if err != nil {
		return nil, err
	}
	staged.MassRateInput().SetState(states.NewMassRatePortState(conf.MassRate))

	eqSys := turbine.GetCycleFitEqSys(
		staged, cycleNode,
		common.Scaler(staged.GetPhiFunc()),
		common.Scaler(staged.GetPsiFunc()),
	)

	solver, err := solverGen(eqSys)
	if err != nil {
		return nil, fmt.Errorf("failed to create solver %s", err.Error())
	}

	_, err = solver.Solve(
		mat.NewVecDense(2, []float64{staged.Ht(), 1}),
		conf.Precision,
		1, 1000,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fit to cycle: %s", err.Error())
	}

	return staged, nil
}

func (conf *TurbineConfig) GetStagedTurbine() (turbine.StagedTurbineNode, error) {
	if err := conf.validate(); err != nil {
		return nil, err
	}
	geomList := make([]turbine.IncompleteStageGeometryGenerator, conf.StageNum)
	for i := 0; i != conf.StageNum; i++ {
		geomList[i] = turbine.NewIncompleteStageGeometryGenerator(
			turbine.NewIncompleteGenerator(
				conf.StatorElongationArr[i], conf.DeltaStatorRelArr[i],
				conf.GammaInArr[i], conf.GammaOutArr[i], conf.ApproxTStatorRel[i],
			),
			turbine.NewIncompleteGenerator(
				conf.RotorElongationArr[i], conf.DeltaRotorRelArr[i],
				conf.GammaInArr[i], conf.GammaOutArr[i], conf.ApproxTRotorRel[i],
			),
		)
	}
	xEnd := float64(conf.StageNum - 1)

	var phiLaw common.Func1D
	if conf.StageNum > 1 {
		phiLaw = ditributions.GetUnitBiParabolic(
			0, xEnd, conf.PhiMaxCoord, conf.PhiStartLoss, conf.PhiEndLoss,
		).Scale(conf.PhiMax)
	} else {
		phiLaw = ditributions.GetConstant(conf.PhiMax)
	}

	var psiLaw common.Func1D
	if conf.StageNum > 1 {
		psiLaw = ditributions.GetUnitBiParabolic(
			0, xEnd, conf.PsiMaxCoord, conf.PsiStartLoss, conf.PsiEndLoss,
		).Scale(conf.PsiMax)
	} else {
		psiLaw = ditributions.GetConstant(conf.PsiMax)
	}

	var reactivityLaw common.Func1D
	if conf.StageNum > 1 {
		reactivityLaw = ditributions.GetLinear(0, conf.ReactivityStart, xEnd, conf.ReactivityEnd)
	} else {
		reactivityLaw = ditributions.GetConstant(conf.ReactivityStart)
	}

	var airGapRelLaw common.Func1D
	if conf.StageNum > 1 {
		airGapRelLaw = ditributions.GetLinear(0, conf.AirGapRelStart, xEnd, conf.AirGapRelEnd)
	} else {
		airGapRelLaw = ditributions.GetConstant(conf.AirGapRelStart)
	}

	heatDropDistributionFunc := ditributions.GetUnitBiParabolic(
		0, xEnd, conf.HtMaxCoord, conf.HtStartLoss, conf.HtEndLoss,
	)

	return turbine.NewStagedTurbineNode(
		conf.RPM, conf.Alpha1, conf.TotalHeatDrop,
		conf.LRelIn, phiLaw, psiLaw, reactivityLaw, airGapRelLaw,
		heatDropDistributionFunc, geomList, conf.Precision,
	), nil
}

func (conf *TurbineConfig) validate() error {
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
	if len(conf.ApproxTRotorRel) < conf.StageNum {
		return fmt.Errorf("invalid number of approxTRotorRel")
	}
	if len(conf.ApproxTStatorRel) < conf.StageNum {
		return fmt.Errorf("invalid number of approxTStatorRel")
	}
	return nil
}
