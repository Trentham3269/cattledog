package models

type User struct {
	ID       int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Email    string `gorm:"type:varchar(80)" json:"email"`
	Password string `gorm:"type:varchar(80)" json:"password"`
}