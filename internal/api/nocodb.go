package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client represents a client for working with NocoDB Meta API v3
type Client struct {
	baseURL    string
	apiToken   string
	baseID     string
	httpClient *resty.Client
}

// NewClient creates a new NocoDB client
func NewClient(baseURL, apiToken, baseID string) *Client {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetHeader("xc-token", apiToken)
	client.SetHeader("Content-Type", "application/json")

	return &Client{
		baseURL:    baseURL,
		apiToken:   apiToken,
		baseID:     baseID,
		httpClient: client,
	}
}

// Table represents a table in NocoDB
type Table struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	BaseID      string  `json:"base_id"`
	Fields      []Field `json:"fields,omitempty"`
}

// TableList represents a list of tables
type TableList struct {
	List []Table `json:"list"`
}

// TableCreate represents a table creation request
type TableCreate struct {
	Title       string        `json:"title"`
	Description string        `json:"description,omitempty"`
	Fields      []FieldCreate `json:"fields,omitempty"`
}

// TableUpdate represents a table update request
type TableUpdate struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// Field represents a field in NocoDB
type Field struct {
	ID           string                 `json:"id"`
	Title        string                 `json:"title"`
	Type         string                 `json:"type"`
	DefaultValue interface{}            `json:"default_value,omitempty"`
	Required     bool                   `json:"required,omitempty"`
	Unique       bool                   `json:"unique,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Options      map[string]interface{} `json:"options,omitempty"`
}

// FieldCreate represents a field creation request
type FieldCreate struct {
	Title        string                 `json:"title"`
	Type         string                 `json:"type"`
	DefaultValue interface{}            `json:"default_value,omitempty"`
	Required     bool                   `json:"required,omitempty"`
	Unique       bool                   `json:"unique,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Options      map[string]interface{} `json:"options,omitempty"`
}

// FieldUpdate represents a field update request
type FieldUpdate struct {
	Title        string                 `json:"title,omitempty"`
	Type         string                 `json:"type,omitempty"`
	DefaultValue interface{}            `json:"default_value,omitempty"`
	Required     *bool                  `json:"required,omitempty"`
	Unique       *bool                  `json:"unique,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Options      map[string]interface{} `json:"options,omitempty"`
}

// Record represents a record in a table
type Record map[string]interface{}

// RecordList represents a list of records
type RecordList struct {
	List     []Record `json:"list"`
	PageInfo struct {
		Page      int `json:"page"`
		PageSize  int `json:"pageSize"`
		TotalRows int `json:"totalRows"`
	} `json:"pageInfo"`
}

// APIError represents an API error
type APIError struct {
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Details interface{} `json:"details"`
}

// ListTables gets a list of all tables in the base
func (c *Client) ListTables() (*TableList, error) {
	var result TableList
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("%s/api/v3/meta/bases/%s/tables", c.baseURL, c.baseID))

	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}

	if resp.IsError() {
		var apiErr APIError
		if err := json.Unmarshal(resp.Body(), &apiErr); err == nil {
			return nil, fmt.Errorf("API error: %s, %s", apiErr.Message, apiErr.Error)
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	return &result, nil
}

// GetTable gets a table schema by ID
func (c *Client) GetTable(tableID string) (*Table, error) {
	var result Table
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("%s/api/v3/meta/bases/%s/tables/%s", c.baseURL, c.baseID, tableID))

	if err != nil {
		return nil, fmt.Errorf("failed to get table: %w", err)
	}

	if resp.IsError() {
		var apiErr APIError
		if err := json.Unmarshal(resp.Body(), &apiErr); err == nil {
			return nil, fmt.Errorf("API error: %s", apiErr.Message)
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	return &result, nil
}

// GetTableByName gets a table by name
func (c *Client) GetTableByName(tableName string) (*Table, error) {
	tables, err := c.ListTables()
	if err != nil {
		return nil, err
	}

	for _, table := range tables.List {
		if table.Title == tableName {
			// Get full schema with fields
			return c.GetTable(table.ID)
		}
	}

	return nil, fmt.Errorf("table '%s' not found", tableName)
}

// CreateTable creates a new table
func (c *Client) CreateTable(req *TableCreate) (*Table, error) {
	var result Table
	resp, err := c.httpClient.R().
		SetBody(req).
		SetResult(&result).
		Post(fmt.Sprintf("%s/api/v3/meta/bases/%s/tables", c.baseURL, c.baseID))

	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	if resp.IsError() {
		var apiErr APIError
		fmt.Println(string(resp.Body()))
		if err := json.Unmarshal(resp.Body(), &apiErr); err == nil {
			return nil, fmt.Errorf("API error: %s", apiErr.Message)
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	return &result, nil
}

// UpdateTable updates a table
func (c *Client) UpdateTable(tableID string, req *TableUpdate) (*Table, error) {
	var result Table
	resp, err := c.httpClient.R().
		SetBody(req).
		SetResult(&result).
		Patch(fmt.Sprintf("%s/api/v3/meta/bases/%s/tables/%s", c.baseURL, c.baseID, tableID))

	if err != nil {
		return nil, fmt.Errorf("failed to update table: %w", err)
	}

	if resp.IsError() {
		var apiErr APIError
		if err := json.Unmarshal(resp.Body(), &apiErr); err == nil {
			return nil, fmt.Errorf("API error: %s", apiErr.Message)
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	return &result, nil
}

// DeleteTable deletes a table
func (c *Client) DeleteTable(tableID string) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("%s/api/v3/meta/bases/%s/tables/%s", c.baseURL, c.baseID, tableID))

	if err != nil {
		return fmt.Errorf("failed to delete table: %w", err)
	}

	if resp.IsError() {
		var apiErr APIError
		if err := json.Unmarshal(resp.Body(), &apiErr); err == nil {
			return fmt.Errorf("API error: %s", apiErr.Message)
		}
		return fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	return nil
}

// CreateField creates a new field in a table
func (c *Client) CreateField(tableID string, req *FieldCreate) (*Field, error) {
	var result Field
	resp, err := c.httpClient.R().
		SetBody(req).
		SetResult(&result).
		Post(fmt.Sprintf("%s/api/v3/meta/bases/%s/tables/%s/fields", c.baseURL, c.baseID, tableID))

	if err != nil {
		return nil, fmt.Errorf("failed to create field: %w", err)
	}

	if resp.IsError() {
		var apiErr APIError
		if err := json.Unmarshal(resp.Body(), &apiErr); err == nil {
			return nil, fmt.Errorf("API error: %s", apiErr.Message)
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	return &result, nil
}

// GetField gets a field by ID
func (c *Client) GetField(fieldID string) (*Field, error) {
	var result Field
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("%s/api/v3/meta/bases/%s/fields/%s", c.baseURL, c.baseID, fieldID))

	if err != nil {
		return nil, fmt.Errorf("failed to get field: %w", err)
	}

	if resp.IsError() {
		var apiErr APIError
		if err := json.Unmarshal(resp.Body(), &apiErr); err == nil {
			return nil, fmt.Errorf("API error: %s", apiErr.Message)
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	return &result, nil
}

// UpdateField updates a field
func (c *Client) UpdateField(fieldID string, req *FieldUpdate) (*Field, error) {
	var result Field
	resp, err := c.httpClient.R().
		SetBody(req).
		SetResult(&result).
		Patch(fmt.Sprintf("%s/api/v3/meta/bases/%s/fields/%s", c.baseURL, c.baseID, fieldID))

	if err != nil {
		return nil, fmt.Errorf("failed to update field: %w", err)
	}

	if resp.IsError() {
		var apiErr APIError
		if err := json.Unmarshal(resp.Body(), &apiErr); err == nil {
			return nil, fmt.Errorf("API error: %s", apiErr.Message)
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	return &result, nil
}

// DeleteField deletes a field
func (c *Client) DeleteField(fieldID string) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("%s/api/v3/meta/bases/%s/fields/%s", c.baseURL, c.baseID, fieldID))

	if err != nil {
		return fmt.Errorf("failed to delete field: %w", err)
	}

	if resp.IsError() {
		var apiErr APIError
		if err := json.Unmarshal(resp.Body(), &apiErr); err == nil {
			return fmt.Errorf("API error: %s", apiErr.Message)
		}
		return fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	return nil
}

// InsertRecord inserts a record into a table
func (c *Client) InsertRecord(tableID string, record Record) (Record, error) {
	var result struct {
		ID     interface{}            `json:"id"`
		Fields map[string]interface{} `json:"fields"`
	}

	resp, err := c.httpClient.R().
		SetBody(map[string]interface{}{"fields": record}).
		SetResult(&result).
		Post(fmt.Sprintf("%s/api/v3/data/%s/%s/records", c.baseURL, c.baseID, tableID))

	if err != nil {
		return nil, fmt.Errorf("failed to insert record: %w", err)
	}

	if resp.IsError() {
		var apiErr APIError
		if err := json.Unmarshal(resp.Body(), &apiErr); err == nil {
			return nil, fmt.Errorf("API error: %s", apiErr.Message)
		}
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode(), string(resp.Body()))
	}

	// Return full record with ID
	resultRecord := make(Record)
	resultRecord["Id"] = result.ID
	resultRecord["id"] = result.ID
	for k, v := range result.Fields {
		resultRecord[k] = v
	}
	return resultRecord, nil
}

// DeleteRecord deletes a record from a table by ID
// In NocoDB v3, deletion is performed via DELETE with body containing an array of objects with id
// Body format: [{"id": "recordID"}]
func (c *Client) DeleteRecord(tableID string, recordID string) error {
	// Convert recordID to the required type (can be string or number)
	var idValue interface{} = recordID

	// Try to convert to number if possible
	if idInt, err := strconv.Atoi(recordID); err == nil {
		idValue = idInt
	}

	// Form body for deletion: array of objects with id
	deleteBody := []map[string]interface{}{
		{"id": idValue},
	}

	resp, err := c.httpClient.R().
		SetBody(deleteBody).
		Delete(fmt.Sprintf("%s/api/v3/data/%s/%s/records", c.baseURL, c.baseID, tableID))

	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	if resp.IsError() {
		var apiErr APIError
		if err := json.Unmarshal(resp.Body(), &apiErr); err == nil {
			return fmt.Errorf("API error: %s", apiErr.Message)
		}
		return fmt.Errorf("API error: status %d, body: %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

// DeleteRecords deletes records by array of IDs
// In NocoDB v3, bulk deletion is performed via DELETE with body containing an array of objects with id
// Body format: [{"id": 1}, {"id": 2}]
func (c *Client) DeleteRecords(tableID string, where map[string]interface{}) error {
	// If where is provided, first get records by condition
	var recordIDs []interface{}

	if len(where) > 0 {
		// Get records by condition via GET request
		req := c.httpClient.R()

		// Convert where to query parameter for GET request
		whereJSON, err := json.Marshal(where)
		if err != nil {
			return fmt.Errorf("failed to marshal where clause: %w", err)
		}
		req.SetQueryParam("where", string(whereJSON))

		var result struct {
			Records []struct {
				ID     interface{}            `json:"id"`
				Fields map[string]interface{} `json:"fields"`
			} `json:"records"`
			Next string `json:"next,omitempty"`
		}

		req.SetResult(&result)
		resp, err := req.Get(fmt.Sprintf("%s/api/v3/data/%s/%s/records", c.baseURL, c.baseID, tableID))

		if err != nil {
			return fmt.Errorf("failed to get records for deletion: %w", err)
		}

		if resp.IsError() {
			var apiErr APIError
			if err := json.Unmarshal(resp.Body(), &apiErr); err == nil {
				return fmt.Errorf("API error getting records: %s", apiErr.Message)
			}
			return fmt.Errorf("API error: status %d, body: %s", resp.StatusCode(), string(resp.Body()))
		}

		// Collect record IDs
		for _, record := range result.Records {
			recordIDs = append(recordIDs, record.ID)
		}
	} else {
		return fmt.Errorf("where condition is required for DeleteRecords")
	}

	if len(recordIDs) == 0 {
		return nil // No records to delete
	}

	// Form body for bulk deletion: array of objects with id
	deleteBody := make([]map[string]interface{}, len(recordIDs))
	for i, id := range recordIDs {
		deleteBody[i] = map[string]interface{}{"id": id}
	}

	// Execute DELETE request with array of IDs in body
	resp, err := c.httpClient.R().
		SetBody(deleteBody).
		Delete(fmt.Sprintf("%s/api/v3/data/%s/%s/records", c.baseURL, c.baseID, tableID))

	if err != nil {
		return fmt.Errorf("failed to delete records: %w", err)
	}

	if resp.IsError() {
		var apiErr APIError
		if err := json.Unmarshal(resp.Body(), &apiErr); err == nil {
			return fmt.Errorf("API error: %s", apiErr.Message)
		}
		return fmt.Errorf("API error: status %d, body: %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

// GetRecords gets records from a table
func (c *Client) GetRecords(tableID string, limit, offset int) (*RecordList, error) {
	var result struct {
		Records []struct {
			ID     interface{}            `json:"id"`
			Fields map[string]interface{} `json:"fields"`
		} `json:"records"`
		Next string `json:"next,omitempty"`
	}

	req := c.httpClient.R().
		SetResult(&result)

	if limit > 0 {
		req.SetQueryParam("limit", fmt.Sprintf("%d", limit))
	}
	if offset > 0 {
		req.SetQueryParam("offset", fmt.Sprintf("%d", offset))
	}

	resp, err := req.Get(fmt.Sprintf("%s/api/v3/data/%s/%s/records", c.baseURL, c.baseID, tableID))

	if err != nil {
		return nil, fmt.Errorf("failed to get records: %w", err)
	}

	if resp.IsError() {
		var apiErr APIError
		if err := json.Unmarshal(resp.Body(), &apiErr); err == nil {
			return nil, fmt.Errorf("API error: %s", apiErr.Message)
		}
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode(), string(resp.Body()))
	}

	// Convert to RecordList format
	recordList := &RecordList{
		List: make([]Record, len(result.Records)),
	}

	for i, item := range result.Records {
		record := make(Record)
		// Add ID in both formats for compatibility
		record["Id"] = item.ID
		record["id"] = item.ID
		for k, v := range item.Fields {
			record[k] = v
		}
		recordList.List[i] = record
	}

	return recordList, nil
}
