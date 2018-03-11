package p3nc

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
	pScheme  free3n.ThreeShaftCoolFreeScheme
	pNetwork graph.Network
	scheme   schemes.ThreeShaftsCoolerScheme

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
		0.1, 2, iterNum, schemePrecision,
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

	bOMR := s.scheme.MainBurner().MassRateOutput().GetState().Value().(float64)

	hptIMR := s.scheme.HPT().MassRateInput().GetState().Value().(float64)

	lptIMR := s.scheme.LPT().MassRateInput().GetState().Value().(float64)

	ftIMR := s.scheme.FT().MassRateInput().GetState().Value().(float64)

	s.InDelta(
		initMassRate,
		s.pScheme.LPC().MassRateInput().GetState().Value().(float64),
		1e-6,
	)
	s.InDelta(s.scheme.LPC().PStagIn(), s.pScheme.LPC().PStagIn(), 1e-6)
	s.InDelta(s.scheme.LPC().PStagOut(), s.pScheme.LPC().PStagOut(), 1e-6)
	s.InDelta(s.scheme.LPC().TStagIn(), s.pScheme.LPC().TStagIn(), 1e-6)
	s.InDelta(s.scheme.LPC().TStagOut(), s.pScheme.LPC().TStagOut(), 1)
	s.InDelta(s.scheme.LPC().PiStag(), s.pScheme.LPC().PiStag(), 1e-8)
	s.InDelta(s.scheme.LPC().Eta(), s.pScheme.LPC().Eta(), 1e-8)
	s.approxEqual(
		s.scheme.LPC().PowerOutput().GetState().Value().(float64),
		s.pScheme.LPC().PowerOutput().GetState().Value().(float64),
		1e-3,
	)

	s.InDelta(
		initMassRate,
		s.pScheme.HPC().MassRateInput().GetState().Value().(float64),
		5e-3,
	)
	s.InDelta(s.scheme.HPC().PStagIn(), s.pScheme.HPC().PStagIn(), 1e-6)
	s.InDelta(s.scheme.HPC().PStagOut(), s.pScheme.HPC().PStagOut(), 1e-6)
	s.InDelta(s.scheme.HPC().TStagIn(), s.pScheme.HPC().TStagIn(), 1)
	s.InDelta(s.scheme.HPC().TStagOut(), s.pScheme.HPC().TStagOut(), 1)
	s.InDelta(s.scheme.HPC().PiStag(), s.pScheme.HPC().PiStag(), 1e-8)
	s.InDelta(s.scheme.HPC().Eta(), s.pScheme.HPC().Eta(), 1e-8)
	s.approxEqual(
		s.scheme.HPC().PowerOutput().GetState().Value().(float64),
		s.pScheme.HPC().PowerOutput().GetState().Value().(float64),
		1e-3,
	)

	s.InDelta(s.scheme.MainBurner().PStagIn(), s.pScheme.Burner().PStagIn(), 1e-6)
	s.InDelta(s.scheme.MainBurner().PStagOut(), s.pScheme.Burner().PStagOut(), 1e-6)
	s.InDelta(s.scheme.MainBurner().TStagIn(), s.pScheme.Burner().TStagIn(), 1)
	s.InDelta(s.scheme.MainBurner().TStagOut(), s.pScheme.Burner().TStagOut(), 5e-1)
	s.approxEqual(initMassRate*bOMR, s.pScheme.Burner().MassRateOutput().GetState().Value().(float64), 1e-4)
	s.InDelta(s.scheme.MainBurner().Alpha(), s.pScheme.Burner().Alpha(), 1e-3)

	s.InDelta(s.scheme.HPT().PStagIn(), s.pScheme.HPT().PStagIn(), 1e-6)
	s.InDelta(s.scheme.HPT().PStagOut(), s.pScheme.HPT().PStagOut(), 1e-6)
	s.InDelta(s.scheme.HPT().TStagIn(), s.pScheme.HPT().TStagIn(), 1)
	s.InDelta(s.scheme.HPT().TStagOut(), s.pScheme.HPT().TStagOut(), 1)
	s.approxEqual(
		initMassRate*hptIMR,
		s.pScheme.HPT().MassRateInput().GetState().Value().(float64), 1e-4,
	)
	s.InDelta(s.scheme.HPT().Eta(), s.pScheme.HPT().Eta(), 1e-6)

	s.InDelta(s.scheme.LPT().PStagIn(), s.pScheme.LPT().PStagIn(), 1e-6)
	s.InDelta(s.scheme.LPT().PStagOut(), s.pScheme.LPT().PStagOut(), 1e-6)
	s.InDelta(s.scheme.LPT().TStagIn(), s.pScheme.LPT().TStagIn(), 1e-1)
	s.InDelta(s.scheme.LPT().TStagOut(), s.pScheme.LPT().TStagOut(), 1e-1)
	s.approxEqual(
		initMassRate*lptIMR,
		s.pScheme.LPT().MassRateInput().GetState().Value().(float64), 5e-5,
	)
	s.InDelta(s.scheme.LPT().Eta(), s.pScheme.LPT().Eta(), 1e-6)

	s.InDelta(s.scheme.FT().PStagIn(), s.pScheme.FT().PStagIn(), 1e-6)
	s.InDelta(s.scheme.FT().PStagOut(), s.pScheme.FT().PStagOut(), 1e-6)
	s.InDelta(s.scheme.FT().TStagIn(), s.pScheme.FT().TStagIn(), 1e-1)
	s.InDelta(s.scheme.FT().TStagOut(), s.pScheme.FT().TStagOut(), 1e-1)
	s.approxEqual(
		initMassRate*ftIMR,
		s.pScheme.FT().MassRateInput().GetState().Value().(float64), 3e-5,
	)
	s.InDelta(s.scheme.FT().Eta(), s.pScheme.FT().Eta(), 1e-6)

	s.approxEqual(
		s.scheme.FT().PowerOutput().GetState().Value().(float64),
		s.pScheme.FT().PowerOutput().GetState().Value().(float64),
		6e-5,
	)

	s.approxEqual(
		-s.scheme.FT().PowerOutput().GetState().Value().(float64) * initMassRate*ftIMR,
		s.pScheme.Payload().PowerOutput().GetState().Value().(float64),
		1e-6,
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
