package article

import (
	"fmt"
)

const (
	power     = 20e6
	relaxCoef = 0.1
	iterNum   = 10000
	precision = 0.01
	piStag = 10
)

func Entry() {
	scheme := get2nScheme(piStag)
	pScheme, pErr := getParametric(scheme)
	if pErr != nil {
		panic(pErr)
	}
	err := solveParametric(pScheme)
	if err != nil {
		panic(err)
	}

	pc := pScheme.Compressor()
	pb := pScheme.Burner()
	pct := pScheme.CompressorTurbine()
	pft := pScheme.FreeTurbine()

	fmt.Println(pc.PiStag())
	fmt.Println(pct.PiTStag())
	fmt.Println(pft.PiTStag())

	fmt.Println(
		pc.PStagIn(), pc.PStagOut(),
		pb.PStagIn(), pb.PStagOut(),
		pct.PStagIn(), pct.PStagOut(),
		pft.PStagIn(), pft.PStagOut(),
	)

	fmt.Println(
		pc.TStagIn(), pc.TStagOut(),
		pb.TStagIn(), pb.TStagOut(),
		pct.TStagIn(), pct.TStagOut(),
		pft.TStagIn(), pft.TStagOut(),
	)

	fmt.Println(pScheme.Payload().Power())
}
