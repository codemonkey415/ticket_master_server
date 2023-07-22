package models

type ShoeData struct {
	Heading     string   `bson:"heading" json:"heading" validate:"required"`
	Description string   `bson:"description" json:"description" validate:"required"`
	Gender      string   `bson:"gender" json:"gender" validate:"required"`
	Img         []string `bson:"img" json:"img" validate:"required"`
}
