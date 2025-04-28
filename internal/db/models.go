package db

import (
	"encoding/json"
	"time"
)

type Order struct {
	ID                  int             `json:"id"`
	CustomerID          int             `json:"customer_id"`
	TotalAmount         float64         `json:"total_amount"`
	Status              string          `json:"status"`
	SpecialInstructions json.RawMessage `json:"special_instructions,omitempty"`
	PaymentMethod       string          `json:"payment_method"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

type MenuItem struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Allergens   []string `json:"allergens"`
	Category    string   `json:"category"`
	Size        string   `json:"size"`
}

type Inventory struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Stock       float64   `json:"stock"`
	Price       float64   `json:"price"`
	UnitType    string    `json:"unit_type"`
	LastUpdated time.Time `json:"last_updated"` // default now
}
