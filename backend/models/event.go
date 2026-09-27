package models

type Event struct {
	ID            uint     `gorm:"primaryKey;column:id_event" json:"id_event"`
	DestinationID *uint    `gorm:"column:destinationId" json:"destinationId"`
	Name          string   `gorm:"column:nameevent;not null" json:"nameevent"`
	StartDate     string   `gorm:"column:start_date" json:"start_date"`
	EndDate       string   `gorm:"column:end_date" json:"end_date"`
	Description   string   `gorm:"column:description" json:"description"`
	StartTime     string   `gorm:"column:start_time" json:"start_time"`
	EndTime       string   `gorm:"column:end_time" json:"end_time"`
	Price         float64  `gorm:"column:price" json:"price"`
	Maps          string   `gorm:"column:maps" json:"maps"`
	Do            string   `gorm:"column:do" json:"do"`
	Dont          string   `gorm:"column:dont" json:"dont"`
	Safety        string   `gorm:"column:safety" json:"safety"`
	Longitude     *float64 `gorm:"column:longitude" json:"longitude"`
	Latitude      *float64 `gorm:"column:latitude" json:"latitude"`
	ImageEvent    string   `gorm:"column:image_event" json:"image_event"`
	Location      string   `gorm:"column:location" json:"location"`

	Destination *Destination `gorm:"foreignKey:DestinationID;references:ID" json:"destination,omitempty"`
}

func (Event) TableName() string {
	return "event"
}
