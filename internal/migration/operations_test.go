package migration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amfaro/nocodb-migrator/internal/api"
)

func TestAlterFieldOrderUpdatesDefaultGridViewColumn(t *testing.T) {
	var fieldUpdated bool
	var viewColumnUpdated bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/v3/meta/bases/base1/tables":
			if r.Method != http.MethodGet {
				t.Fatalf("unexpected method for list tables: %s", r.Method)
			}
			_, _ = w.Write([]byte(`{"list":[{"id":"tbl1","title":"T"}]}`))
		case "/api/v3/meta/bases/base1/tables/tbl1":
			if r.Method != http.MethodGet {
				t.Fatalf("unexpected method for get table: %s", r.Method)
			}
			_, _ = w.Write([]byte(`{"id":"tbl1","title":"T","fields":[{"id":"fld1","title":"c","type":"Number"}]}`))
		case "/api/v3/meta/bases/base1/fields/fld1":
			if r.Method != http.MethodPatch {
				t.Fatalf("unexpected method for update field: %s", r.Method)
			}
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode field update: %v", err)
			}
			if body["order"] != 2.5 {
				t.Fatalf("field order = %v, want 2.5", body["order"])
			}
			fieldUpdated = true
			_, _ = w.Write([]byte(`{"id":"fld1","title":"c","type":"Number"}`))
		case "/api/v2/meta/tables/tbl1/views":
			if r.Method != http.MethodGet {
				t.Fatalf("unexpected method for list views: %s", r.Method)
			}
			_, _ = w.Write([]byte(`{"list":[{"id":"view1","title":"Grid view","type":"grid","is_default":true}]}`))
		case "/api/v2/meta/views/view1/columns":
			if r.Method != http.MethodGet {
				t.Fatalf("unexpected method for list view columns: %s", r.Method)
			}
			_, _ = w.Write([]byte(`{"list":[{"id":"vc1","fk_column_id":"fld1","order":1}]}`))
		case "/api/v2/meta/views/view1/columns/vc1":
			if r.Method != http.MethodPatch {
				t.Fatalf("unexpected method for update view column: %s", r.Method)
			}
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode view column update: %v", err)
			}
			if body["order"] != 2.5 {
				t.Fatalf("view column order = %v, want 2.5", body["order"])
			}
			viewColumnUpdated = true
			_, _ = w.Write([]byte(`1`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	order := 2.5
	executor := NewExecutor(api.NewClient(server.URL, "token", "base1"))
	err := executor.ExecuteOperation(&Operation{
		Type:  "alter_field",
		Table: "T",
		Column: &ColumnDefinition{
			Name:  "c",
			Type:  "Number",
			Order: &order,
		},
	})
	if err != nil {
		t.Fatalf("ExecuteOperation returned error: %v", err)
	}
	if !fieldUpdated {
		t.Fatal("field update was not called")
	}
	if !viewColumnUpdated {
		t.Fatal("view column update was not called")
	}
}

func TestSetDisplayFieldByID(t *testing.T) {
	var primaryFieldSet bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/v3/meta/bases/base1/tables":
			if r.Method != http.MethodGet {
				t.Fatalf("unexpected method for list tables: %s", r.Method)
			}
			_, _ = w.Write([]byte(`{"list":[{"id":"tbl1","title":"T"}]}`))
		case "/api/v3/meta/bases/base1/tables/tbl1":
			if r.Method != http.MethodGet {
				t.Fatalf("unexpected method for get table: %s", r.Method)
			}
			_, _ = w.Write([]byte(`{"id":"tbl1","title":"T","fields":[{"id":"fld1","title":"c","type":"Number"}]}`))
		case "/api/v2/meta/columns/fld1/primary":
			if r.Method != http.MethodPost {
				t.Fatalf("unexpected method for set primary field: %s", r.Method)
			}
			primaryFieldSet = true
			_, _ = w.Write([]byte(`{}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	executor := NewExecutor(api.NewClient(server.URL, "token", "base1"))
	err := executor.ExecuteOperation(&Operation{
		Type:    "set_display_field",
		Table:   "T",
		FieldID: "fld1",
	})
	if err != nil {
		t.Fatalf("ExecuteOperation returned error: %v", err)
	}
	if !primaryFieldSet {
		t.Fatal("set primary field was not called")
	}
}

func TestSetDisplayFieldByName(t *testing.T) {
	var primaryFieldSet bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/v3/meta/bases/base1/tables":
			if r.Method != http.MethodGet {
				t.Fatalf("unexpected method for list tables: %s", r.Method)
			}
			_, _ = w.Write([]byte(`{"list":[{"id":"tbl1","title":"T"}]}`))
		case "/api/v3/meta/bases/base1/tables/tbl1":
			if r.Method != http.MethodGet {
				t.Fatalf("unexpected method for get table: %s", r.Method)
			}
			_, _ = w.Write([]byte(`{"id":"tbl1","title":"T","fields":[{"id":"fld1","title":"c","type":"Number"}]}`))
		case "/api/v2/meta/columns/fld1/primary":
			if r.Method != http.MethodPost {
				t.Fatalf("unexpected method for set primary field: %s", r.Method)
			}
			primaryFieldSet = true
			_, _ = w.Write([]byte(`{}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	executor := NewExecutor(api.NewClient(server.URL, "token", "base1"))
	err := executor.ExecuteOperation(&Operation{
		Type:  "set_display_field",
		Table: "T",
		Column: &ColumnDefinition{
			Name: "c",
		},
	})
	if err != nil {
		t.Fatalf("ExecuteOperation returned error: %v", err)
	}
	if !primaryFieldSet {
		t.Fatal("set primary field was not called")
	}
}
