package ratingutil

import (
	"time"

	"github.com/mashiike/rating"
)

//RatingPeriod constants
//can multiple float64
//  PeriodDay * 3.0  => 3 days
const (
	PeriodDay   time.Duration = 24 * time.Hour
	PeriodWeek                = 7 * PeriodDay
	PeriodMonth               = 30 * PeriodDay
	PeriodYear                = 365 * PeriodDay
)

//Clock is a clock used in this package. The default is to use time.Now()
type Clock interface {
	Now() time.Time
}

type defaultClock struct{}

func (c defaultClock) Now() time.Time {
	return time.Now()
}

//Config is service configuration
type Config struct {
	//Service Clock
	Clock
	Tau float64

	//RatingPeriod is the fixed interval of Rating.
	//All matches played between this interval are considered to occur simultaneously and are calculated.
	//In RatingPeriod, the period in which the players play about 15 times is good.
	RatingPeriod time.Duration

	//It will return to the initial deviation if you have not played for about this period.
	//This period is a guideline, and the time to return to the actual initial deviation is determined by the player's Volatility here.
	//And initial Volatility is calculated based on this period.
	PeriodToResetDeviation time.Duration

	//Fighting prosperity strategy uses round robin by default
	DefaultApplyStrategy ApplyStrategy
}

//NewConfig is default configuration
func NewConfig() *Config {
	return &Config{
		Clock:                  defaultClock{},
		RatingPeriod:           PeriodWeek,
		PeriodToResetDeviation: PeriodYear,
		Tau:                    0.5,
		DefaultApplyStrategy:   AsRoundrobin,
	}
}

//InitialVolatility calculates the initial rating fluctuation according to the setting
func (c *Config) InitialVolatility() float64 {
	count := c.PeriodToResetDeviation.Seconds() / c.RatingPeriod.Seconds()
	return rating.NewVolatility(50.0, count)
}

//WithClock is set clock to config
func (c *Config) WithClock(clock Clock) *Config {
	c.Clock = clock
	return c
}

//WithRatingPeriod is set RatingPeriod to config
func (c *Config) WithRatingPeriod(period time.Duration) *Config {
	c.RatingPeriod = period
	return c
}

//WithTau is set Tau to config
func (c *Config) WithTau(tau float64) *Config {
	c.Tau = tau
	return c
}

//WithApplyStrategy is set DefaultApplyStrategy to config
func (c *Config) WithApplyStrategy(strategy ApplyStrategy) *Config {
	if strategy != nil {
		c.DefaultApplyStrategy = strategy
	}
	return c
}
