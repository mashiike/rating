package ratingutil

import (
	"time"

	"github.com/mashiike/rating"
	"github.com/pkg/errors"
)

type Match struct {
	outcome map[Element]float64
}

func (m *Match) AddScore(element Element, score float64) error {
	if _, ok := m.outcome[element]; !ok {
		return errors.New("this element not join match")
	}
	m.outcome[element] += score
	return nil
}

func (m *Match) ResetScore() {
	for element, _ := range m.outcome {
		m.outcome[element] = 0.0
	}
}

func (m *Match) Ratings() map[string]rating.Rating {
	ratings := make(map[string]rating.Rating, len(m.outcome))
	for target, _ := range m.outcome {
		ratings[target.Name()] = target.Rating()
	}
	return ratings
}

func (m *Match) Apply(outcomeAt time.Time, config *Config) error {
	for target, _ := range m.outcome {
		if err := target.Prepare(outcomeAt, config); err != nil {
			return errors.Wrapf(err, "failed prepare %v", target.Name())
		}
	}
	opponents := m.Ratings()
	for target, score1 := range m.outcome {
		for opponent, score2 := range m.outcome {
			if target == opponent {
				continue
			}

			score := rating.ScoreLose
			if score1 > score2 {
				score = rating.ScoreWin
			}
			if score1 == score2 {
				score = rating.ScoreDraw
			}
			if err := target.ApplyMatch(opponents[opponent.Name()], score); err != nil {
				return errors.Wrapf(err, "failed apply %v vs %v", target.Name(), opponent.Name())
			}
		}
	}
	m.ResetScore()
	return nil
}
