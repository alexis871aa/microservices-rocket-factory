package model

import "time"

type Part struct {
	Uuid          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          *[]string
	Metadata      map[string]Value
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
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
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Value struct {
	Str   *string
	Int   *int64
	Float *float64
	Bool  *bool
}

type PartsFilter struct {
	Uuids                 *[]string
	Names                 *[]string
	Categories            *[]Category
	ManufacturerCountries *[]string
	Tags                  *[]string
}

type PartInfo struct {
	Part Part
}

type PartsInfoFilter struct {
	Parts []Part
}
