package models

type Event struct {
	EventId   string `bson:"event_id" json:"event_id"`
	Name      string `bson:"name" json:"name"`
	Label     string `bson:"label" json:"label"`
	Venue     string `bson:"venue" json:"venue"`
	IsActive  int32  `bson:"is_active" json:"is_active"`
	EventDate any    `bson:"event_date" json:"event_date"`
}
