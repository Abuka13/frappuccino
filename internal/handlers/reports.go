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
// func PopularItems(dbc *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodGet {
//             http.Error(w, "Invalid Method Request", http.StatusMethodNotAllowed)
//             return
//         }
// 		query
// 	}
// }