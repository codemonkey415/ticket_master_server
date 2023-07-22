package models

type Women struct {
	Title       string   `bson:"title" json:"title"`
	Gender      string   `bson:"gender" json:"gender"`
	Description string   `bson:"description" json:"description"`
	Category    string   `bson:"category" json:"category"`
	Price       float64  `bson:"price" json:"price"`
	Size        []string `bson:"size" json:"size"`
	Color       string   `bson:"color" json:"color"`
	Rating      float64  `bson:"rating" json:"rating"`
	Img         []string `bson:"img" json:"img"`
}
