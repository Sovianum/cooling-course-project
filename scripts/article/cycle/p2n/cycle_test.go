package p2n

import (
	"github.com/stretchr/testify/suite"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/turbocycle/library/parametric/free2n"
	"testing"
	"github.com/Sovianum/turbocycle/common"
	"fmt"
)

type ParametricTestSuite struct {
	suite.Suite
	scheme schemes.TwoShaftsScheme
	pScheme free2n.DoubleShaftFreeScheme
}

func (s *ParametricTestSuite) SetupTest() {
	s.scheme = GetScheme(piStag)
	s.pScheme, _=  GetParametric(s.scheme)
}

func (s *ParametricTestSuite) TestSamePoint() {
	pc := s.pScheme.Compressor()
	pb := s.pScheme.Burner()
	pct := s.pScheme.CompressorTurbine()

	n, _ := s.pScheme.GetNetwork()
	_, e := n.Solve(0.1, 2, 100, 0.01)
	s.Require().Nil(e)

	s.InDelta(pc.MassRate(), pb.MassRateInput().GetState().Value().(float64), 1e-7)
	s.True(
		common.ApproxEqual(
			pb.MassRateInput().GetState().Value().(float64),
			pct.MassRateInput().GetState().Value().(float64),
			5e-2,
		),
		"expected: %f, got: %f",
		pb.MassRateInput().GetState().Value().(float64),
		pct.MassRateInput().GetState().Value().(float64),
	)

	for _, nv := range s.pScheme.Assembler().GetNamedReport() {
		fmt.Printf("%s:\t\t\t%f\n", nv.Name, nv.Value)
	}
}

func TestParametricTestSuite(t *testing.T) {
	suite.Run(t, new(ParametricTestSuite))
}
