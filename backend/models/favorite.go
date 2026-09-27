package models

type Favorite struct {
	ID            uint `gorm:"primaryKey;column:id_bookmark" json:"id_bookmark"`
	UserID        uint `gorm:"column:userId" json:"userId"`
	DestinationID *uint `gorm:"column:destinationId" json:"destinationId"`
	EventID       *uint `gorm:"column:eventId" json:"eventId,omitempty"`

	// Relations
	User        *User        `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Destination *Destination `gorm:"foreignKey:DestinationID;references:ID" json:"destination,omitempty"`
	Event       *Event       `gorm:"foreignKey:EventID;references:ID" json:"event,omitempty"`
}

func (Favorite) TableName() string {
	return "bookmark"
}
