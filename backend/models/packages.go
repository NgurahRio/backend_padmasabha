package models

type Packages struct {
	ID             uint   `gorm:"primaryKey;column:id_packages" json:"id_packages"`
	VillaID        uint   `gorm:"column:villaId" json:"villaId"`
	SubPackageData string `gorm:"column:subpackage_data" json:"subpackage_data"`

	Villa Villa `gorm:"foreignKey:VillaID;references:ID" json:"villa,omitempty"`
}

func (Packages) TableName() string {
	return "packages"
}
