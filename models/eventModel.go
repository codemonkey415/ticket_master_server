package models

type Event struct {
	EventId   string `bson:"event_id" json:"event_id"`
	Name      string `bson:"name" json:"name"`
	Venue     string `bson:"venue" json:"venue"`
	EventDate any    `bson:"event_date" json:"event_date"`
}
