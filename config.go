package golicko

import (
	"time"

	"github.com/mashiike/golicko/rating"
)

//RatingPeriod constans
//can multiple float64
//  PeriodDay * 3.0  => 3 days
const (
	PeriodDay   time.Duration = 24 * time.Hour
	PeriodWeek                = 7 * PeriodDay
	PeriodMonth               = 30 * PeriodDay
	PeriodYear                = 365 * PeriodDay
)

//Config is rating system propaties
type Config struct {
	//RatingPeriod is the fixed interval of Rating.
	//All matches played between this interval are considered to occur simultaneously and are calculated.
	//In RatingPeriod, the period in which the players play about 15 times is good.
	RatingPeriod time.Duration

	//It will return to the initial deviation if you have not played for about this period.
	//This period is a guideline, and the time to return to the actual initial deviation is determined by the player's Volatility here.
	//And initial Volatility is calculated based on this period.
	PeriodToResetDeviation time.Duration

	Tau float64
}

//NewConfig is default configuration
func NewConfig() *Config {
	return &Config{
		RatingPeriod:           PeriodWeek,
		PeriodToResetDeviation: PeriodYear,
		Tau:                    0.5,
	}
}

//InitialVolatility detamine initial volatility from rating period and period to reset deviation
func (c *Config) InitialVolatility() float64 {
	count := c.PeriodToResetDeviation.Seconds() / c.RatingPeriod.Seconds()
	return rating.ComputeInitialVolatility(50.0, rating.InitialDeviation, count)
}
