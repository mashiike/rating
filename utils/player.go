package utils

import (
	"strconv"
	"time"

	"github.com/mashiike/rating/rating"
	"github.com/pkg/errors"
)

//GameCount is (win/lose/draw) count
type GameCount struct {
	Win  uint64 `json:"win"`
	Lose uint64 `json:"lose"`
	Draw uint64 `json:"draw"`
}

func (g GameCount) byteLength() int {
	return 13
}

func (g GameCount) appendByte(b []byte) []byte {
	b = append(b, '(')
	b = append(b, strconv.FormatUint(g.Win, 10)...)
	b = append(b, '/')
	b = append(b, strconv.FormatUint(g.Lose, 10)...)
	b = append(b, '/')
	b = append(b, strconv.FormatUint(g.Draw, 10)...)
	b = append(b, ')')
	return b
}

//String is implemets fmt.Stringer
func (g GameCount) String() string {
	b := make([]byte, 0, g.byteLength())
	b = g.appendByte(b)
	return string(b)
}

//Player is rating holder, implements IPlayer
type Player struct {
	Estimated *rating.Estimated `json:"estimated"`
	FixedAt   time.Time         `json:"fixed_at"`
	GameCount
}

//NewPlayer is Player constractor
func NewPlayer(createdAt time.Time, config *Config) *Player {

	return &Player{
		Estimated: rating.NewEstimated(
			rating.Default(config.InitialVolatility()),
			config.Tau,
		),
		FixedAt: createdAt.Truncate(config.RatingPeriod),
	}
}

//Rating is current this player's estimated rating
func (p *Player) Rating() rating.Rating {
	return p.Estimated.Rating()
}

//Prepare must do before Update
func (p *Player) Prepare(now time.Time, config *Config) error {
	//update system tou: maybe not change often
	p.Estimated.Tau = config.Tau

	//Reflects the previous non-match period.
	for now.Sub(p.FixedAt) > config.RatingPeriod {
		if err := p.Estimated.Fix(); err != nil {
			return err
		}
		p.FixedAt = p.FixedAt.Add(config.RatingPeriod)
	}
	return nil
}

//Update do Player's rating update.
func (p *Player) Update(result *MatchResult, _ *Config) error {
	if p.FixedAt.After(result.OutcomeAt) {
		return errors.New("a match from the pasted")
	}
	if err := p.Estimated.ApplyMatch(result.Opponent, result.Score); err != nil {
		return errors.Wrap(err, "player update")
	}
	switch result.Score {
	case rating.ScoreDraw:
		p.GameCount.Draw++
	case rating.ScoreLose:
		p.GameCount.Lose++
	case rating.ScoreWin:
		p.GameCount.Win++
	}
	return nil
}

//String is implements fmt.Stringer
func (p *Player) String() string {
	b := make([]byte, 0, len(rating.PlusMinusFormat)+p.GameCount.byteLength())
	b = p.Rating().AppendFormat(b, rating.PlusMinusFormat)
	b = p.GameCount.appendByte(b)
	return string(b)
}
