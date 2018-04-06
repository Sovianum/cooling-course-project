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
}

func (s *BuilderTestSuite) TestConsistencySingleIter() {
	n, _ := s.scheme.GetNetwork()
	n.Solve(1, 2, 100, 1e-5)

	s.Require().Nil(s.pNetwork.Solve(0.1, 2, 1000, 0.1))

	initMassRate := schemes.GetMassRate(power, s.scheme)
	bOMR := s.scheme.Burner().MassRateOutput().GetState().Value().(float64)
	ctIMR := s.scheme.TurboCascade().Turbine().MassRateInput().GetState().Value().(float64)
	ftIMR := s.scheme.FreeTurbineBlock().FreeTurbine().MassRateInput().GetState().Value().(float64)

	s.InDelta(
		initMassRate,
		s.pScheme.Compressor().MassRateInput().GetState().Value().(float64),
		1e-6,
	)
	s.InDelta(s.scheme.Compressor().PStagIn(), s.pScheme.Compressor().PStagIn(), 1e-6)
	s.InDelta(s.scheme.Compressor().PStagOut(), s.pScheme.Compressor().PStagOut(), 1e-6)
	s.InDelta(s.scheme.Compressor().TStagIn(), s.pScheme.Compressor().TStagIn(), 1e-6)
	s.InDelta(s.scheme.Compressor().TStagOut(), s.pScheme.Compressor().TStagOut(), 0.1)
	s.InDelta(s.scheme.Compressor().PiStag(), s.pScheme.Compressor().PiStag(), 1e-8)
	s.InDelta(s.scheme.Compressor().Eta(), s.pScheme.Compressor().Eta(), 1e-8)
	s.approxEqual(
		s.scheme.Compressor().PowerOutput().GetState().Value().(float64),
		s.pScheme.Compressor().PowerOutput().GetState().Value().(float64),
		1e-3,
	)

	s.InDelta(s.scheme.Burner().PStagIn(), s.pScheme.Burner().PStagIn(), 1e-6)
	s.InDelta(s.scheme.Burner().PStagOut(), s.pScheme.Burner().PStagOut(), 1e-6)
	s.InDelta(s.scheme.Burner().TStagIn(), s.pScheme.Burner().TStagIn(), 0.1)
	s.InDelta(s.scheme.Burner().TStagOut(), s.pScheme.Burner().TStagOut(), 0.1)
	s.approxEqual(initMassRate*bOMR, s.pScheme.Burner().MassRateOutput().GetState().Value().(float64), 1e-6)
	s.InDelta(s.scheme.Burner().Alpha(), s.pScheme.Burner().Alpha(), 5e-3)

	s.InDelta(s.scheme.TurboCascade().Turbine().PStagIn(), s.pScheme.CompressorTurbine().PStagIn(), 1e-6)
	s.InDelta(s.scheme.TurboCascade().Turbine().PStagOut(), s.pScheme.CompressorTurbine().PStagOut(), 1e-6)
	s.InDelta(s.scheme.TurboCascade().Turbine().TStagIn(), s.pScheme.CompressorTurbine().TStagIn(), 0.1)
	s.InDelta(s.scheme.TurboCascade().Turbine().TStagOut(), s.pScheme.CompressorTurbine().TStagOut(), 0.1)
	s.approxEqual(
		initMassRate*ctIMR,
		s.pScheme.CompressorTurbine().MassRateInput().GetState().Value().(float64), 1e-4,
	)
	s.InDelta(s.scheme.TurboCascade().Turbine().Eta(), s.pScheme.CompressorTurbine().Eta(), 1e-6)

	s.InDelta(s.scheme.CompressorTurbinePipe().PStagIn(), s.pScheme.CompressorTurbinePipe().PStagIn(), 1e-6)
	s.InDelta(s.scheme.CompressorTurbinePipe().PStagOut(), s.pScheme.CompressorTurbinePipe().PStagOut(), 1e-6)
	s.InDelta(s.scheme.CompressorTurbinePipe().TStagIn(), s.pScheme.CompressorTurbinePipe().TStagIn(), 0.1)
	s.InDelta(s.scheme.CompressorTurbinePipe().TStagOut(), s.pScheme.CompressorTurbinePipe().TStagOut(), 0.1)

	s.InDelta(s.scheme.FreeTurbineBlock().FreeTurbine().PStagIn(), s.pScheme.FreeTurbine().PStagIn(), 1e-6)
	s.InDelta(s.scheme.FreeTurbineBlock().FreeTurbine().PStagOut(), s.pScheme.FreeTurbine().PStagOut(), 1e-6)
	s.InDelta(s.scheme.FreeTurbineBlock().FreeTurbine().TStagIn(), s.pScheme.FreeTurbine().TStagIn(), 0.1)
	s.InDelta(s.scheme.FreeTurbineBlock().FreeTurbine().TStagOut(), s.pScheme.FreeTurbine().TStagOut(), 0.1)
	s.approxEqual(
		initMassRate*ftIMR,
		s.pScheme.FreeTurbine().MassRateInput().GetState().Value().(float64), 1e-4,
	)
	s.InDelta(s.scheme.FreeTurbineBlock().FreeTurbine().Eta(), s.pScheme.FreeTurbine().Eta(), 1e-6)

	s.approxEqual(
		s.scheme.FreeTurbineBlock().FreeTurbine().PowerOutput().GetState().Value().(float64),
		s.pScheme.FreeTurbine().PowerOutput().GetState().Value().(float64),
		1e-4,
	)

	s.InDelta(
		s.pScheme.FreeTurbine().PowerOutput().GetState().Value().(float64)*
			s.pScheme.FreeTurbine().MassRateInput().GetState().Value().(float64)+
			s.pScheme.Payload().PowerOutput().GetState().Value().(float64),
		0, 200,
	)

	s.approxEqual(
		-s.scheme.FreeTurbineBlock().FreeTurbine().PowerOutput().GetState().Value().(float64)*initMassRate*ftIMR,
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

func TestP2NBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(BuilderTestSuite))
}
