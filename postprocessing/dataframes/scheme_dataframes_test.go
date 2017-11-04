package dataframes

import (
	"testing"
	"github.com/Sovianum/cooling-course-project/core/schemes/three_shafts"
	"github.com/Sovianum/cooling-course-project/core"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"os"
)

const (
	power = 6000e3
	relaxCoef = 0.1
	iterNum = 100
)

func TestNewThreeShaftsDF_Smoke(t *testing.T) {
	var scheme = three_shafts.GetInitedThreeShaftsScheme()
	var pi = 10.
	var piFactor = 0.5

	var generator = core.GetDoubleCompressorDataGenerator(scheme, power, relaxCoef, iterNum)
	_, err := generator(pi, piFactor)
	assert.Nil(t, err)

	var df = NewThreeShaftsDF(power, scheme)
	var b, _ = json.MarshalIndent(df, "", "    ")
	os.Stdout.Write(b)
}
