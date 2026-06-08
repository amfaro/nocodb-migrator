package migration

import (
	"fmt"

	"github.com/memclutter/nocodb-migrator/internal/api"
)

// Executor executes migration operations
type Executor struct {
	client *api.Client
}

// NewExecutor creates a new migration executor
func NewExecutor(client *api.Client) *Executor {
	return &Executor{
		client: client,
	}
}

// ExecuteOperation executes a single migration operation
func (e *Executor) ExecuteOperation(op *Operation) error {
	switch op.Type {
	case "create_table":
		return e.createTable(op)
	case "alter_table":
		return e.alterTable(op)
	case "drop_table":
		return e.dropTable(op)
	case "create_field":
		return e.createField(op)
	case "alter_field":
		return e.alterField(op)
	case "drop_field":
		return e.dropField(op)
	case "insert_row":
		return e.insertRow(op)
	case "delete_row":
		return e.deleteRow(op)
	default:
		return fmt.Errorf("unknown operation type: %s", op.Type)
	}
}

// createTable creates a table
func (e *Executor) createTable(op *Operation) error {
	// Convert columns to API format
	fields := make([]api.FieldCreate, len(op.Columns))
	for i, col := range op.Columns {
		fields[i] = api.FieldCreate{
			Title:        col.Name,
			Type:         col.Type,
			DefaultValue: col.DefaultValue,
			Required:     col.Required,
			Unique:       col.Unique,
			Description:  col.Description,
			Options:      col.Options,
			Order:        col.Order,
		}
	}

	req := &api.TableCreate{
		Title:       op.Table,
		Description: "",
		Fields:      fields,
	}

	_, err := e.client.CreateTable(req)
	return err
}

// alterTable modifies a table
func (e *Executor) alterTable(op *Operation) error {
	// Get table by name
	table, err := e.client.GetTableByName(op.Table)
	if err != nil {
		return fmt.Errorf("failed to get table: %w", err)
	}

	req := &api.TableUpdate{}

	// Update only specified fields from data
	if op.Data != nil {
		if title, ok := op.Data["title"].(string); ok {
			req.Title = title
		}
		if desc, ok := op.Data["description"].(string); ok {
			req.Description = desc
		}
	}

	_, err = e.client.UpdateTable(table.ID, req)
	return err
}

// dropTable deletes a table
func (e *Executor) dropTable(op *Operation) error {
	table, err := e.client.GetTableByName(op.Table)
	if err != nil {
		return fmt.Errorf("failed to get table: %w", err)
	}

	return e.client.DeleteTable(table.ID)
}

// createField creates a field in a table
func (e *Executor) createField(op *Operation) error {
	table, err := e.client.GetTableByName(op.Table)
	if err != nil {
		return fmt.Errorf("failed to get table: %w", err)
	}

	// For LinkToAnotherRecord, need to convert relatedTable (name) to related_table_id (ID)
	// and ensure relation_type exists
	options := op.Column.Options
	if op.Column.Type == "LinkToAnotherRecord" && options != nil {
		newOptions := make(map[string]interface{})

		// Copy all existing options except relatedTable
		for k, v := range options {
			if k != "relatedTable" {
				newOptions[k] = v
			}
		}

		// Convert relatedTable (name) to related_table_id (ID)
		if relatedTableName, ok := options["relatedTable"].(string); ok {
			relatedTable, err := e.client.GetTableByName(relatedTableName)
			if err != nil {
				return fmt.Errorf("failed to get related table %s: %w", relatedTableName, err)
			}
			newOptions["related_table_id"] = relatedTable.ID
		}

		// Check for relation_type (required parameter)
		if _, exists := newOptions["relation_type"]; !exists {
			// Default to "hm" (has-many) if not specified
			newOptions["relation_type"] = "hm"
		}

		options = newOptions
	}

	req := &api.FieldCreate{
		Title:        op.Column.Name,
		Type:         op.Column.Type,
		DefaultValue: op.Column.DefaultValue,
		Required:     op.Column.Required,
		Unique:       op.Column.Unique,
		Description:  op.Column.Description,
		Options:      options,
		Order:        op.Column.Order,
	}

	_, err = e.client.CreateField(table.ID, req)
	return err
}

// alterField modifies a field
func (e *Executor) alterField(op *Operation) error {
	table, err := e.client.GetTableByName(op.Table)
	if err != nil {
		return fmt.Errorf("failed to get table: %w", err)
	}

	var fieldID string
	if op.FieldID != "" {
		fieldID = op.FieldID
	} else if op.Column != nil && op.Column.Name != "" {
		// Find field by name
		for _, field := range table.Fields {
			if field.Title == op.Column.Name {
				fieldID = field.ID
				break
			}
		}
		if fieldID == "" {
			return fmt.Errorf("field '%s' not found in table '%s'", op.Column.Name, op.Table)
		}
	} else {
		return fmt.Errorf("field_id or column.name is required")
	}

	req := &api.FieldUpdate{}
	if op.Column != nil {
		if op.Column.Name != "" {
			req.Title = op.Column.Name
		}
		if op.Column.Type != "" {
			req.Type = op.Column.Type
		}
		if op.Column.DefaultValue != nil {
			req.DefaultValue = op.Column.DefaultValue
		}
		if op.Column.Description != "" {
			req.Description = op.Column.Description
		}
		if op.Column.Options != nil {
			req.Options = op.Column.Options
		}
		if op.Column.Order != nil {
			req.Order = op.Column.Order
		}
		req.Required = &op.Column.Required
		req.Unique = &op.Column.Unique
	}

	_, err = e.client.UpdateField(fieldID, req)
	return err
}

// dropField deletes a field
func (e *Executor) dropField(op *Operation) error {
	table, err := e.client.GetTableByName(op.Table)
	if err != nil {
		return fmt.Errorf("failed to get table: %w", err)
	}

	var fieldID string
	if op.FieldID != "" {
		fieldID = op.FieldID
	} else if op.Column != nil && op.Column.Name != "" {
		// Find field by name
		for _, field := range table.Fields {
			if field.Title == op.Column.Name {
				fieldID = field.ID
				break
			}
		}
		if fieldID == "" {
			return fmt.Errorf("field '%s' not found in table '%s'", op.Column.Name, op.Table)
		}
	} else {
		return fmt.Errorf("field_id or column.name is required")
	}

	return e.client.DeleteField(fieldID)
}

// insertRow inserts a record into a table
func (e *Executor) insertRow(op *Operation) error {
	table, err := e.client.GetTableByName(op.Table)
	if err != nil {
		return fmt.Errorf("failed to get table: %w", err)
	}

	_, err = e.client.InsertRecord(table.ID, op.Data)
	return err
}

// deleteRow deletes a record from a table
func (e *Executor) deleteRow(op *Operation) error {
	table, err := e.client.GetTableByName(op.Table)
	if err != nil {
		return fmt.Errorf("failed to get table: %w", err)
	}

	if op.RecordID != "" {
		return e.client.DeleteRecord(table.ID, op.RecordID)
	}

	if len(op.Where) > 0 {
		return e.client.DeleteRecords(table.ID, op.Where)
	}

	return fmt.Errorf("either record_id or where condition is required")
}
