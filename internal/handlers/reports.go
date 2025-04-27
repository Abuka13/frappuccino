package handlers
import (
	"database/sql"
	"encoding/json"
	"net/http"
	"log"
)
func TotalAmount(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
            http.Error(w, "Invalid Method Request", http.StatusMethodNotAllowed)
            return
        }
        query := `Select SUM(total_amount) from orders`
		row := dbc.QueryRow(query)
		
		var sum float64
		if err := row.Scan(&sum); err != nil {
			http.Error(w, "Failed to scan ordered items", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sum)
	}
}
func PopularItems(dbc *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
            http.Error(w, "Invalid Method Request", http.StatusMethodNotAllowed)
            return
        }
		query := `SELECT 
					mi.id AS menuItemID,
					mi.name,
					mi.description
					
				FROM 
					order_items oi
				JOIN 
					menu_items mi ON oi.menu_item_id = mi.id
				GROUP BY 
					mi.id, mi.name, mi.description`
		rows, err := dbc.Query(query)
		if err != nil {
			http.Error(w, "Failed to fetch ordered items", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer rows.Close()
		type MenuItem struct {
			ID        string `json:"menuItemID"`
			Name     string    `json:"name"`
			Description string `json:"description"`
		}
		var items []MenuItem
		for rows.Next() {
			var menuItemID string
			var name string
			var description string

			if err := rows.Scan(&menuItemID, &name, &description); err != nil {
				http.Error(w, "Failed to scan ordered items", http.StatusInternalServerError)
				log.Println(err)
				return
			}
			item := MenuItem{
				ID:        menuItemID,
				Name:     name,
				Description: description,
			}
			items = append(items, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
}
}