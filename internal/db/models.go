package db

import "time"

type Order struct {
    ID                int
    CustomerID        int
    TotalAmount       float64
    Status            string
    SpecialInstructions map[string]string
    PaymentMethod     string
    CreatedAt         time.Time
    UpdatedAt         time.Time
}

type MenuItem struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Allergens   []string
	Category    string
	Size        string
}
