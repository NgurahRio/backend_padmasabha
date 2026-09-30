package models

type Villa struct {
	ID            uint   `gorm:"primaryKey;column:id_villa" json:"id_villa"`
	SubcategoryID string `gorm:"column:subcategoryId;type:longtext" json:"subcategoryId"`
	Name          string `gorm:"column:namevilla" json:"namevilla"`
	Description   string `gorm:"column:description" json:"description"`
	Imagedata     string `gorm:"column:imagedata" json:"imagedata"`
	CreatedAt     string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     string `gorm:"column:updated_at" json:"updated_at"`
	FacilityID    string `gorm:"column:facilityId;type:longtext" json:"facilityId"`
	Operational   string `gorm:"column:operational" json:"operational"`

	Subcategories []Subcategory `gorm:"many2many:villa_subcategories;joinForeignKey:ID;joinReferences:ID" json:"subcategories,omitempty"`
	Facilities    []Facility    `gorm:"many2many:villa_facilities;joinForeignKey:ID;joinReferences:IDFacility" json:"facilities,omitempty"`
}

func (Villa) TableName() string {
	return "villa"
}
