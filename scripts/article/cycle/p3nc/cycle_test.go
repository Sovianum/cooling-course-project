package p3nc

import (
	"fmt"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/library/parametric/free3n"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ParametricTestSuite struct {
	suite.Suite
	scheme  schemes.ThreeShaftsCoolerScheme
	pScheme free3n.ThreeShaftCoolFreeScheme
}

func (s *ParametricTestSuite) SetupTest() {
	s.scheme = GetScheme(lpcPiStag, hpcPiStag)
	s.pScheme, _ = GetParametric(s.scheme)
}

func (s *ParametricTestSuite) TestSamePoint() {
	lpc := s.pScheme.LPC()
	//hpc := s.pScheme.LPC()

	b := s.pScheme.Burner()

	hpt := s.pScheme.HPT()
	//lpt := s.pScheme.LPT()

	n, _ := s.pScheme.GetNetwork()
	e := n.Solve(0.1, 2, 100, 0.01)
	s.Require().Nil(e)

	s.InDelta(lpc.MassRate(), b.MassRateInput().GetState().Value().(float64), 1e-7)
	s.True(
		common.ApproxEqual(
			b.MassRateOutput().GetState().Value().(float64),
			hpt.MassRateInput().GetState().Value().(float64),
			1e-5,
		),
		"expected: %f, got: %f",
		b.MassRateOutput().GetState().Value().(float64),
		hpt.MassRateInput().GetState().Value().(float64),
	)

	for _, nv := range s.pScheme.Assembler().GetNamedReport() {
		fmt.Printf("%s:\t\t\t%f\n", nv.Name, nv.Value)
	}
}

func TestParametricTestSuite(t *testing.T) {
	suite.Run(t, new(ParametricTestSuite))
}
