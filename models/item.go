package models

type Item struct {
	ID          int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Title       string `gorm:"type:varchar(50)" json:"title"`
	Description string `gorm:"type:varchar(250)" json:"description"`
	CatID       int    `json:"cat_id"`
	UserID      int    `json:"-"`
	User        User   `gorm:"foreignkey:UserID" json:"-"`
}
