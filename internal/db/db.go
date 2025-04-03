package db

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
    "frappuccino/internal/config" // Импортируем пакет config
)

// Connect устанавливает соединение с базой данных.
func Connect() (*sql.DB, error) {
    // Загружаем конфигурацию базы данных
    cfg := config.LoadConfig()

    // Формируем строку подключения
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
    )

    // Открываем соединение с базой данных
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }

    // Проверяем соединение
    if err := db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}