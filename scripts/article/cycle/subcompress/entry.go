package subcompress

import (
	"github.com/Sovianum/cooling-course-project/core/schemes/s3nsc"
	common2 "github.com/Sovianum/cooling-course-project/scripts/article/cycle/common"
	"github.com/Sovianum/turbocycle/common"
)

func Entry() error {
	scheme := s3nsc.GetInitedThreeShaftsSubCompressScheme()
	data := getSchemeDataTemplate(
		common.Arange(12, 0.5, 30),
		[]float64{0.5},
		[]float64{1.01, 2},
		common.Arange(0.1, 0.01, 7),
	)
	if e := updateSchemeData(scheme, data); e != nil {
		return e
	}
	if e := common2.SaveData(data, common2.DataRoot+"3nsc_simple.json"); e != nil {
		return e
	}
	return nil
}
