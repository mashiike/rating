package golicko_test

import (
	"testing"

	"github.com/mashiike/golicko"
)

func TestConfig_InititalVolatility(t *testing.T) {
	c := golicko.NewConfig()
	got := c.InitialVolatility()
	expected := 0.276152
	if got != expected {
		t.Errorf("Config.InitialVolatility got %v, expected %v", got, expected)
	}
}
