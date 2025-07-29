package export

import (
	"time"
)

func (e *excelExporter) Cast(value any, fieldConfig FieldConfig) any {
	if fieldConfig.As == "boolean" || fieldConfig.As == "bool" {
		if v, ok := value.(bool); ok {
			if v {
				return "✅"
			} else {
				return "❌"
			}
		}
	}

	if fieldConfig.As == "date" {

		if t, ok := value.(*time.Time); ok && t == nil {
			return fieldConfig.Default
		}

		if fieldConfig.DateFormatLocation == nil {
			fieldConfig.DateFormatLocation = time.Local
		}

		if fieldConfig.DateFormat != "" {
			if fieldConfig.DateFormat == "DATETIME" {
				fieldConfig.DateFormat = "2006-01-02 15:04:05"
			} else if fieldConfig.DateFormat == "DATE" {
				fieldConfig.DateFormat = "2006-01-02"
			}
		} else {
			fieldConfig.DateFormat = time.RFC3339
		}

		if t, ok := value.(time.Time); ok {
			return t.In(fieldConfig.DateFormatLocation).Format(fieldConfig.DateFormat)
		}

		if tString, ok := value.(string); ok {
			if fieldConfig.DateParseLocation == nil {
				fieldConfig.DateParseLocation = time.Local
			}

			if fieldConfig.DateParseLayout == "" {
				fieldConfig.DateParseLayout = time.RFC3339
			}

			if t, err := time.ParseInLocation(time.RFC3339, tString, fieldConfig.DateParseLocation); err == nil {
				return t.In(fieldConfig.DateFormatLocation).Format(fieldConfig.DateFormat)
			}
		}

	}

	return value
}
