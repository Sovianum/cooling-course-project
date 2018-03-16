package p2n

import (
	"github.com/Sovianum/cooling-course-project/scripts/article/cycle/common"
)

func Entry() {
	scheme := GetScheme(piStag)
	schemeData, err := GetSchemeData(scheme)
	if err != nil {
		panic(err)
	}

	if err := common.SaveData(schemeData, common.DataRoot + "2n_simple.json"); err != nil {
		panic(err)
	}

	OptimizeScheme(scheme, schemeData)

	pScheme, pErr := GetParametric(scheme)
	if pErr != nil {
		panic(pErr)
	}

	pData, err := SolveParametric(pScheme)
	if err != nil {
		panic(err)
	}

	if err := common.SaveData(pData, common.DataRoot + "2n.json"); err != nil {
		panic(err)
	}
}
