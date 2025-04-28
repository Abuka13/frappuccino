package main

import (
    "log"
    "net/http"

    "frappuccino/internal/config"
    "frappuccino/internal/db"
    "frappuccino/internal/handlers"
)

func main() {
    // Загружаем конфигурацию базы данных
    cfg := config.LoadConfig()

    // Подключаемся к базе данных
    dbConn, err := db.Connect(cfg.DB)
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
    }
    defer dbConn.Close()

    // Регистрируем обработчики
    http.HandleFunc("GET /orders", handlers.GetOrders(dbConn))
    http.HandleFunc("POST /orders", handlers.CreateOrder(dbConn))
    http.HandleFunc("DELETE /orders/", handlers.DeleteOrder(dbConn))
    http.HandleFunc("GET /orders/", handlers.GetOrderByID(dbConn))
    http.HandleFunc("PUT /orders/", handlers.UpdateOrderByID(dbConn))
    http.HandleFunc("POST /orders/close/", handlers.CloseOrder(dbConn))go install mvdan.cc/gofumpt@latest

    // Inventory routes
    http.HandleFunc("GET /inventory", handlers.GetInventoryItems(dbConn))
    http.HandleFunc("POST /inventory", handlers.CreateInventoryItem(dbConn))
    http.HandleFunc("GET /inventory/", handlers.GetInventoryItemByID(dbConn))
    http.HandleFunc("PUT /inventory/", handlers.UpdateInventoryItem(dbConn))
    http.HandleFunc("DELETE /inventory/", handlers.DeleteInventoryItem(dbConn))
    // http.HandleFunc("POST /inventory/restock/", handlers.RestockInventoryItem(dbConn))

    // Menu Items routes
    http.HandleFunc("GET /menu", handlers.GetMenuItems(dbConn))
    http.HandleFunc("POST /menu", handlers.CreateMenuItem(dbConn))
    http.HandleFunc("GET /menu/", handlers.GetMenuItemByID(dbConn))
    http.HandleFunc("PUT /menu/", handlers.UpdateMenuItem(dbConn))
    http.HandleFunc("DELETE /menu/", handlers.DeleteMenuItem(dbConn))
    // http.HandleFunc("POST /menu_items/toggle/", handlers.ToggleMenuItemAvailability(dbConn))


    // Report routes
    http.HandleFunc("GET /reports/total-sales", handlers.TotalAmount(dbConn))
    http.HandleFunc("GET /reports/popular-items", handlers.PopularItems(dbConn))

    http.HandleFunc("GET /orders/numberOfOrderedItems", handlers.GetNumberOfOrderedItems(dbConn))

    http.HandleFunc("GET /reports/search", handlers.FullTextSearchReport(dbConn))
    http.HandleFunc("GET /reports/orderedItemsByPeriod", handlers.OrderedItemsByPeriod(dbConn))
    http.HandleFunc("POST /orders/batch-process", handlers.BulkOrderProcess(dbConn))
    http.HandleFunc("GET /inventory/getLeftOvers", handlers.GetLeftovers(dbConn))
        // Запускаем HTTP-сервер    
    log.Println("Server is running on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}