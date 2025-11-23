package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all server configuration settings
type Config struct {
	UploadDir         string
	UploadPassword    string
	MaxFileSizeMB     int
	FileExpiryMinutes int
	Port              string
	IsDefaultPassword bool
}

// LoadFromEnv loads configuration from environment variables with sensible defaults
func LoadFromEnv() (*Config, error) {
	uploadPassword := os.Getenv("UPLOAD_PASSWORD")
	isDefaultPassword := uploadPassword == "" || uploadPassword == "demo"
	if uploadPassword == "" {
		uploadPassword = "demo"
	}

	cfg := &Config{
		UploadDir:         "./uploads",
		UploadPassword:    uploadPassword,
		MaxFileSizeMB:     getIntEnvOrDefault("MAX_FILE_SIZE_MB", 100),
		FileExpiryMinutes: getIntEnvOrDefault("FILE_EXPIRY_MINUTES", 10),
		Port:              getEnvOrDefault("PORT", "8088"),
		IsDefaultPassword: isDefaultPassword,
	}

	if cfg.Port != "" && cfg.Port[0] != ':' {
		cfg.Port = ":" + cfg.Port
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Validate checks that all configuration values are valid
func (c *Config) Validate() error {
	if c.UploadPassword == "" {
		return fmt.Errorf("upload password cannot be empty")
	}
	if c.MaxFileSizeMB <= 0 {
		return fmt.Errorf("max file size must be positive, got %d", c.MaxFileSizeMB)
	}
	if c.FileExpiryMinutes <= 0 {
		return fmt.Errorf("file expiry minutes must be positive, got %d", c.FileExpiryMinutes)
	}
	return nil
}

// MaxFileBytes returns the maximum file size in bytes
func (c *Config) MaxFileBytes() int64 {
	return int64(c.MaxFileSizeMB) << 20 // Convert MB to bytes
}

// LogSummary logs the configuration (without sensitive data)
func (c *Config) LogSummary() []string {
	return []string{
		fmt.Sprintf("Password: %s", "***"), // Don't log actual password
		fmt.Sprintf("Max file size: %d MB", c.MaxFileSizeMB),
		fmt.Sprintf("File expiry: %d minutes", c.FileExpiryMinutes),
		fmt.Sprintf("Port: %s", c.Port),
		fmt.Sprintf("Upload directory: %s", c.UploadDir),
	}
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnvOrDefault returns environment variable as integer or default if not set/invalid
func getIntEnvOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil && intValue > 0 {
			return intValue
		}
	}
	return defaultValue
}
