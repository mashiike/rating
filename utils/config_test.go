package utils_test

import (
	"testing"

	"github.com/mashiike/rating"
)

func TestConfig_InititalVolatility(t *testing.T) {
	c := utils.NewConfig()
	got := c.InitialVolatility()
	expected := 0.276152
	if got != expected {
		t.Errorf("Config.InitialVolatility got %v, expected %v", got, expected)
	}
}
