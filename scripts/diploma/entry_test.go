package diploma

import (
	"encoding/json"
	"github.com/Sovianum/cooling-course-project/core"
	"github.com/Sovianum/cooling-course-project/core/profiling"
	"github.com/Sovianum/cooling-course-project/core/schemes/s3n"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestDataGeneration(t *testing.T) {
	scheme := getScheme(s3n.PiDiplomaLow, s3n.PiDiplomaHigh)
	data, err := getThreeShaftsSchemeData(
		scheme,
		power/etaR,
		2, 8, 40,
		2, 8, 40,
	)
	assert.NoError(t, err)

	newData := core.ConvertDoubleCompressorDataPoints(data)

	b, err := json.Marshal(newData)
	assert.NoError(t, err)

	err = profiling.SaveString("../../"+dataDir+"3n.json", string(b))
	assert.NoError(t, err)
}

func getThreeShaftsSchemeData(
	scheme schemes.ThreeShaftsScheme,
	power float64,
	piLowStart, piLowEnd float64, piLowStepNum int,
	piHighStart, piHighEnd float64, piHighStepNum int,
) ([]core.DoubleCompressorDataPoint, error) {
	var points []core.DoubleCompressorDataPoint
	generator := core.GetDoubleCompressorDataGenerator(scheme, power, relaxCoef, iterNum)

	for _, piLow := range common.LinSpace(piLowStart, piLowEnd, piLowStepNum) {
		for _, piHigh := range common.LinSpace(piHighStart, piHighEnd, piHighStepNum) {
			pi := piHigh * piLow
			piFactor := math.Log(piLow) / math.Log(pi)
			point, err := generator(pi, piFactor)
			if err != nil {
				return nil, err
			}
			points = append(points, point)
		}
	}
	return points, nil
}
