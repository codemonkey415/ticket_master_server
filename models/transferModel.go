package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transfer struct {
	TransferID          int64              `bson:"transfer_id" json:"transfer_id"`
	StartTime           time.Time          `bson:"start_time" json:"start_time"`
	EndDate             time.Time          `bson:"end_date" json:"end_date"`
	MarketplaceID       int64              `bson:"marketplace_id" json:"marketplace_id"`
	Marketplace         string             `bson:"marketplace" json:"marketplace"`
	UserID              int64              `bson:"user_id" json:"user_id"`
	OrderID             int64              `bson:"order_id" json:"order_id"`
	InvoiceID           string             `bson:"invoice_id" json:"invoice_id"`
	RecipientEmail      string             `bson:"recipient_email" json:"recipient_email"`
	RecipientFirstName  string             `bson:"recipient_first_name" json:"recipient_first_name"`
	RecipientLastName   string             `bson:"recipient_last_name" json:"recipient_last_name"`
	Note                string             `bson:"note" json:"note"`
	Status              string             `bson:"status" json:"status"`
	StatusCode          int64              `bson:"status_code" json:"status_code"`
	MappingError        string             `bson:"mapping_error" json:"mapping_error"`
	AppliedMappingRules []string           `bson:"applied_mapping_rules" json:"applied_mapping_rules"`
	TransferType        string             `bson:"transfer_type" json:"transfer_type"`
	InitiatedBy         string             `bson:"initiated_by" json:"initiated_by"`
	RecalledBy          string             `bson:"recalled_by" json:"recalled_by"`
	OpuUserID           int64              `bson:"opu_user_id" json:"opu_user_id"`
	Confirmed           bool               `bson:"confirmed" json:"confirmed"`
	PrimaryEventID      string             `bson:"primary_event_id" json:"primary_event_id"`
	Requests            []interface{}      `bson:"requests" json:"requests"`
	User                primitive.ObjectID `bson:"user,omitempty" json:"-"`
}
