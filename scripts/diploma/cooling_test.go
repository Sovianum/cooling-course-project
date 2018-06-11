package diploma

import (
	"fmt"
	"github.com/Sovianum/cooling-course-project/core/midall/inited"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/impl/stage/geometry"
	"github.com/Sovianum/turbocycle/impl/stage/turbine"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/gap"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/profile"
	"github.com/Sovianum/turbocycle/utils/turbine/geom"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profiles"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
	"testing"
)

type coolingTestDataPack struct {
	stage        turbine.StageNode
	bladeProfile profiles.BladeProfile
	gapPack      gap.DataPack
}

func TestOptimizePSCoolingSystem(t *testing.T) {
	pack := getDataPack()
	//posVec0 := mat.NewVecDense(7, []float64{5e-3, 10e-3, 17e-3, 22e-3, 30e-3, 37e-3, 42e-3})
	posVec0 := mat.NewVecDense(1, []float64{10e-3})

	slitGeom := []SlitGeom{
		{0, 0.45e-3},
		//{0, 0.30e-3},
		//{0, 0.30e-3},
		//{0, 0.30e-3},
		//{0, 0.35e-3},
		//{0, 0.35e-3},
		//{0, 0.35e-3},
	}

	sysFunc := func(posVec *mat.VecDense) profile.TemperatureSystem {
		fmt.Println("pos vec: ", posVec)
		for i := range slitGeom {
			slitGeom[i].Coord = posVec.At(i, 0)
		}
		return getPSConvFilmTemperatureSystem(
			pack.gapPack.AlphaGas,
			pack.stage,
			pack.bladeProfile,
			slitGeom,
		)
	}

	res, err := optimizeCoolingSystem(sysFunc, posVec0, 1e-3, 1e-5, 0.00001, 100, newton.DetailedLog)
	assert.NoError(t, err)

	fmt.Println("res: ", res)
}

func getDataPack() coolingTestDataPack {
	initedMachines, err := inited.GetInitedStagedNodes()
	if err != nil {
		panic(err)
	}

	stage := initedMachines.HPT.Stages()[0]
	statorProfiler := getStatorProfiler(stage)

	statorMidProfile := profiles.NewBladeProfileFromProfiler(
		0.5,
		0.01, 0.01,
		0.2, 0.2,
		statorProfiler,
	)
	stagePack := stage.GetDataPack()
	statorMidProfile.Transform(geom.Scale(geometry.ChordProjection(stagePack.StageGeometry.StatorGeometry())))

	gapCalculator := getGapCalculator(stage, statorMidProfile)
	gapPack := gapCalculator.GetPack(coolAirMassRate)

	return coolingTestDataPack{
		stage:        stage,
		bladeProfile: statorMidProfile,
		gapPack:      gapPack,
	}
}
