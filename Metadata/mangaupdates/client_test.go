package mangaupdates

import (
	"testing"
)	

func TestSearchSeriesID(t *testing.T) {
	title := "Naruto"

	id, error := SearchSeriesID(title)
	if error != nil {
		t.Fatalf("Failed to search series ID for %s: %v", title, error)
	}

	if id <= 0 {
		t.Fatalf("Expected a valid series ID, got %d", id)
	}

	t.Logf("Found series ID %d for title %s", id, title)
}

func TestGetSeriesPublications(t *testing.T) {
	var id = 63868402874

	publications, publishers, err := GetSeriesPublications(id)
	if err != nil {
		t.Fatalf("Failed to get series publications for ID %d: %v", id, err)
	}

	if publications == "" || publishers == "" {
		t.Fatalf("Expected non-empty publications and publishers, got publications: %s, publishers: %s", publications, publishers)
	}

	t.Logf("Publications: %s, Publishers: %s", publications, publishers)

}