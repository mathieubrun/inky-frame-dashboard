package weather

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeocode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "Zurich" {
			_, _ = fmt.Fprintln(w, `{"results":[{"name":"Zurich","latitude":47.3717,"longitude":8.5422,"country":"Switzerland"}]}`)
		} else {
			_, _ = fmt.Fprintln(w, `{"results":[]}`)
		}
	}))
	defer ts.Close()

	p := NewOpenMeteoProvider()
	p.geocodingURL = ts.URL

	tests := []struct {
		city    string
		wantErr bool
		name    string
	}{
		{"Zurich", false, "Zurich"},
		{"Unknown", true, ""},
	}

	for _, tt := range tests {
		loc, err := p.geocode(tt.city)
		if (err != nil) != tt.wantErr {
			t.Errorf("geocode(%q) error = %v, wantErr %v", tt.city, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && loc.City != tt.name {
			t.Errorf("geocode(%q) = %v, want name %v", tt.city, loc.City, tt.name)
		}
	}
}
