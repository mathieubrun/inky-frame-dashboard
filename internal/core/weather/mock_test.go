package weather

import (
	"testing"
	"time"
)

func TestNewMockProvider(t *testing.T) {
	provider := NewMockProvider()
	if provider == nil {
		t.Fatal("expected NewMockProvider to return a non-nil provider")
	}
}

func TestMockProvider_GetForecast(t *testing.T) {
	provider := NewMockProvider()
	city := "Zurich"

	forecast, err := provider.GetForecast(city)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if forecast == nil {
		t.Fatal("expected GetForecast to return a non-nil forecast")
	}

	// Verify Location
	if forecast.Location.City != city {
		t.Errorf("expected city %s, got %s", city, forecast.Location.City)
	}
	if forecast.Location.Country != "CH" {
		t.Errorf("expected country CH, got %s", forecast.Location.Country)
	}

	// Verify Current data
	if forecast.Current.Timestamp.IsZero() {
		t.Error("expected current timestamp to be non-zero")
	}
	if forecast.Current.Temperature < 10.0 || forecast.Current.Temperature > 20.0 {
		t.Errorf("expected current temperature between 10.0 and 20.0, got %f", forecast.Current.Temperature)
	}
	if forecast.Current.WindSpeed < 5.0 || forecast.Current.WindSpeed > 20.0 {
		t.Errorf("expected current wind speed between 5.0 and 20.0, got %f", forecast.Current.WindSpeed)
	}
	if forecast.Current.WindDirection < 0.0 || forecast.Current.WindDirection > 360.0 {
		t.Errorf("expected current wind direction between 0.0 and 360.0, got %f", forecast.Current.WindDirection)
	}
	if forecast.Current.Precipitation < 0.0 || forecast.Current.Precipitation > 5.0 {
		t.Errorf("expected current precipitation between 0.0 and 5.0, got %f", forecast.Current.Precipitation)
	}
	if forecast.Current.PrecipitationProb < 0.0 || forecast.Current.PrecipitationProb > 100.0 {
		t.Errorf("expected current precipitation probability between 0.0 and 100.0, got %f", forecast.Current.PrecipitationProb)
	}

	// Verify Hourly data
	if len(forecast.Hourly) != 24 {
		t.Errorf("expected 24 hourly records, got %d", len(forecast.Hourly))
	}

	// Verify hourly timestamps and ranges
	now := time.Now().Truncate(time.Hour)
	for i, record := range forecast.Hourly {
		expectedTs := now.Add(time.Duration(i) * time.Hour)
		// Allow small difference because of execution time
		if !record.Timestamp.Equal(expectedTs) && !record.Timestamp.After(expectedTs.Add(-time.Second)) && !record.Timestamp.Before(expectedTs.Add(time.Second)) {
			t.Errorf("at index %d: expected timestamp %v, got %v", i, expectedTs, record.Timestamp)
		}

		// Range checks (approximate based on current + delta)
		if record.Temperature < 8.0 || record.Temperature > 22.0 {
			t.Errorf("at index %d: expected temperature between 8.0 and 22.0, got %f", i, record.Temperature)
		}
		if record.WindSpeed < 2.5 || record.WindSpeed > 22.5 {
			t.Errorf("at index %d: expected wind speed between 2.5 and 22.5, got %f", i, record.WindSpeed)
		}
		if record.Precipitation < 0.0 || record.Precipitation > 2.0 {
			t.Errorf("at index %d: expected precipitation between 0.0 and 2.0, got %f", i, record.Precipitation)
		}
		if record.PrecipitationProb < 0.0 || record.PrecipitationProb > 100.0 {
			t.Errorf("at index %d: expected precipitation probability between 0.0 and 100.0, got %f", i, record.PrecipitationProb)
		}
	}

	// Verify FetchedAt
	if forecast.FetchedAt.IsZero() {
		t.Error("expected FetchedAt to be non-zero")
	}
}

func TestMockProvider_InterfaceImplementation(t *testing.T) {
	var _ Provider = (*MockProvider)(nil)
}
