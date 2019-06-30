package rating_test

import (
	"fmt"

	"github.com/mashiike/rating"
)

func ExampleRating() {
	player := rating.New(1500.0, 200.0, 0.06)
	opponents := []rating.Rating{
		rating.New(1400.0, 30.0, 0.06),
		rating.New(1550.0, 100.0, 0.06),
		rating.New(1700.0, 300.0, 0.06),
	}
	scores := []float64{
		rating.ScoreWin,
		rating.ScoreLose,
		rating.ScoreLose,
	}
	updated, _ := player.Update(opponents, scores, 0.5)
	fmt.Println(updated)
	fmt.Printf("strength  : %f\n", updated.Strength())
	fmt.Printf("deviation : %f\n", updated.Deviation())
	fmt.Printf("volatility: %f\n", updated.Volatility())
	//Output:
	//1464.0 (1161.0-1767.1 v=0.059996)
	//strength  : 1464.050000
	//deviation : 151.510000
	//volatility: 0.059996
}

func ExampleEstimated_Rating() {
	player := rating.New(1500.0, 200.0, 0.06)
	opponents := []rating.Rating{
		rating.New(1400.0, 30.0, 0.06),
		rating.New(1550.0, 100.0, 0.06),
		rating.New(1700.0, 300.0, 0.06),
	}
	scores := []float64{
		rating.ScoreWin,
		rating.ScoreLose,
		rating.ScoreLose,
	}
	prev := player
	e := rating.NewEstimated(prev)
	for i := 0; i < len(opponents); i++ {
		e.ApplyMatch(opponents[i], scores[i])
		updated := e.Rating()
		fmt.Println(updated)
		fmt.Printf("strength  diff : %f\n", updated.Strength()-prev.Strength())
		fmt.Printf("deviation diff : %f\n", updated.Deviation()-prev.Deviation())
		fmt.Println("---")
		prev = updated
	}
	e.Fix(0.5)
	fmt.Println(e.Fixed)
	//Output:
	//1563.6 (1212.8-1914.4 v=0.06)
	//strength  diff : 63.560000
	//deviation diff : -24.600000
	//---
	//1492.4 (1175.7-1809.1 v=0.06)
	//strength  diff : -71.170000
	//deviation diff : -17.070000
	//---
	//1464.0 (1161.0-1767.1 v=0.06)
	//strength  diff : -28.340000
	//deviation diff : -6.820000
	//---
	//1464.0 (1161.0-1767.1 v=0.059996)
}
