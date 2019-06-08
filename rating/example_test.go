package rating_test

import (
	"fmt"

	"github.com/mashiike/golicko/rating"
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
