package domain

type ServerPackageRequest struct {
	Name         string   `json:"name"`
	Benefit      []string `json:"benefit"`
	Price        float64  `json:"price"`
	Validity     int64    `json:"validity"`
	Customizable bool     `bson:"customizable" json:"customizable"`
}

type ServerPackageUpdate struct {
	Name         string   `json:"name"`
	Benefit      []string `json:"benefit"`
	Price        float64  `json:"price"`
	Validity     int64    `json:"validity"`
	Customizable bool     `bson:"customizable" json:"customizable"`
}

type ServerPackageStatusUpdate struct {
	Status string `json:"status"`
}
