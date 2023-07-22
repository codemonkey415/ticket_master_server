package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderSummary struct {
	SubTotal int64 `bson:"subTotal" json:"subTotal" validate:"required"`
	Quantity int64 `bson:"quantity" json:"quantity" validate:"required"`
	Shipping int64 `bson:"shipping" json:"shipping" validate:"required"`
	Discount int64 `bson:"discount" json:"discount" validate:"required"`
	Total    int64 `bson:"total" json:"total" validate:"required"`
}

type CartProduct struct {
	Title       string             `bson:"title" json:"title" validate:"required"`
	Gender      string             `bson:"gender" json:"gender" validate:"required"`
	Description string             `bson:"description" json:"description" validate:"required"`
	Category    string             `bson:"category" json:"category" validate:"required"`
	Price       int64              `bson:"price" json:"price" validate:"required"`
	Size        string             `bson:"size" json:"size" validate:"required"`
	Color       string             `bson:"color" json:"color" validate:"required"`
	Rating      int64              `bson:"rating" json:"rating" validate:"required"`
	Img         []primitive.Binary `bson:"img" json:"img" validate:"required"`
	Quantity    int64              `bson:"quantity" json:"quantity" validate:"required"`
}

type PaymentDetails struct {
	OrderID           string `bson:"orderId" json:"orderId" validate:"required"`
	RazorpayOrderID   string `bson:"razorpayOrderId" json:"razorpayOrderId" validate:"required"`
	RazorpayPaymentID string `bson:"razorpayPaymentId" json:"razorpayPaymentId" validate:"required"`
}

type ShippingDetails struct {
	FirstName    string `bson:"firstName" json:"firstName" validate:"required"`
	LastName     string `bson:"lastName" json:"lastName" validate:"required"`
	AddressLine1 string `bson:"addressLine1" json:"addressLine1" validate:"required"`
	AddressLine2 string `bson:"addressLine2,omitempty" json:"addressLine2,omitempty"`
	Locality     string `bson:"locality" json:"locality" validate:"required"`
	PinCode      int64  `bson:"pinCode" json:"pinCode" validate:"required"`
	State        string `bson:"state" json:"state" validate:"required"`
	Country      string `bson:"country" json:"country" validate:"required"`
	Email        string `bson:"email" json:"email" validate:"required,email"`
	Mobile       int64  `bson:"mobile" json:"mobile" validate:"required,len=10"`
}

type Order struct {
	OrderSummary    OrderSummary       `bson:"orderSummary" json:"orderSummary"`
	CartProducts    []CartProduct      `bson:"cartProducts" json:"cartProducts"`
	PaymentDetails  PaymentDetails     `bson:"paymentDetails,omitempty" json:"-"`
	ShippingDetails ShippingDetails    `bson:"shippingDetails,omitempty" json:"-"`
	User            primitive.ObjectID `bson:"user,omitempty" json:"-"`
}
