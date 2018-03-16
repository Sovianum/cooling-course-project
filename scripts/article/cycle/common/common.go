package common

import (
	"os"
	"encoding/json"
)

func SaveData(data interface{}, path string) error {
	b, _ := json.Marshal(data)
	f, e := os.Create(path)
	if e != nil {
		return e
	}
	_, e = f.WriteString(string(b))
	return e
}

type FloatArr []float64

func NewFloatArr() *FloatArr {
	result := FloatArr(make([]float64, 0))
	return &result
}

func (fArr *FloatArr) Append(f float64) *FloatArr {
	if fArr == nil {
		*fArr = FloatArr([]float64{f})
		return fArr
	}
	*fArr = FloatArr(append(*fArr, f))
	return fArr
}

func (fArr *FloatArr) Len() int  {
	if fArr == nil {
		return 0
	}
	return len(*fArr)
}

func (fArr *FloatArr) At(i int) float64 {
	return (*fArr)[i]
}
