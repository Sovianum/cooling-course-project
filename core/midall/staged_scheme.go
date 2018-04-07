package midall

import (
	"fmt"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/impl/stage/compressor"
	"github.com/Sovianum/turbocycle/impl/stage/turbine"
	"github.com/Sovianum/turbocycle/library/schemes"
)

func NewStagedScheme3n(
	source schemes.ThreeShaftsScheme,
	lpcConfig, hpcConfig CompressorConfig,
	hptConfig, lptConfig, ftConfig TurbineConfig,
) (*StagedScheme3n, error) {
	msg := ""
	var err error
	solverGen := newton.NewUniformNewtonSolverGen(1e-5, newton.NoLog)
	result := &StagedScheme3n{}

	result.LPC, err = lpcConfig.GetFittedStagedCompressor(source.LPC(), solverGen)
	if err != nil {
		msg += fmt.Sprintf("lpcErr: %s;\n", err.Error())
		err = nil
	}
	result.HPC, err = hpcConfig.GetFittedStagedCompressor(source.HPC(), solverGen)
	if err != nil {
		msg += fmt.Sprintf("hpcErr: %s;\n", err.Error())
		err = nil
	}
	result.HPT, err = hptConfig.GetFittedStagedTurbine(source.HPT(), solverGen)
	if err != nil {
		msg += fmt.Sprintf("hptErr: %s;\n", err.Error())
		err = nil
	}
	result.LPT, err = lptConfig.GetFittedStagedTurbine(source.LPT(), solverGen)
	if err != nil {
		msg += fmt.Sprintf("lptErr: %s;\n", err.Error())
		err = nil
	}
	result.FT, err = ftConfig.GetFittedStagedTurbine(source.FT(), solverGen)
	if err != nil {
		msg += fmt.Sprintf("ftErr: %s;\n", err.Error())
		err = nil
	}
	if msg != "" {
		return nil, fmt.Errorf(msg)
	}
	return result, nil
}

type StagedScheme3n struct {
	LPC compressor.StagedCompressorNode
	HPC compressor.StagedCompressorNode

	HPT turbine.StagedTurbineNode
	LPT turbine.StagedTurbineNode
	FT  turbine.StagedTurbineNode
}
