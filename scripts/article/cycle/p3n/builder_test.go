package p3n

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbocycle/core/math"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/stretchr/testify/suite"
	math2 "math"
	"testing"
	"github.com/Sovianum/turbocycle/library/parametric/free3n"
)

type BuilderTestSuite struct {
	suite.Suite
	pScheme  free3n.ThreeShaftFreeScheme
	pNetwork graph.Network
	scheme   schemes.ThreeShaftsScheme

	vSolver math.Solver
}

func (s *BuilderTestSuite) SetupTest() {
	s.scheme = GetScheme(lpcPiStag, hpcPiStag)

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
}

func (s *BuilderTestSuite) TestConsistency() {
	err := s.pNetwork.Solve(1, 2, 100, 1e-5)
	s.Require().Nil(err)

	initMassRate := schemes.GetMassRate(power, s.scheme)

	lpcIMR := s.scheme.LPC().MassRateInput().GetState().Value().(float64)
	lpcOMR := s.scheme.LPC().MassRateOutput().GetState().Value().(float64)

	hpcIMR := s.scheme.HPC().MassRateInput().GetState().Value().(float64)
	hpcOMR := s.scheme.HPC().MassRateOutput().GetState().Value().(float64)

	bOMR := s.scheme.MainBurner().MassRateOutput().GetState().Value().(float64)

	hptIMR := s.scheme.HPT().MassRateInput().GetState().Value().(float64)
	hptOMR := s.scheme.HPT().MassRateOutput().GetState().Value().(float64)

	lptIMR := s.scheme.LPT().MassRateInput().GetState().Value().(float64)
	lptOMR := s.scheme.LPT().MassRateOutput().GetState().Value().(float64)

	ftIMR := s.scheme.FT().MassRateInput().GetState().Value().(float64)
	ftOMR := s.scheme.FT().MassRateOutput().GetState().Value().(float64)

	s.approxEqual(initMassRate, lpcIMR, 3e-2)
	s.approxEqual(initMassRate, lpcOMR, 3e-2)
	s.InDelta(s.scheme.LPC().PiStag(), s.pScheme.LPC().PiStag(), 1e-6)
	s.approxEqual(
		s.scheme.LPC().PowerOutput().GetState().Value().(float64),
		s.pScheme.LPC().PowerOutput().GetState().Value().(float64),
		5e-2,
	)

	s.approxEqual(initMassRate, hpcIMR, 3e-2)
	s.approxEqual(initMassRate, hpcOMR, 3e-2)
	s.InDelta(s.scheme.HPC().PiStag(), s.pScheme.HPC().PiStag(), 1e-6)
	s.approxEqual(
		s.scheme.HPC().PowerOutput().GetState().Value().(float64),
		s.pScheme.HPC().PowerOutput().GetState().Value().(float64),
		5e-2,
	)

	s.approxEqual(initMassRate*bOMR, s.pScheme.Burner().MassRateOutput().GetState().Value().(float64), 1e-2)
	s.InDelta(s.scheme.MainBurner().Alpha(), s.pScheme.Burner().Alpha(), 1e-3)

	s.approxEqual(
		initMassRate*hptIMR,
		s.pScheme.HPT().MassRateInput().GetState().Value().(float64), 3e-2,
	)
	s.approxEqual(
		initMassRate*hptOMR,
		s.pScheme.HPT().MassRateOutput().GetState().Value().(float64), 3e-2,
	)
	s.approxEqual(
		s.scheme.HPT().PowerOutput().GetState().Value().(float64),
		s.pScheme.HPT().PowerOutput().GetState().Value().(float64),
		1e-2,
	)
	s.approxEqual(
		s.scheme.HPT().PiTStag(),
		s.pScheme.HPT().PiTStag(),
		1e-4,
	)
	s.approxEqual(
		s.scheme.HPT().TStagOut(),
		s.pScheme.HPT().TStagOut(),
		1e-2,
	)

	s.approxEqual(
		initMassRate*lptIMR,
		s.pScheme.LPT().MassRateInput().GetState().Value().(float64), 3e-2,
	)
	s.approxEqual(
		initMassRate*lptOMR,
		s.pScheme.LPT().MassRateOutput().GetState().Value().(float64), 3e-2,
	)
	s.approxEqual(
		s.scheme.LPT().PowerOutput().GetState().Value().(float64),
		s.pScheme.LPT().PowerOutput().GetState().Value().(float64),
		1e-2,
	)
	s.approxEqual(
		s.scheme.LPT().PiTStag(),
		s.pScheme.LPT().PiTStag(),
		1e-4,
	)
	s.approxEqual(
		s.scheme.LPT().TStagOut(),
		s.pScheme.LPT().TStagOut(),
		1e-2,
	)

	s.approxEqual(
		initMassRate*ftIMR,
		s.pScheme.FT().MassRateInput().GetState().Value().(float64), 3e-2,
	)
	s.approxEqual(
		initMassRate*ftOMR,
		s.pScheme.FT().MassRateOutput().GetState().Value().(float64), 3e-2,
	)
	s.approxEqual(
		s.scheme.FT().PowerOutput().GetState().Value().(float64),
		s.pScheme.FT().PowerOutput().GetState().Value().(float64),
		1e-2,
	)
	s.approxEqual(
		s.scheme.FT().PiTStag(),
		s.pScheme.FT().PiTStag(),
		1e-4,
	)
	s.approxEqual(
		s.scheme.FT().TStagOut(),
		s.pScheme.FT().TStagOut(),
		1e-2,
	)
}

func (s *BuilderTestSuite) approxEqual(x1, x2, precision float64) {
	s.True(
		common.ApproxEqual(x1, x2, precision), "need %f got %f",
		precision,
		math2.Abs(x1-x2)/math2.Max(math2.Abs(x1), math2.Abs(x2)),
	)
}

func TestBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(BuilderTestSuite))
}
