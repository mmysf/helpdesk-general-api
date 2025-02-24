package export

import "github.com/xuri/excelize/v2"

func (e *excelExporter) AddStyle() (err error) {
	style, errStyle := e.xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 15,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E0EBF5"},
			Pattern: 1,
		},
	})
	if errStyle != nil {
		err = errStyle
		return
	}

	lastHeaderCol, _ := excelize.CoordinatesToCellName(len(e.configFields), 1)
	e.xlsx.SetCellStyle(e.workingSheetName, "A1", lastHeaderCol, style)

	// set column width
	for i, fieldConfig := range e.configFields {
		// get cell name
		colName, _ := excelize.ColumnNumberToName(i + 1)
		e.xlsx.SetColWidth(e.workingSheetName, colName, colName, float64(fieldConfig.LongestChar)*1.123)
	}

	return
}
