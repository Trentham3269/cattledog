package models

type Item struct {
	ID          int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Title       string `gorm:"type:varchar(50)" json:"title"`
	Description string `gorm:"type:varchar(250)" json:"desc"`
	Cat_id      int    `json:"cat_id"`
}