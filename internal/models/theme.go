package models

type Theme struct {
	Base
	Name        string `gorm:"size:50;not null;unique" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	PrimaryColor string `gorm:"size:20" json:"primary_color"`
	LogoURL     string `gorm:"size:255" json:"logo_url"`
	Active      bool   `gorm:"default:true" json:"active"`
	Users       []User `gorm:"foreignKey:ThemeID" json:"-"`
}