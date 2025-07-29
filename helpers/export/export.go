package export

import "encoding/base64"

func (e *excelExporter) ToBase64() (base64File string, err error) {
	// pre render
	err = e._preRender()
	if err != nil {
		return
	}

	// write to buffer
	buff, err := e.xlsx.WriteToBuffer()
	if err != nil {
		return
	}

	// convert to base64
	base64File = base64.StdEncoding.EncodeToString(buff.Bytes())
	return
}

func (e *excelExporter) Save(filepath string) (err error) {
	// pre render
	err = e._preRender()
	if err != nil {
		return
	}

	// saving a file
	err = e.xlsx.SaveAs(filepath)
	if err != nil {
		return
	}

	return
}

func (e *excelExporter) _preRender() (err error) {
	// add header
	err = e.SetHeader(e.configFields)
	if err != nil {
		return
	}

	// add style if needed
	if e.withStyle {
		err = e.AddStyle()
		if err != nil {
			return
		}
	}

	// add filter if needed
	if e.withFilter {
		err = e.AddFilter()
		if err != nil {
			return
		}
	}

	return
}
