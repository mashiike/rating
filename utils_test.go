package rating_test

import (
	"testing"

	"github.com/mashiike/rating"
)

func TestComputeInitialVolatility(t *testing.T) {
	v := rating.ComputeInitialVolatility(50.0, 350.0, 100.0)
	expected := 0.199409
	if v != expected {
		t.Errorf("rating.ComputeInitialVolatility got %v, expected %v", v, expected)
	}
}
