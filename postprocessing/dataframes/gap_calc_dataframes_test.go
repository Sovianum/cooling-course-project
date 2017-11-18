package dataframes

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGapGasDF_AirDataIterator(t *testing.T) {
	var df = GapGasDF{
		AirMassRate:[]float64{1, 1, 1},
		DCoef:[]float64{1, 1, 1},
		EpsCoef:[]float64{1, 1, 1},
		AirGap:[]float64{1, 1, 1},
	}
	var i = 0
	for range df.TableRows() {
		i++
	}
	assert.Equal(t, 3, i)
}