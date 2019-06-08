package golicko_test

import (
	"fmt"

	"github.com/mashiike/golicko"
)

func ExampleRating() {
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
			},
			Score: golicko.ScoreWin,
		},
		golicko.Result{
			Opponent: golicko.Rating{
				Value:      1550.0,
				Deviation:  100.0,
				Volatility: 0.06,
			},
			Score: golicko.ScoreLose,
		},
		golicko.Result{
			Opponent: golicko.Rating{
				Value:      1700.0,
				Deviation:  300.0,
				Volatility: 0.06,
			},
			Score: golicko.ScoreLose,
		},
	}
	fmt.Println(player.Update(results, 0.5))
	//Output:
	//rating:1464.05 (1161.02-1767.08), volatility:0.059993
}
