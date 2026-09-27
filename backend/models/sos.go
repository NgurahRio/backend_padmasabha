package models

type SOS struct {
	ID      uint   `gorm:"primaryKey;column:id_sos" json:"id_sos"`
	Name    string `gorm:"column:name_sos;not null" json:"name_sos"`
	Alamat  string `gorm:"column:alamat_sos;not null" json:"alamat_sos"`
	Telepon string `gorm:"column:telepon;not null" json:"telepon"`
}

func (SOS) TableName() string {
	return "sos"
}
