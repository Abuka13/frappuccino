### New Endpoints

As part of the task, the following endpoints must be rewritten to work using SQL queries:

    Orders:
        POST /orders: Create a new order.  ++++++++++
        GET /orders: Retrieve all orders. +++++++++++
        GET /orders/{id}: Retrieve a specific order by ID.          ++++++
        PUT /orders/{id}: Update an existing order.    ++++      
        DELETE /orders/{id}: Delete an order.        +++
        POST /orders/{id}/close: Close an order.                 ++

    Menu Items:
        POST /menu: Add a new menu item.
        GET /menu: Retrieve all menu items.
        GET /menu/{id}: Retrieve a specific menu item.
        PUT /menu/{id}: Update a menu item.
        DELETE /menu/{id}: Delete a menu item.

    Inventory:
        POST /inventory: Add a new inventory item.
        GET /inventory: Retrieve all inventory items.
        GET /inventory/{id}: Retrieve a specific inventory item.
        PUT /inventory/{id}: Update an inventory item.
        DELETE /inventory/{id}: Delete an inventory item.

    Aggregations:
        GET /reports/total-sales: Get the total sales amount.
        GET /reports/popular-items: Get a list of popular menu items.



func CreateOrder(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Verify content type
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}
		var order db.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		// Validate required fields
		if order.CustomerID == 0 {
			http.Error(w, "customer_id is required", http.StatusBadRequest)
			return
		}
		if order.TotalAmount <= 0 {
			http.Error(w, "total_amount must be greater than 0", http.StatusBadRequest)
			return
		}
		if order.PaymentMethod == "" {
			http.Error(w, "payment_method is required", http.StatusBadRequest)
			return
		}
		// Set default status if not provided
		if order.Status == "" {
			order.Status = "open"
		}
		// Insert into database
		query := `
            INSERT INTO orders (customer_id, total_amount, status, special_instructions, payment_method)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id
        `
		var orderID int
		err := dbc.QueryRowContext(r.Context(), query,
			order.CustomerID,
			order.TotalAmount,
			order.Status,
			order.SpecialInstructions, // This will be stored as JSONB in PostgreSQL
			order.PaymentMethod,
		).Scan(&orderID)
		if err != nil {
			http.Error(w, "Failed to create order: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"order_id": orderID})
	}
}
