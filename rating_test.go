package golicko_test

import (
	"testing"

	"github.com/mashiike/golicko"
)

func TestPaperExample(t *testing.T) {
	player := golicko.Rating{
		Value:      1500.0,
		Deviation:  200.0,
		Volatility: 0.06,
	}
	results := []golicko.Result{
		golicko.Result{
			Opponent: golicko.Rating{
				Value:      1400.0,
				Deviation:  30.0,
				Volatility: 0.06,
			}.ToGlicko2(),
			Score: golicko.ScoreWin,
		},
		golicko.Result{
			Opponent: golicko.Rating{
				Value:      1550.0,
				Deviation:  100.0,
				Volatility: 0.06,
			}.ToGlicko2(),
			Score: golicko.ScoreLose,
		},
		golicko.Result{
			Opponent: golicko.Rating{
				Value:      1700.0,
				Deviation:  300.0,
				Volatility: 0.06,
			}.ToGlicko2(),
			Score: golicko.ScoreLose,
		},
	}
	gplayer := player.ToGlicko2()
	quantity := gplayer.ComputeQuantity(results)
	got := gplayer.ApplyQuantity(quantity, 0.5).ToRating()
	expected := golicko.Rating{
		Value:      1464.06,
		Deviation:  151.52,
		Volatility: 0.05999,
	}
	if !got.Equals(expected) {
		t.Errorf("unexpected rating: %+v", got)
		t.Logf("quantity: %+v", quantity)
	}
}
