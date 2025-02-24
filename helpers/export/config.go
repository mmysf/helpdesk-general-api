package export

type Opts struct {
	configFields     []FieldConfig
	withStyle        bool
	withFilter       bool
	defaultRow       int
	defaultSheetName string
}

type OptFunc func(e *Opts)

func DefaultConfig() Opts {
	return Opts{
		configFields:     []FieldConfig{},
		defaultRow:       1,
		withStyle:        false,
		defaultSheetName: "Sheet1",
	}
}

func WithHeaders(headers []FieldConfig) func(*Opts) {
	return func(o *Opts) {
		o.configFields = headers
		o.defaultRow = 2
	}
}

func WithDefaultSheetName(sheetName string) func(*Opts) {
	return func(o *Opts) {
		o.defaultSheetName = sheetName
	}
}

func WithStyle(enable bool) func(o *Opts) {
	return func(o *Opts) {
		o.withStyle = enable
	}
}

func WithFilter(enable bool) func(o *Opts) {
	return func(o *Opts) {
		o.withFilter = enable
	}
}
