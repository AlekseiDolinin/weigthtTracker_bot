package parse

import (
	"testing"
	"time"
	"weightTrack_bot/models"
)

func TestParseRecord(t *testing.T) {
	value := "5218837036 58.00 2025-07-09T21:56:18+04:00 0"

	expectedTime, err := time.Parse(time.RFC3339, "2025-07-09T21:56:18+04:00")
	if err != nil {
		t.Fatalf("Failed to parse expected time: %v", err)
	}

	want := models.NewRecord(5218837036, 58.0, expectedTime, 0)

	got, _ := ParseRecord(value)

	// Сравниваем структуры поэлементно
	if got.GetId() != want.GetId() {
		t.Errorf("UserID mismatch: got %v, want %v", got.GetId(), want.GetId())
	}
	if got.GetWeight() != want.GetWeight() {
		t.Errorf("Weight mismatch: got %v, want %v", got.GetWeight(), want.GetWeight())
	}
	if !got.GetTime().Equal(want.GetTime()) {
		t.Errorf("DateTime mismatch: got %v, want %v", got.GetTime(), want.GetTime())
	}
}
