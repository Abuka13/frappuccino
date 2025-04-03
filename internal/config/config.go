package config

type Config struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
}

// LoadConfig загружает конфигурацию базы данных.
func LoadConfig() *Config {
    return &Config{
        Host:     "db",          // Имя сервиса базы данных из docker-compose.yml
        Port:     "5432",        // Порт PostgreSQL
        User:     "latte",       // Имя пользователя
        Password: "latte",       // Пароль
        Name:     "frappuccino", // Имя базы данных
    }
}