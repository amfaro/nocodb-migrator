package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// NewCreateCommand creates the create command
func NewCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new migration",
		Long:  "Create a new migration with up and down files",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			return createMigration(name)
		},
	}

	return cmd
}

func createMigration(name string) error {
	// Get migrations directory from environment variable or use default value
	migrationsDir := getMigrationsDir()

	// If path is relative, make it absolute relative to current directory
	if !filepath.IsAbs(migrationsDir) {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		migrationsDir = filepath.Join(cwd, migrationsDir)
	}

	// Create migrations directory if it doesn't exist
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Generate timestamp
	timestamp := time.Now().Unix()

	// Create file names
	upFileName := fmt.Sprintf("%d-%s.up.json", timestamp, name)
	downFileName := fmt.Sprintf("%d-%s.down.json", timestamp, name)

	upPath := filepath.Join(migrationsDir, upFileName)
	downPath := filepath.Join(migrationsDir, downFileName)

	// Create file contents
	upContent := `{
  "operations": [
    {
      "type": "create_table",
      "table": "ExampleTable",
      "columns": [
        {
          "name": "Id",
          "type": "ID",
          "required": true
        },
        {
          "name": "Name",
          "type": "SingleLineText",
          "required": true
        }
      ]
    }
  ]
}
`

	downContent := `{
  "operations": [
    {
      "type": "drop_table",
      "table": "ExampleTable"
    }
  ]
}
`

	// Write files
	if err := os.WriteFile(upPath, []byte(upContent), 0644); err != nil {
		return fmt.Errorf("failed to create up migration file: %w", err)
	}

	if err := os.WriteFile(downPath, []byte(downContent), 0644); err != nil {
		return fmt.Errorf("failed to create down migration file: %w", err)
	}

	fmt.Printf("Created migration files:\n")
	fmt.Printf("  Up:   %s\n", upPath)
	fmt.Printf("  Down: %s\n", downPath)

	return nil
}
