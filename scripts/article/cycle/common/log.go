package common

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
)

func DetailedLog2Shaft(iterNum int, precision float64, residual *mat.VecDense) {
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

func DetailedLog3Shaft(iterNum int, precision float64, residual *mat.VecDense) {
	result := fmt.Sprintf("i: %d\t", iterNum)
	result += fmt.Sprintf("precision: %f\t", precision)
	result += fmt.Sprintf(
		"hpc_mr: %.5f\t hpt_mr: %.5f\t lpt_mr: %.5f\t ft_mr: %.5f\t hp_po: %.5f\t lp_po: %.5f\t ft_po: %.5f\t ft_p: %.5f\t burn: %.5f\t",
		residual.At(0, 0),
		residual.At(1, 0),
		residual.At(2, 0),
		residual.At(3, 0),
		residual.At(4, 0),
		residual.At(5, 0),
		residual.At(6, 0),
		residual.At(7, 0),
		residual.At(8, 0),
	)
	result += fmt.Sprintf("residual: %f", mat.Norm(residual, 2))
	fmt.Println(result)
}

func DetailedLog3ShaftMidBurn(iterNum int, precision float64, residual *mat.VecDense) {
	result := fmt.Sprintf("i: %d\t", iterNum)
	result += fmt.Sprintf("precision: %f\t", precision)
	result += fmt.Sprintf(
		"hpc_mr: %.5f\t hpt_mr: %.5f\t lpt_mr: %.5f\t ft_mr: %.5f\t hp_po: %.5f\t lp_po: %.5f\t ft_po: %.5f\t ft_p: %.5f\t burn: %.5f\t mid_burn: %.5f",
		residual.At(0, 0),
		residual.At(1, 0),
		residual.At(2, 0),
		residual.At(3, 0),
		residual.At(4, 0),
		residual.At(5, 0),
		residual.At(6, 0),
		residual.At(7, 0),
		residual.At(8, 0),
		residual.At(9, 0),
	)
	result += fmt.Sprintf("residual: %f", mat.Norm(residual, 2))
	fmt.Println(result)
}
