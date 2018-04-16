package profiling

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/stage/states"
	"github.com/Sovianum/turbocycle/impl/stage/turbine"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/laws"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profilers"
)

func GetInitedStatorProfiler(
	geomGen turbine.BladingGeometryGenerator,
	meanInletTriangle, meanOutletTriangle states.VelocityTriangle,
) profilers.Profiler {
	return profilers.NewProfiler(
		profilers.ProfilerConfig{
			Windage:    1,
			ApproxTRel: 0.7,

			Behavior: profilers.NewStatorProfilingBehavior(),
			GeomGen:  geomGen,

			MeanInletTriangle:  meanInletTriangle,
			MeanOutletTriangle: meanOutletTriangle,
			InletVelocityLaw:   laws.NewConstantAbsoluteAngleLaw(),
			OutletVelocityLaw:  laws.NewConstantAbsoluteAngleLaw(),

			InletProfileAngleFunc: func(characteristicAngle, hRel float64) float64 {
				return characteristicAngle
			},
			OutletProfileAngleFunc: func(characteristicAngle, hRel float64) float64 {
				return characteristicAngle
			},

			InstallationAngleFunc: func(hRel float64) float64 {
				return common.InterpTolerate(
					hRel,
					[]float64{0, 1},
					[]float64{common.ToRadians(50), common.ToRadians(50)},
				)
			},

			InletExpansionAngleFunc: func(hRel float64) float64 {
				return common.InterpTolerate(
					hRel,
					[]float64{0, 1},
					[]float64{common.ToRadians(15), common.ToRadians(30)},
				)
			},
			OutletExpansionAngleFunc: func(hRel float64) float64 {
				return common.ToRadians(5)
			},

			InletPSAngleFractionFunc: func(hRel float64) float64 {
				return 0.5
			},
			OutletPSAngleFractionFunc: func(hRel float64) float64 {
				return 1 / 3
			},
		},
	)
}

func GetInitedRotorProfiler(
	geomGen turbine.BladingGeometryGenerator,
	meanInletTriangle, meanOutletTriangle states.VelocityTriangle,
) profilers.Profiler {
	var inletLaw = laws.NewConstantAbsoluteAngleLaw()
	var outletLaw = laws.NewConstantLabourLaw(inletLaw, meanInletTriangle)

	return profilers.NewProfiler(
		profilers.ProfilerConfig{
			Windage:    1,
			ApproxTRel: 0.7,

			Behavior: profilers.NewRotorProfilingBehavior(),
			GeomGen:  geomGen,

			MeanInletTriangle:  meanInletTriangle,
			MeanOutletTriangle: meanOutletTriangle,
			InletVelocityLaw:   inletLaw,
			OutletVelocityLaw:  outletLaw,

			InletProfileAngleFunc: func(characteristicAngle, hRel float64) float64 {
				return characteristicAngle + common.ToRadians(2)
			},
			OutletProfileAngleFunc: func(characteristicAngle, hRel float64) float64 {
				return characteristicAngle
			},

			InstallationAngleFunc: func(hRel float64) float64 {
				return common.InterpTolerate(
					hRel,
					[]float64{0, 0.5, 1},
					[]float64{
						common.ToRadians(68),
						common.ToRadians(55),
						common.ToRadians(50),
					},
				)
			},

			InletExpansionAngleFunc: func(hRel float64) float64 {
				return common.InterpTolerate(
					hRel,
					[]float64{0, 1},
					[]float64{common.ToRadians(20), common.ToRadians(15)},
				)
			},
			OutletExpansionAngleFunc: func(hRel float64) float64 {
				return common.ToRadians(5)
			},

			InletPSAngleFractionFunc: func(hRel float64) float64 {
				return common.InterpTolerate(
					hRel,
					[]float64{0, 1},
					[]float64{0.7, 0.9},
				)
			},
			OutletPSAngleFractionFunc: func(hRel float64) float64 {
				return common.InterpTolerate(
					hRel,
					[]float64{0, 1},
					[]float64{1 / 3, 1},
				)
			},
		},
	)
}
