# goliko
[![GoDoc](https://godoc.org/github.com/mashiike/golicko?status.svg)](https://godoc.org/github.com/mashiike/golicko)
[![Go Report Card](https://goreportcard.com/badge/github.com/mashiike/golicko)](https://goreportcard.com/report/github.com/mashiike/golicko)
[![CircleCI](https://circleci.com/gh/mashiike/golicko/tree/master.svg?style=svg)](https://circleci.com/gh/mashiike/golicko/tree/master)


This is the Go implementation of Gliko2 Rating

## Usage

### Import the package

```go
import "github.com/mashiike/golicko"
```

### Use Rating and Result struct

```go
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
		Score: golicko.ScoreDraw,
	},
}
player = player.Update(results, golicko.DefaultSetting))
```
