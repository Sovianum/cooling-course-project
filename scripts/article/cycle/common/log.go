package common

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
)

func DetailedLog(iterNum int, precision float64, residual *mat.VecDense) {
	result := fmt.Sprintf("i: %d\t", iterNum)
	result += fmt.Sprintf("precision: %f\t", precision)
	result += fmt.Sprintf(
		"ggmr: %.5f\t ggPower: %.5f\t ftMR: %.5f\t ftPower: %.5f\t ftPressure: %.5f\t ggTemp: %.5f\t",
		residual.At(0, 0),
		residual.At(1, 0),
		residual.At(2, 0),
		residual.At(3, 0),
		residual.At(4, 0),
		residual.At(5, 0),
	)
	result += fmt.Sprintf("residual: %f", mat.Norm(residual, 2))
	fmt.Println(result)
}
