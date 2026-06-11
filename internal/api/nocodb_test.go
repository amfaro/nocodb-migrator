package api

import (
	"encoding/json"
	"testing"
)

func TestViewUnmarshalAcceptsBoolishIsDefault(t *testing.T) {
	tests := []struct {
		name string
		json string
		want bool
	}{
		{name: "bool true", json: `{"is_default":true}`, want: true},
		{name: "bool false", json: `{"is_default":false}`, want: false},
		{name: "numeric one", json: `{"is_default":1}`, want: true},
		{name: "numeric zero", json: `{"is_default":0}`, want: false},
		{name: "float one", json: `{"is_default":1.0}`, want: true},
		{name: "float zero", json: `{"is_default":0.0}`, want: false},
		{name: "string one", json: `{"is_default":"1"}`, want: true},
		{name: "string false", json: `{"is_default":"false"}`, want: false},
		{name: "null", json: `{"is_default":null}`, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var view View
			if err := json.Unmarshal([]byte(tt.json), &view); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if got := bool(view.IsDefault); got != tt.want {
				t.Fatalf("IsDefault = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestViewUnmarshalAcceptsNumericType(t *testing.T) {
	var view View
	if err := json.Unmarshal([]byte(`{"type":3}`), &view); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if view.Type != viewKindGrid {
		t.Fatalf("Type = %q, want %q", view.Type, viewKindGrid)
	}
}
