package model

import "time"

type Part struct {
	Uuid          string           `bson:"_id,omitempty"`
	Name          string           `bson:"name"`
	Description   string           `bson:"description"`
	Price         float64          `bson:"price"`
	StockQuantity int64            `bson:"stock_quantity"`
	Category      Category         `bson:"category"`
	Dimensions    Dimensions       `bson:"dimensions"`
	Manufacturer  Manufacturer     `bson:"manufacturer"`
	Tags          *[]string        `bson:"tags,omitempty"`
	Metadata      map[string]Value `bson:"metadata"`
	CreatedAt     *time.Time       `bson:"created_at,omitempty"`
	UpdatedAt     *time.Time       `bson:"updated_at,omitempty"`
}

type Category int

const (
	CategoryUnknown Category = iota
	CategoryEngine
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

type Value struct {
	Str   *string  `bson:"str,omitempty"`
	Int   *int64   `bson:"int,omitempty"`
	Float *float64 `bson:"float,omitempty"`
	Bool  *bool    `bson:"bool,omitempty"`
}

type PartsFilter struct {
	Uuids                 *[]string
	Names                 *[]string
	Categories            *[]Category
	ManufacturerCountries *[]string
	Tags                  *[]string
}
