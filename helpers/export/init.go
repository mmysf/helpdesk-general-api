package export

import (
	"github.com/xuri/excelize/v2"
)

type Exporter interface {
	// SetHeader(headers []FieldConfig) (err error)
	// AddStyle() (err error)
	AddRow(mapData map[string]any) (err error)
	// Cast(value any, fieldConfig FieldConfig) any
	ToBase64() (base64File string, err error)
	Save(filepath string) error
}

type excelExporter struct {
	xlsx              *excelize.File
	workingSheetName  string
	workingSheetIndex int
	currentRow        map[int]int

	Opts
}

func NewExporter(opts ...OptFunc) Exporter {

	option := DefaultConfig()

	for _, fn := range opts {
		fn(&option)
	}

	xlsx := excelize.NewFile()

	// set default sheet name
	xlsx.SetSheetName(xlsx.GetSheetName(0), option.defaultSheetName)

	activeSheetIndex := xlsx.GetActiveSheetIndex()

	return &excelExporter{
		Opts:              option,
		xlsx:              xlsx,
		workingSheetName:  option.defaultSheetName,
		workingSheetIndex: activeSheetIndex,
		currentRow: map[int]int{
			activeSheetIndex: option.defaultRow,
		},
	}
}
