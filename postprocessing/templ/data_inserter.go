package templ

import (
	"io/ioutil"
	"os"
)

const (
	templateName = "template"
)

func NewDataInserter(templateFilePath, outputFilePath string) DataInserter {
	return &dataInserter{
		templateFilePath: templateFilePath,
		outputFilePath:   outputFilePath,
	}
}

type DataInserter interface {
	Insert(data interface{}) error
}

type dataInserter struct {
	templateFilePath string
	outputFilePath   string
}

func (inserter *dataInserter) Insert(data interface{}) error {
	var f, fErr = ioutil.ReadFile(inserter.templateFilePath)
	if fErr != nil {
		return fErr
	}
	var funcMap = GetFuncMap()
	var t, tErr = GetTemplate(
		templateName,
		string(f),
		funcMap,
	)
	if tErr != nil {
		return tErr
	}

	var out, outErr = os.Create(inserter.outputFilePath)
	defer out.Close()

	if outErr != nil {
		return outErr
	}
	var execErr = t.Execute(out, data)
	if execErr != nil {
		return execErr
	}
	return nil
}
