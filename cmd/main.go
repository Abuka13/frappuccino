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
    http.HandleFunc("POST /orders/close/", handlers.CloseOrder(dbConn))
    http.HandleFunc("GET /orders/numberOfOrderedItems", handlers.GetNumberOfOrderedItems(dbConn))

    // Запускаем HTTP-сервер
    log.Println("Server is running on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}