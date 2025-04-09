package db

import (
    "time"
    "encoding/json"
)

type Order struct {
    ID                 int             `json:"id"`
    CustomerID         int             `json:"customer_id"`
    TotalAmount        float64         `json:"total_amount"`
    Status             string          `json:"status"`
    SpecialInstructions json.RawMessage `json:"special_instructions,omitempty"`
    PaymentMethod      string          `json:"payment_method"`
    CreatedAt          time.Time       `json:"created_at"`
    UpdatedAt          time.Time       `json:"updated_at"`
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
