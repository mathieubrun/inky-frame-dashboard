package agenda

import "testing"

func TestMockCalendarProvider_GetAgenda(t *testing.T) {
	p := NewMockCalendarProvider()
	count := 5
	forecast, err := p.GetAgenda("test", count)
	if err != nil {
		t.Fatalf("GetAgenda failed: %v", err)
	}

	if len(forecast.Events) != count {
		t.Errorf("Expected %d events, got %d", count, len(forecast.Events))
	}
}
