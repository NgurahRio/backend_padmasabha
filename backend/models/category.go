package models

type Category struct {
	ID   uint   `gorm:"primaryKey;column:id_categories" json:"id_categories"`
	Name string `gorm:"column:name;not null" json:"name"`
}

func (Category) TableName() string {
	return "categories"
}
