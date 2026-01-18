package migration

import (
	"fmt"
	"log"
)

// ExecuteMigration executes a migration
func (e *Executor) ExecuteMigration(m *Migration) error {
	log.Printf("Executing migration with %d operations", len(m.Operations))

	for i, op := range m.Operations {
		log.Printf("Executing operation %d/%d: %s", i+1, len(m.Operations), op.Type)

		if err := e.ExecuteOperation(&op); err != nil {
			return fmt.Errorf("operation %d failed: %w", i+1, err)
		}

		log.Printf("Operation %d completed successfully", i+1)
	}

	log.Printf("Migration completed successfully")
	return nil
}
