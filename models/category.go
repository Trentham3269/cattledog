package models

type Category struct {
	ID int `gorm:"primary_key" json:"id"`
	Name string `gorm:"size:50" json:"name"`
}