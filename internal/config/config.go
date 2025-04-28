package config

type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
}

// LoadConfig loads the application configuration.
func LoadConfig() *Config {
	return &Config{
		DB: struct {
			Host     string
			Port     string
			User     string
			Password string
			Name     string
		}{
			Host:     "db",          // Matches the service name in docker-compose.yml
			Port:     "5432",        // Default PostgreSQL port
			User:     "latte",       // Database user
			Password: "latte",       // Database password
			Name:     "frappuccino", // Database name
		},
	}
}
