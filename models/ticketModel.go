package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ticket struct {
	Id                 primitive.ObjectID `bson:"_id,omitempty"`
	EventId            string             `bson:"event_id,omitempty" json:"event_id,omitempty"`
	TmSeatId           string             `bson:"tm_seat_id,omitempty" json:"tm_seat_id,omitempty"`
	BasicEventScrapeAt string             `bson:"basicevent_scrape_at,omitempty" json:"basicevent_scrape_at,omitempty"`
	LineNumber         int                `bson:"line_number,omitempty" json:"line_number,omitempty"`
	Position           int                `bson:"position,omitempty" json:"position,omitempty"`
	RowName            string             `bson:"row_name,omitempty" json:"row_name,omitempty"`
	SeatGroup          any                `bson:"seat_group,omitempty" json:"seat_group,omitempty"`
	SeatName           any                `bson:"seat_name,omitempty" json:"seat_name,omitempty"`
	SectionName        string             `bson:"section_name,omitempty" json:"section_name,omitempty"`
	UniqueId           any                `bson:"unique_id,omitempty" json:"unique_id,omitempty"`
	Accessibility      any                `bson:"accessibility,omitempty" json:"accessibility,omitempty"`
	Attributes         string             `bson:"attributes,omitempty" json:"attributes,omitempty"`
	Currency           string             `bson:"currency,omitempty" json:"currency,omitempty"`
	IsAvailable        int                `bson:"is_available,omitempty" json:"is_available,omitempty"`
	Notes              string             `bson:"notes,omitempty" json:"notes,omitempty"`
	OfferId1           string             `bson:"offer_id1,omitempty" json:"offer_id1,omitempty"`
	OfferId2           string             `bson:"offer_id2,omitempty" json:"offer_id2,omitempty"`
	Price              any                `bson:"price,omitempty" json:"price,omitempty"`
}

// type Ticket struct {
// 	ID               primitive.ObjectID  `bson:"_id,omitempty"`
// 	AccountID        int                 `bson:"account_id,omitempty"`
// 	AccountType      string              `bson:"account_type,omitempty"`
// 	AccountTypeID    int                 `bson:"account_type_id,omitempty"`
// 	Barcode          string              `bson:"barcode,omitempty"`
// 	ConfNum          string              `bson:"conf_num,omitempty"`
// 	Event            string              `bson:"event,omitempty"`
// 	EventID          int                 `bson:"event_id,omitempty"`
// 	Expired          bool                `bson:"expired,omitempty"`
// 	HasAssets        bool                `bson:"has_assets,omitempty"`
// 	HasBarcode       bool                `bson:"has_barcode,omitempty"`
// 	IsParking        *bool               `bson:"is_parking,omitempty"`
// 	LastModified     *primitive.DateTime `bson:"last_modified,omitempty"`
// 	LocalDate        *primitive.DateTime `bson:"local_date,omitempty"`
// 	Locked           *bool               `bson:"locked,omitempty"`
// 	PrimaryEventID   string              `bson:"primary_event_id,omitempty"`
// 	Quantity         int                 `bson:"quantity,omitempty"`
// 	Row              string              `bson:"row,omitempty"`
// 	Seat             string              `bson:"seat,omitempty"`
// 	Section          string              `bson:"section,omitempty"`
// 	TicketStatusCode int                 `bson:"ticket_status_code,omitempty"`
// 	Transferred      *bool               `bson:"transferred,omitempty"`
// 	UserID           int                 `bson:"user_id,omitempty"`
// 	Username         string              `bson:"username,omitempty"`
// 	Venue            string              `bson:"venue,omitempty"`
// 	AppUser          primitive.ObjectID  `bson:"app_user,omitempty"`
// 	AppTransfer      primitive.ObjectID  `bson:"app_transfer,omitempty"`
// }
