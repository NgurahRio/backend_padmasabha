package models

type Review struct {
	ID            uint          `gorm:"primaryKey;column:id_review" json:"id_review"`
	UserID        uint          `gorm:"column:userId" json:"userId"`
	User          User          `gorm:"foreignKey:UserID;references:ID" json:"user"`
	DestinationID *uint         `gorm:"column:destinationId" json:"destinationId"`
	Destination   *Destination  `gorm:"foreignKey:DestinationID;references:ID" json:"destination,omitempty"`
	EventID       *uint         `gorm:"column:eventId" json:"eventId"`
	Event         *Event        `gorm:"foreignKey:EventID;references:ID" json:"event,omitempty"`
	Rating        int           `gorm:"column:rating" json:"rating"`
	Comment       string        `gorm:"column:comment" json:"comment"`
	CreatedAt     string        `gorm:"column:created_at" json:"created_at"`
}

func (Review) TableName() string {
	return "review"
}
