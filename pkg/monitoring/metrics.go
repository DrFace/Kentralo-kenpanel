package monitoring

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// MetricSample represents a single point-in-time telemetry datum.
type MetricSample struct {
	Timestamp   time.Time `json:"timestamp"`
	CPUUsagePct float64   `json:"cpu_usage_pct"`
	RAMUsagePct float64   `json:"ram_usage_pct"`
	DiskUsagePct float64  `json:"disk_usage_pct"`
	NetInBps    int64     `json:"net_in_bps"`
	NetOutBps   int64     `json:"net_out_bps"`
	LoadAvg1m   float64   `json:"load_avg_1m"`
}

// AlertRule defines a threshold condition.
type AlertRule struct {
	ID        string  `json:"id"`
	Metric    string  `json:"metric"` // cpu, ram, disk, uptime
	Condition string  `json:"condition"` // gt, lt
	Threshold float64 `json:"threshold"`
	Duration  time.Duration `json:"duration"`
	Severity  string  `json:"severity"` // critical, warning, info
	Channel   string  `json:"channel"`  // webhook, email, slack
}

// UptimeCheck defines an endpoint probe.
type UptimeCheck struct {
	ID             string        `json:"id"`
	TargetURL      string        `json:"target_url"`
	Interval       time.Duration `json:"interval"`
	ExpectedStatus int           `json:"expected_status"`
	LastStatus     int           `json:"last_status"`
	LastLatency    time.Duration `json:"last_latency"`
	IsUp           bool          `json:"is_up"`
	LastChecked    time.Time     `json:"last_checked"`
}

// MetricsEngine collects telemetry, runs uptime probes and evaluates alerts.
type MetricsEngine struct {
	mu          sync.RWMutex
	samples     []MetricSample
	alerts      map[string]*AlertRule
	uptimeChecks map[string]*UptimeCheck
}

// NewMetricsEngine initializes monitoring.
func NewMetricsEngine() *MetricsEngine {
	return &MetricsEngine{
		samples:      make([]MetricSample, 0),
		alerts:       make(map[string]*AlertRule),
		uptimeChecks: make(map[string]*UptimeCheck),
	}
}

// RecordSample appends telemetry and enforces a retention buffer.
func (e *MetricsEngine) RecordSample(s MetricSample) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.samples = append(e.samples, s)
	if len(e.samples) > 1000 {
		e.samples = e.samples[len(e.samples)-1000:] // retain latest 1000 samples
	}
}

// CheckUptime executes an HTTP probe.
func (e *MetricsEngine) CheckUptime(ctx context.Context, check *UptimeCheck) {
	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", check.TargetURL, nil)
	if err != nil {
		check.IsUp = false
		check.LastStatus = 0
		check.LastLatency = time.Since(start)
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	check.LastLatency = time.Since(start)
	check.LastChecked = time.Now().UTC()

	if err != nil {
		check.IsUp = false
		check.LastStatus = 0
		return
	}
	defer resp.Body.Close()

	check.LastStatus = resp.StatusCode
	check.IsUp = (resp.StatusCode == check.ExpectedStatus || (check.ExpectedStatus == 0 && resp.StatusCode >= 200 && resp.StatusCode < 400))
}
