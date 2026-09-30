package models

type Restaurant struct {
	ID            uint   `gorm:"primaryKey;column:id_restaurant" json:"id_restaurant"`
	SubcategoryID string `gorm:"column:subcategoryId;type:longtext" json:"subcategoryId"`
	Name          string `gorm:"column:namerestaurant" json:"namerestaurant"`
	Description   string `gorm:"column:description" json:"description"`
	Imagedata     string `gorm:"column:imagedata" json:"imagedata"`
	Do            string `gorm:"column:do" json:"do"`
	Dont          string `gorm:"column:dont" json:"dont"`
	Safety        string `gorm:"column:safety" json:"safety"`
	Maps          string `gorm:"column:maps" json:"maps"`
	CreatedAt     string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     string `gorm:"column:updated_at" json:"updated_at"`
	FacilityID    string `gorm:"column:facilityId;type:longtext" json:"facilityId"`
	Operational   string `gorm:"column:operational" json:"operational"`
}

func (Restaurant) TableName() string { return "restaurant" }
