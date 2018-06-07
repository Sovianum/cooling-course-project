package subcompress

import (
	"github.com/Sovianum/cooling-course-project/core/schemes/s3nsc"
	common2 "github.com/Sovianum/cooling-course-project/scripts/article/cycle/common"
	"github.com/Sovianum/turbocycle/common"
)

func Entry() error {
	scheme := s3nsc.GetInitedThreeShaftsSubCompressScheme()
	//scheme := s3n.GetDiplomaInitedThreeShaftsScheme()
	data := getSchemeDataTemplate(
		//common.Arange(12, 0.5, 30),
		[]float64{19},
		[]float64{0.529},
		[]float64{1.01, 1.2, 2},
		common.Arange(0.09, 0.01, 7),
	)
	if e := updateSchemeData(scheme, data); e != nil {
		return e
	}
	if e := common2.SaveData(data, common2.DataRoot+"3nsc_simple.json"); e != nil {
		return e
	}
	return nil
}
