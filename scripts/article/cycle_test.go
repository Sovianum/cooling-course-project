package article

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
	scheme *schemes.TwoShaftsSchemeImpl
	pScheme free2n.DoubleShaftFreeScheme
}

func (s *ParametricTestSuite) SetupTest() {
	s.scheme = get2nScheme(piStag).(*schemes.TwoShaftsSchemeImpl)
	s.pScheme, _=  getParametric(s.scheme)
}

func (s *ParametricTestSuite) TestSamePoint() {
	pc := s.pScheme.Compressor()
	pb := s.pScheme.Burner()
	pct := s.pScheme.CompressorTurbine()
	payload := s.pScheme.Payload()
	pft := s.pScheme.FreeTurbine()

	n, _ := s.pScheme.GetNetwork()
	_, e := n.Solve(0.1, 2, 100, 0.01)
	s.Require().Nil(e)

	s.InDelta(pc.MassRate(), pb.MassRateInput().GetState().Value().(float64), 1e7)
	s.True(common.ApproxEqual(
		pb.MassRateInput().GetState().Value().(float64),
		pct.MassRateInput().GetState().Value().(float64),
		3e-2,
	))

	fmt.Println(
		pc.MassRateInput().GetState().Value().(float64),
		pc.MassRateOutput().GetState().Value().(float64),

		pb.MassRateInput().GetState().Value().(float64),
		pb.MassRateOutput().GetState().Value().(float64),

		pct.MassRateInput().GetState().Value().(float64),
		pct.MassRateOutput().GetState().Value().(float64),

		pft.MassRateInput().GetState().Value().(float64),
		pft.MassRateOutput().GetState().Value().(float64),
	)
	fmt.Println(
		pc.PowerOutput().GetState().Value().(float64) * pc.MassRateInput().GetState().Value().(float64),
		pct.PowerOutput().GetState().Value().(float64) * pct.MassRateInput().GetState().Value().(float64),
		pc.PowerOutput().GetState().Value().(float64) * pc.MassRateInput().GetState().Value().(float64) +
			pct.PowerOutput().GetState().Value().(float64) * pct.MassRateInput().GetState().Value().(float64),
	)
	fmt.Println(
		payload.PowerOutput().GetState().Value().(float64),
		pft.PowerOutput().GetState().Value().(float64) * pft.MassRateInput().GetState().Value().(float64),
		payload.Power() + pft.PowerOutput().GetState().Value().(float64) * pft.MassRateInput().GetState().Value().(float64),
		payload.Power() - pft.PowerOutput().GetState().Value().(float64) * pft.MassRateInput().GetState().Value().(float64),
	)

	for _, nv := range s.pScheme.Assembler().GetNamedReport() {
		fmt.Printf("%s:\t\t\t%f\n", nv.Name, nv.Value)
	}
}

func TestParametricTestSuite(t *testing.T) {
	suite.Run(t, new(ParametricTestSuite))
}
