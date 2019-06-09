# rating
[![GoDoc](https://godoc.org/github.com/mashiike/rating?status.svg)](https://godoc.org/github.com/mashiike/rating)
[![Go Report Card](https://goreportcard.com/badge/github.com/mashiike/rating)](https://goreportcard.com/report/github.com/mashiike/rating)
[![CircleCI](https://circleci.com/gh/mashiike/rating/tree/master.svg?style=svg)](https://circleci.com/gh/mashiike/rating/tree/master)


This is the Go implementation of Gliko2 Rating
In a simple use case we use as follows

### Import the package

```go
import "github.com/mashiike/rating"
```

### Usage: Batch by Rating Period.
At the end of each rating period, it reflects the results of the game played in that period.
```go
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
```

### Usage: Sequentially for each game
Use Esteminated struct as follows.
Then save the information in a database or file system etc.  

```go
player := rating.New(1500.0, 200.0, 0.06)
e := rating.NewEstimated(r, 0.5)
opponent := rating.New(1400.0, 30.0, 0.06)
err := e.ApplyMatch(opponent, rating.ScoreWin)
//when the rating period is over
err = e.Fix()
updated := e.Fixed
```

