package p2nr

import (
	"fmt"
	common2 "github.com/Sovianum/cooling-course-project/scripts/article/cycle/common"
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
	pScheme  free2n.DoubleShaftRegFreeScheme
	pNetwork graph.Network
	scheme   schemes.TwoShaftsRegeneratorScheme

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
	ftIMR := s.scheme.FreeTurbineBlock().FreeTurbine().MassRateInput().GetState().Value().(float64)

	s.InDelta(
		initMassRate,
		s.pScheme.Compressor().MassRateInput().GetState().Value().(float64),
		1e-6,
	)
	s.InDelta(s.scheme.Compressor().PStagIn(), s.pScheme.Compressor().PStagIn(), 1e-6)
	s.InDelta(s.scheme.Compressor().PStagOut(), s.pScheme.Compressor().PStagOut(), 1e-6)
	s.InDelta(s.scheme.Compressor().TStagIn(), s.pScheme.Compressor().TStagIn(), 1e-6)
	s.InDelta(s.scheme.Compressor().TStagOut(), s.pScheme.Compressor().TStagOut(), 1.5e-3)
	s.InDelta(s.scheme.Compressor().PiStag(), s.pScheme.Compressor().PiStag(), 1e-8)
	s.InDelta(s.scheme.Compressor().Eta(), s.pScheme.Compressor().Eta(), 1e-8)
	s.approxEqual(
		s.scheme.Compressor().PowerOutput().GetState().Value().(float64),
		s.pScheme.Compressor().PowerOutput().GetState().Value().(float64),
		1e-3,
	)

	s.InDelta(
		s.scheme.Regenerator().ColdInput().PressureInput().GetState().Value().(float64),
		s.pScheme.Regenerator().ColdInput().PressureInput().GetState().Value().(float64),
		1e-6,
	)
	s.InDelta(
		s.scheme.Regenerator().ColdInput().TemperatureInput().GetState().Value().(float64),
		s.pScheme.Regenerator().ColdInput().TemperatureInput().GetState().Value().(float64),
		0.1,
	)
	s.InDelta(
		s.scheme.Regenerator().HotInput().PressureInput().GetState().Value().(float64),
		s.pScheme.Regenerator().HotInput().PressureInput().GetState().Value().(float64),
		1e-6,
	)
	s.InDelta(
		s.scheme.Regenerator().HotInput().TemperatureInput().GetState().Value().(float64),
		s.pScheme.Regenerator().HotInput().TemperatureInput().GetState().Value().(float64),
		1e-3,
	)
	s.InDelta(
		s.scheme.Regenerator().Sigma(), s.pScheme.Regenerator().Sigma(), 3e-3,
	)

	s.InDelta(
		s.scheme.Regenerator().ColdOutput().PressureOutput().GetState().Value().(float64),
		s.pScheme.Regenerator().ColdOutput().PressureOutput().GetState().Value().(float64),
		1e-6,
	)
	s.InDelta(
		s.scheme.Regenerator().ColdOutput().TemperatureOutput().GetState().Value().(float64),
		s.pScheme.Regenerator().ColdOutput().TemperatureOutput().GetState().Value().(float64),
		1,
	)

	s.InDelta(s.scheme.Burner().PStagIn(), s.pScheme.Burner().PStagIn(), 1e-6)
	s.InDelta(s.scheme.Burner().PStagOut(), s.pScheme.Burner().PStagOut(), 1e-6)
	s.InDelta(s.scheme.Burner().TStagIn(), s.pScheme.Burner().TStagIn(), 1)
	s.InDelta(s.scheme.Burner().TStagOut(), s.pScheme.Burner().TStagOut(), 1)
	s.approxEqual(initMassRate*bOMR, s.pScheme.Burner().MassRateOutput().GetState().Value().(float64), 1e-6)
	s.InDelta(s.scheme.Burner().Alpha(), s.pScheme.Burner().Alpha(), 5e-4)

	s.InDelta(s.scheme.TurboCascade().Turbine().PStagIn(), s.pScheme.CompressorTurbine().PStagIn(), 1e-6)
	s.InDelta(s.scheme.TurboCascade().Turbine().PStagOut(), s.pScheme.CompressorTurbine().PStagOut(), 1e-6)
	s.InDelta(s.scheme.TurboCascade().Turbine().TStagIn(), s.pScheme.CompressorTurbine().TStagIn(), 1)
	s.InDelta(s.scheme.TurboCascade().Turbine().TStagOut(), s.pScheme.CompressorTurbine().TStagOut(), 1)
	s.approxEqual(
		initMassRate*ctIMR,
		s.pScheme.CompressorTurbine().MassRateInput().GetState().Value().(float64), 1e-3,
	)
	s.InDelta(s.scheme.TurboCascade().Turbine().Eta(), s.pScheme.CompressorTurbine().Eta(), 1e-6)

	s.InDelta(
		0,
		s.pScheme.Compressor().MassRateInput().GetState().Value().(float64)*
			s.pScheme.Compressor().PowerOutput().GetState().Value().(float64)+
			s.pScheme.CompressorTurbine().MassRateInput().GetState().Value().(float64)*
				s.pScheme.CompressorTurbine().PowerOutput().GetState().Value().(float64)*0.99,
		1e2,
	)

	common2.DetailedLog2Shaft(0, 0, s.pScheme.Assembler().GetVectorPort().GetState().(graph.VectorPortState).Vec)

	loss := s.scheme.FreeTurbineBlock().OutletPressureLoss()
	pLoss := s.pScheme.FreeTurbinePipe()
	s.InDelta(loss.PressureOutput().GetState().Value().(float64), pLoss.PressureOutput().GetState().Value().(float64), 1e-6)
	fmt.Println(s.pScheme.Regenerator().HotOutput().PressureOutput().GetState().Value().(float64))

	fmt.Println(
		s.pScheme.FreeTurbinePipe().PressureInput().GetState().Value().(float64),
		s.pScheme.FreeTurbinePipe().PressureOutput().GetState().Value().(float64),
	)

	s.InDelta(s.scheme.CompressorTurbinePipe().PStagIn(), s.pScheme.CompressorTurbinePipe().PStagIn(), 1e-6)
	s.InDelta(s.scheme.CompressorTurbinePipe().PStagOut(), s.pScheme.CompressorTurbinePipe().PStagOut(), 1e-6)
	s.InDelta(s.scheme.CompressorTurbinePipe().TStagIn(), s.pScheme.CompressorTurbinePipe().TStagIn(), 1)
	s.InDelta(s.scheme.CompressorTurbinePipe().TStagOut(), s.pScheme.CompressorTurbinePipe().TStagOut(), 1)

	s.InDelta(s.scheme.FreeTurbineBlock().FreeTurbine().PStagIn(), s.pScheme.FreeTurbine().PStagIn(), 1e-6)
	s.InDelta(s.scheme.FreeTurbineBlock().FreeTurbine().PStagOut(), s.pScheme.FreeTurbine().PStagOut(), 1e-6)
	s.InDelta(s.scheme.FreeTurbineBlock().FreeTurbine().TStagIn(), s.pScheme.FreeTurbine().TStagIn(), 1)
	s.InDelta(s.scheme.FreeTurbineBlock().FreeTurbine().TStagOut(), s.pScheme.FreeTurbine().TStagOut(), 1)
	s.approxEqual(
		initMassRate*ftIMR,
		s.pScheme.FreeTurbine().MassRateInput().GetState().Value().(float64), 1e-3,
	)
	s.InDelta(s.scheme.FreeTurbineBlock().FreeTurbine().Eta(), s.pScheme.FreeTurbine().Eta(), 1e-6)

	s.approxEqual(
		s.scheme.FreeTurbineBlock().FreeTurbine().PowerOutput().GetState().Value().(float64),
		s.pScheme.FreeTurbine().PowerOutput().GetState().Value().(float64),
		1e-3,
	)

	s.approxEqual(
		-s.scheme.FreeTurbineBlock().FreeTurbine().PowerOutput().GetState().Value().(float64)*initMassRate*ftIMR,
		s.pScheme.Payload().PowerOutput().GetState().Value().(float64),
		1e-6,
	)
	s.approxEqual(
		s.scheme.Regenerator().Sigma(),
		s.pScheme.Regenerator().Sigma(),
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
