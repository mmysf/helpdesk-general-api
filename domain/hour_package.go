package domain

type HourPackageRequest struct {
	Name          string   `json:"name"`
	Benefit       []string `json:"benefit"`
	Price         float64  `json:"price"`
	DurationHours int64    `json:"durationHours"`
}

type HourPackageUpdate struct {
	Name          string   `json:"name"`
	Benefit       []string `json:"benefit"`
	Price         float64  `json:"price"`
	DurationHours int64    `json:"durationHours"`
}

type HourPackageStatusUpdate struct {
	Status string `json:"status"`
}
