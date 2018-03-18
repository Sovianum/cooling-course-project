package dataframes

import "github.com/Sovianum/cooling-course-project/core"

type VariantDF struct {
	PiTotal   float64
	PiLow     float64
	PiHigh    float64
	MaxEta    core.DoubleCompressorDataPoint
	MaxLabour core.DoubleCompressorDataPoint
}
