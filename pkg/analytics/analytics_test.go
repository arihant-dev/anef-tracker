package analytics_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/analytics"
	"testing"
)

func TestMedianDuration(t *testing.T) {
	durations := []float64{5.0, 10.0, 15.0, 20.0, 25.0}
	med := analytics.CalculateMedian(durations)
	if med != 15.0 {
		t.Errorf("expected median 15.0, got %.1f", med)
	}

	evenDurations := []float64{10.0, 20.0, 30.0, 40.0}
	evenMed := analytics.CalculateMedian(evenDurations)
	if evenMed != 25.0 {
		t.Errorf("expected median 25.0 for even elements, got %.1f", evenMed)
	}
}

func TestPercentile(t *testing.T) {
	durations := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	p75 := analytics.CalculatePercentile(durations, 75.0)
	if p75 < 7.0 || p75 > 8.0 {
		t.Errorf("expected 75th percentile ~8.0, got %.1f", p75)
	}
}

func TestConfidenceScore(t *testing.T) {
	score := analytics.ComputeConfidenceScore(150, 2.5)
	if score < 0.80 {
		t.Errorf("expected high confidence score for 150 samples, got %.2f", score)
	}
}
