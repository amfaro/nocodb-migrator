package cmd

import (
	"fmt"
	"os"

	"github.com/amfaro/nocodb-migrator/internal/api"
)

func initClient() (*api.Client, error) {
	baseURL := os.Getenv("NOCODB_URL")
	apiToken := os.Getenv("NOCODB_API_TOKEN")
	baseID := os.Getenv("NOCODB_BASE_ID")

	if baseURL == "" {
		return nil, fmt.Errorf("NOCODB_URL environment variable is required")
	}
	if apiToken == "" {
		return nil, fmt.Errorf("NOCODB_API_TOKEN environment variable is required")
	}
	if baseID == "" {
		return nil, fmt.Errorf("NOCODB_BASE_ID environment variable is required")
	}

	return api.NewClient(baseURL, apiToken, baseID), nil
}

// getMigrationsDir returns the migrations directory from environment variable
// or default value "./migrations"
func getMigrationsDir() string {
	dir := os.Getenv("NOCODB_MIGRATIONS_DIR")
	if dir == "" {
		return "./migrations"
	}
	return dir
}
