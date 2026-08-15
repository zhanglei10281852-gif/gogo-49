package detect

import (
	"testing"
	"time"

	"LogPilot/internal/config"
	"LogPilot/internal/model"
)

func TestRepeatedDetectionDoesNotInflateIncidentCount(t *testing.T) {
	base := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	rule := config.Rule{ID: "errors", Title: "Repeated errors", Enabled: true, Severity: "ERROR", All: []config.Predicate{{Field: "level", Op: "gte", Value: "ERROR"}}, GroupBy: []string{"source"}, Threshold: 2, Window: "5m", Cooldown: "10m", ResolveAfter: "30m"}
	events := []model.Event{
		model.NewEvent(base, model.LevelError, "api", "failed one", nil),
		model.NewEvent(base.Add(time.Minute), model.LevelError, "api", "failed two", nil),
	}
	engine := New([]config.Rule{rule})
	engine.Now = func() time.Time { return base.Add(2 * time.Minute) }
	first, err := engine.Detect(events, nil)
	if err != nil {
		t.Fatal(err)
	}
	second, err := engine.Detect(events, first.Incidents)
	if err != nil {
		t.Fatal(err)
	}
	if len(second.Incidents) != 1 || second.Incidents[0].Count != 2 {
		t.Fatalf("reprocessing unchanged events inflated incident: %+v", second.Incidents)
	}
}
