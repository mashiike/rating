package utils

import (
	"reflect"
	"time"

	"github.com/mashiike/golicko/rating"
	"github.com/pkg/errors"
)

//MatchResult represents the Match result from the viewpoint of a player
type MatchResult struct {
	Opponent  rating.Rating
	Score     float64
	OutcomeAt time.Time
}

//IPlayer is an opponent's interface, which can be Player / Team, etc.
type IPlayer interface {
	Rating() rating.Rating
	Prepare(time.Time, *Config) error
	Update(*MatchResult, *Config) error
}

//Match represents a match from the perspective of the system.
//Specify nil if there is no winner
type Match struct {
	Left      IPlayer
	Right     IPlayer
	Winner    IPlayer
	OutcomeAt time.Time
}

func (m *Match) newResult(opponent IPlayer) *MatchResult {
	score := rating.ScoreWin
	if m.Winner == nil {
		score = rating.ScoreDraw
	}
	if m.Winner == opponent {
		score = rating.ScoreLose
	}
	return &MatchResult{
		Opponent:  opponent.Rating(),
		Score:     score,
		OutcomeAt: m.OutcomeAt,
	}
}

//LeftResult is create Left Player's /Team's match result
func (m *Match) LeftResult() *MatchResult {
	return m.newResult(m.Right)
}

//RightResult is create Left Player's /Team's match result
func (m *Match) RightResult() *MatchResult {
	return m.newResult(m.Left)
}

//Apply will reflect the game results in the rating of the player
func (m *Match) Apply(config *Config) error {
	if m.Left == nil {
		return errors.New("left player is nil")
	}
	if m.Right == nil {
		return errors.New("right player is nil")
	}
	if m.Left == m.Right {
		return errors.New("both player is same")
	}
	if !(isNil(m.Winner) || m.Left == m.Winner || m.Right == m.Winner) {
		return errors.New("not match player win")
	}
	if config == nil {
		return errors.New("config is nil")
	}

	if err := m.Left.Prepare(m.OutcomeAt, config); err != nil {
		return errors.Wrap(err, "left player prepare")
	}
	if err := m.Right.Prepare(m.OutcomeAt, config); err != nil {
		return errors.Wrap(err, "right player prepare")
	}

	//before apply strength snapshot
	leftResult := m.LeftResult()
	rightResult := m.RightResult()

	if err := m.Left.Update(leftResult, config); err != nil {
		return errors.Wrap(err, "left player update")
	}
	if err := m.Right.Update(rightResult, config); err != nil {
		return errors.Wrap(err, "right player update")
	}

	return nil
}

func isNil(x interface{}) bool {
	if x == nil || reflect.ValueOf(x).IsNil() {
		return true
	}
	return false
}
