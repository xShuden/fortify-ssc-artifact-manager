package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	SSCUrl   string
	SSCToken string
}

func LoadConfig() (*Config, error) {
	// Try to load from .env file first
	godotenv.Load(".env")
	
	// Try to load from ssc_env.txt if exists
	godotenv.Load("ssc_env.txt")

	config := &Config{
		SSCUrl:   os.Getenv("FORTIFY_SSC_URL"),
		SSCToken: os.Getenv("FORTIFY_SSC_TOKEN"),
	}

	// Validate configuration
	if config.SSCUrl == "" {
		return nil, fmt.Errorf("FORTIFY_SSC_URL environment variable is not set")
	}
	
	if config.SSCToken == "" {
		return nil, fmt.Errorf("FORTIFY_SSC_TOKEN environment variable is not set")
	}

	// Ensure URL ends without trailing slash for consistency
	config.SSCUrl = strings.TrimRight(config.SSCUrl, "/")

	return config, nil
}

func LoadConfigWithOverrides(urlOverride, tokenOverride string) (*Config, error) {
	// First load from environment
	config, _ := LoadConfig()
	
	// If config is nil, create new one
	if config == nil {
		config = &Config{}
	}

	// Override with command line parameters if provided
	if urlOverride != "" {
		config.SSCUrl = strings.TrimRight(urlOverride, "/")
	}
	
	if tokenOverride != "" {
		config.SSCToken = tokenOverride
	}

	// Validate configuration
	if config.SSCUrl == "" {
		return nil, fmt.Errorf("SSC URL is required. Use --url flag or set FORTIFY_SSC_URL environment variable")
	}
	
	if config.SSCToken == "" {
		return nil, fmt.Errorf("SSC Token is required. Use --token flag or set FORTIFY_SSC_TOKEN environment variable")
	}

	return config, nil
}