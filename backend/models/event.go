package models

type Event struct {
	ID          uint    `gorm:"primaryKey;column:id_event" json:"id_event"`
	Name        string  `gorm:"column:nameevent;not null" json:"nameevent"`
	StartDate   string  `gorm:"column:start_date" json:"start_date"`
	EndDate     string  `gorm:"column:end_date" json:"end_date"`
	Description string  `gorm:"column:description" json:"description"`
	StartTime   string  `gorm:"column:start_time" json:"start_time"`
	EndTime     string  `gorm:"column:end_time" json:"end_time"`
	Price       float64 `gorm:"column:price" json:"price"`
	ImageEvent  string  `gorm:"column:image_event" json:"image_event"`
}

func (Event) TableName() string {
	return "event"
}
