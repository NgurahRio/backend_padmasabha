package models

type Subcategory struct {
	ID         uint     `gorm:"primaryKey;column:id_subcategories" json:"id_subcategories"`
	CategoryID uint     `gorm:"column:categoriesId;not null" json:"categoriesId"`
	Name       string   `gorm:"column:namesubcategories;not null" json:"namesubcategories"`
	Category   Category `gorm:"foreignKey:CategoryID;references:ID" json:"category"`
}

func (Subcategory) TableName() string {
	return "subcategories"
}
