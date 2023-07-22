package models

type Kid struct {
	Title       string   `bson:"title" json:"title" validate:"required"`
	Gender      string   `bson:"gender" json:"gender" validate:"required"`
	Description string   `bson:"description" json:"description" validate:"required"`
	Category    string   `bson:"category" json:"category" validate:"required"`
	Price       int64    `bson:"price" json:"price" validate:"required"`
	Size        []string `bson:"size" json:"size" validate:"required"`
	Color       string   `bson:"color" json:"color" validate:"required"`
	Rating      int64    `bson:"rating" json:"rating" validate:"required"`
	Img         []string `bson:"img" json:"img" validate:"required"`
}
