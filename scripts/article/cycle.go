package article

import (
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/cooling-course-project/core/schemes/two_shafts"
	"github.com/Sovianum/turbocycle/library/parametric/free2n"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"fmt"
)

const (
	etaM = 0.99
	cRpm0 = 10000
	cLambdaIn0 = 0.3

	ctID = 0.3
	ctLambdaU0 = 0.3
	ctStageNum = 1

	ftID = 0.5
	ftLambdaU0 = 0.3
	ftStageNum = 1

	payloadRpm0 = 3000

	t0 = 300
	p0 = 1e5
)

func solveParametric(pScheme free2n.DoubleShaftFreeScheme) error {
	network, pErr := pScheme.GetNetwork()
	if pErr != nil {
		return pErr
	}

	sysCall := variator.SysCallFromNetwork(
		network, pScheme.Assembler().GetVectorPort(),
		relaxCoef, 2, iterNum, precision,
	)
	vSolver := variator.NewVariatorSolver(
		sysCall, pScheme.Variators(),
		newton.NewUniformNewtonSolverGen(1e-5, newton.DefaultLog),
	)

	_, sErr := vSolver.Solve(vSolver.GetInit(), 1e-6, 0.5, 10000)
	if sErr != nil {
		return sErr
	}

	return nil
}

func getParametric(scheme *schemes.TwoShaftsSchemeImpl) (free2n.DoubleShaftFreeScheme, error) {
	network, err := scheme.GetNetwork()
	if err != nil {
		return nil, err
	}
	converged, solveErr := network.Solve(relaxCoef, 2, iterNum, precision)
	if solveErr != nil {
		return nil, solveErr
	}
	if !converged {
		return nil, fmt.Errorf("failed to converge")
	}

	return get2nParametricScheme(scheme), nil
}

func get2nParametricScheme(scheme schemes.TwoShaftsScheme) free2n.DoubleShaftFreeScheme {
	casted := scheme.(*schemes.TwoShaftsSchemeImpl)
	builder := NewParametric2NBuilder(
		casted, power, t0, p0, cRpm0, cLambdaIn0,
		ctID, ctLambdaU0, ctStageNum,
		ftID, ftLambdaU0, ftStageNum,
		payloadRpm0, etaM, precision, relaxCoef, iterNum,
	)
	return builder.Build()
}

func get2nScheme(piStag float64) schemes.TwoShaftsScheme {
	scheme := two_shafts.GetInitedTwoShaftsScheme()
	scheme.Compressor().SetPiStag(piStag)
	return scheme
}
