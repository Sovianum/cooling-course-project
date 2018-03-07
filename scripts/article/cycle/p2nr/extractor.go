package p2nr

import (
	"github.com/Sovianum/cooling-course-project/scripts/article/cycle/common"
	"github.com/Sovianum/cooling-course-project/scripts/article/cycle/p2n"
	"github.com/Sovianum/turbocycle/library/parametric/free2n"
)

func NewData2nr() Data2nr {
	return Data2nr{Data2n: p2n.NewData2n()}
}

type Data2nr struct {
	p2n.Data2n
	Sigma common.FloatArr `json:"sigma"`
}

func (data *Data2nr) Load(scheme free2n.DoubleShaftRegFreeScheme) {
	data.Data2n.Load(scheme)
	data.Sigma.Append(scheme.Regenerator().Sigma())
}
