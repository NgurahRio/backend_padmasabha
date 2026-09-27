package models

type SubPackage struct {
	ID          uint   `gorm:"primaryKey;column:id_subpackage" json:"id_subpackage"`
	Packagetype string `gorm:"column:jenispackage" json:"jenispackage"`
	Image       []byte `gorm:"column:image" json:"image"`
}

func (SubPackage) TableName() string {
	return "subpackage"
}
