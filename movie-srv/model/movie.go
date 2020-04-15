package model

type Movie struct {
	Id     string
	Title  string `gorm:"size:65535"`
	Url    string `gorm:"size:65535"`
	Status string
}
