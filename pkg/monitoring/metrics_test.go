package monitoring

import (
	"testing"
	"time"
)

func TestMetricsBuffer(t *testing.T) {
	engine := NewMetricsEngine()

	sample := MetricSample{
		Timestamp:   time.Now().UTC(),
		CPUUsagePct: 15.5,
		RAMUsagePct: 42.0,
		LoadAvg1m:   0.85,
	}

	engine.RecordSample(sample)
	if len(engine.samples) != 1 {
		t.Errorf("expected 1 sample, got %d", len(engine.samples))
	}
}
