package p3nc

func Entry() {
	scheme := GetScheme(lpcPiStag, hpcPiStag)
	pScheme, pErr := GetParametric(scheme)
	if pErr != nil {
		panic(pErr)
	}
	err := SolveParametric(pScheme)
	if err != nil {
		panic(err)
	}
}
