package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Database DatabaseConfig `json:"database"`
	Stripe   StripeConfig  `json:"stripe"`
	Xendit   XenditConfig  `json:"xendit"`
	Server   ServerConfig  `json:"server"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

type StripeConfig struct {
	Secret string `json:"secret"`
	Public string `json:"public"`
}

type XenditConfig struct {
	Secret string `json:"secret"`
	Public string `json:"public"`
}

type ServerConfig struct {
	Port string `json:"port"`
}

// LoadConfig loads configuration from a JSON file and environment variables
func LoadConfig() (*Config, error) {
	config := &Config{}
	
	// Get the absolute path to the config directory
	configDir, err := filepath.Abs("config")
	if err != nil {
		return nil, err
	}
	
	configFile := filepath.Join(configDir, "config.json")
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	// Override with environment variables if present
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		fmt.Sscanf(dbPort, "%d", &config.Database.Port)
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		config.Database.User = dbUser
	}
	if dbPass := os.Getenv("DB_PASSWORD"); dbPass != "" {
		config.Database.Password = dbPass
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.Database.DBName = dbName
	}
	if dbSSLMode := os.Getenv("DB_SSLMODE"); dbSSLMode != "" {
		config.Database.SSLMode = dbSSLMode
	}
	if stripeKey := os.Getenv("STRIPE_API_KEY"); stripeKey != "" {
		config.Stripe.Secret = stripeKey
	}
	if xenditKey := os.Getenv("XENDIT_API_KEY"); xenditKey != "" {
		config.Xendit.Secret = xenditKey
	}
	if port := os.Getenv("PORT"); port != "" {
		config.Server.Port = port
	}

	// Set defaults if not configured
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.Database.Port == 0 {
		config.Database.Port = 5432
	}
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}

	return config, nil
}

func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}
