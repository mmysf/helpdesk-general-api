package export

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

type FieldConfig struct {
	Key                string
	Label              string
	As                 string
	Default            string
	DateFormat         string
	DateFormatLocation *time.Location
	DateParseLayout    string
	DateParseLocation  *time.Location
	LongestChar        int
}

func (e *excelExporter) SetHeader(configFields []FieldConfig) (err error) {

	e.configFields = configFields

	// header
	for i, fieldConfig := range e.configFields {
		// get cell name
		col, _ := excelize.CoordinatesToCellName(i+1, 1)

		// set value header
		e.xlsx.SetCellValue(e.workingSheetName, col, fieldConfig.Label)

		// set longest char
		if len(fieldConfig.Label) > fieldConfig.LongestChar {
			e.configFields[i].LongestChar = len(fieldConfig.Label)
		}
	}

	return
}

func (e *excelExporter) AddRow(mapData map[string]any) (err error) {
	rowValues := make([]any, 0)
	for _, fieldConfig := range e.configFields {
		var cellValue any

		// get from mapData
		value, ok := mapData[fieldConfig.Key]
		if !ok {
			cellValue = fieldConfig.Default
		} else {
			cellValue = e.Cast(value, fieldConfig)
		}

		// cast cell value
		rowValues = append(rowValues, cellValue)
	}

	// set row value
	e.xlsx.SetSheetRow(e.workingSheetName, fmt.Sprintf("A%d", e.currentRow[e.workingSheetIndex]), &rowValues)

	// set longest char
	for i, fieldConfig := range e.configFields {
		col, _ := excelize.CoordinatesToCellName(i+1, e.currentRow[e.workingSheetIndex])
		val, _ := e.xlsx.GetCellValue(e.workingSheetName, col)

		if len(val) > fieldConfig.LongestChar {
			e.configFields[i].LongestChar = len(val)
		}
	}

	e.currentRow[e.workingSheetIndex]++

	return
}
