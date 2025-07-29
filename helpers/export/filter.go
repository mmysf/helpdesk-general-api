package export

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func (e *excelExporter) AddFilter() (err error) {
	if len(e.configFields) == 0 {
		return
	}

	lastCol, _ := excelize.ColumnNumberToName(len(e.configFields))
	rangeRef := fmt.Sprintf("A1:%s1", lastCol)

	err = e.xlsx.AutoFilter(e.workingSheetName, rangeRef, []excelize.AutoFilterOptions{})
	if err != nil {
		return
	}

	return
}
