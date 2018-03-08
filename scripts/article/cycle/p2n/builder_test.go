package p2n

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/core/math"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/library/parametric/free2n"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/stretchr/testify/suite"
	math2 "math"
	"testing"
)

type BuilderTestSuite struct {
	suite.Suite
	pScheme  free2n.DoubleShaftFreeScheme
	pNetwork graph.Network
	scheme   schemes.TwoShaftsScheme

	vSolver math.Solver
}

func (s *BuilderTestSuite) SetupTest() {
	s.scheme = GetScheme(piStag)

	var err error
	s.pScheme, err = GetParametric(s.scheme)
	s.Require().Nil(err)

	s.pNetwork, _ = s.pScheme.GetNetwork()

	sysCall := variator.SysCallFromNetwork(
		s.pNetwork, s.pScheme.Assembler().GetVectorPort(),
		relaxCoef, 2, iterNum, precision,
	)
	s.vSolver = variator.NewVariatorSolver(
		sysCall, s.pScheme.Variators(),
		newton.NewUniformNewtonSolverGen(1e-5, newton.DefaultLog),
	)

	//_, sErr := vSolver.Solve(vSolver.GetInit(), 1e-6, 1, 10000)
	//s.Require().Nil(sErr)
}

func (s *BuilderTestSuite) TestConsistency() {
	err := s.pNetwork.Solve(1, 2, 100, 1e-5)
	s.Require().Nil(err)

	initMassRate := schemes.GetMassRate(power, s.scheme)
	bOMR := s.scheme.Burner().MassRateOutput().GetState().Value().(float64)
	ctIMR := s.scheme.TurboCascade().Turbine().MassRateInput().GetState().Value().(float64)
	ctOMR := s.scheme.TurboCascade().Turbine().MassRateOutput().GetState().Value().(float64)
	ftIMR := s.scheme.FreeTurbineBlock().FreeTurbine().MassRateInput().GetState().Value().(float64)
	ftOMR := s.scheme.FreeTurbineBlock().FreeTurbine().MassRateOutput().GetState().Value().(float64)

	s.approxEqual(initMassRate, s.pScheme.Compressor().MassRateInput().GetState().Value().(float64), 3e-2)
	s.approxEqual(initMassRate, s.pScheme.Compressor().MassRateOutput().GetState().Value().(float64), 3e-2)
	s.InDelta(s.scheme.Compressor().PiStag(), s.pScheme.Compressor().PiStag(), 1e-6)
	s.approxEqual(
		s.scheme.Compressor().PowerOutput().GetState().Value().(float64),
		s.pScheme.Compressor().PowerOutput().GetState().Value().(float64),
		5e-2,
	)

	s.approxEqual(initMassRate*bOMR, s.pScheme.Burner().MassRateOutput().GetState().Value().(float64), 1e-2)
	s.InDelta(s.scheme.Burner().Alpha(), s.pScheme.Burner().Alpha(), 1e-3)

	s.approxEqual(
		initMassRate*ctIMR,
		s.pScheme.CompressorTurbine().MassRateInput().GetState().Value().(float64), 3e-2,
	)
	s.approxEqual(
		initMassRate*ctOMR,
		s.pScheme.CompressorTurbine().MassRateOutput().GetState().Value().(float64), 3e-2,
	)
	s.approxEqual(
		s.scheme.TurboCascade().Turbine().PowerOutput().GetState().Value().(float64),
		s.pScheme.CompressorTurbine().PowerOutput().GetState().Value().(float64),
		2e-2,
	)
	s.approxEqual(
		s.scheme.TurboCascade().Turbine().PiTStag(),
		s.pScheme.CompressorTurbine().PiTStag(),
		1e-4,
	)
	s.approxEqual(
		s.scheme.TurboCascade().Turbine().TStagOut(),
		s.pScheme.CompressorTurbine().TStagOut(),
		2e-2,
	)

	s.approxEqual(
		initMassRate*ftIMR,
		s.pScheme.FreeTurbine().MassRateInput().GetState().Value().(float64), 5e-2,
	)
	s.approxEqual(
		initMassRate*ftOMR,
		s.pScheme.FreeTurbine().MassRateOutput().GetState().Value().(float64), 5e-2,
	)
	s.approxEqual(
		s.scheme.FreeTurbineBlock().FreeTurbine().PowerOutput().GetState().Value().(float64),
		s.pScheme.FreeTurbine().PowerOutput().GetState().Value().(float64),
		4e-2,
	)
	s.approxEqual(
		s.scheme.FreeTurbineBlock().FreeTurbine().PiTStag(),
		s.pScheme.FreeTurbine().PiTStag(),
		1e-4,
	)
	s.approxEqual(
		s.scheme.FreeTurbineBlock().FreeTurbine().TStagOut(),
		s.pScheme.FreeTurbine().TStagOut(),
		2e-2,
	)
}

func (s *BuilderTestSuite) approxEqual(x1, x2, precision float64) {
	s.True(
		common.ApproxEqual(x1, x2, precision), "need %f got %f",
		precision,
		math2.Abs(x1-x2)/math2.Max(math2.Abs(x1), math2.Abs(x2)),
	)
}

func TestP2NBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(BuilderTestSuite))
}
