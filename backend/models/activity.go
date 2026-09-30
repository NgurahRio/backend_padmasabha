package models

type Activity struct {
	ID            uint   `gorm:"primaryKey;column:id_activity" json:"id_activity"`
	SubcategoryID string `gorm:"column:subcategoryId;type:longtext" json:"subcategoryId"`
	Name          string `gorm:"column:nameactivity" json:"nameactivity"`
	Description   string `gorm:"column:description" json:"description"`
	Imagedata     string `gorm:"column:imagedata" json:"imagedata"`
	CreatedAt     string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     string `gorm:"column:updated_at" json:"updated_at"`
	FacilityID    string `gorm:"column:facilityId;type:longtext" json:"facilityId"`
	Operational   string `gorm:"column:operational" json:"operational"`
}

func (Activity) TableName() string { return "activity" }
