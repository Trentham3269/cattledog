package models

type Category struct {
	ID    int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name  string `gorm:"type:varchar(50)" json:"name"`
	Items []Item `gorm:"foreignkey:CatId" json:"items"`
}
