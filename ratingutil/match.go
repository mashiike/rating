package ratingutil

import (
	"fmt"
	"sort"
	"time"

	"github.com/mashiike/rating"
	"github.com/pkg/errors"
)

//ApplyStrategy is an alias for a function.
//This function shows how to reflect in a multiplayer game when reflecting the result of Match in Rating.
type ApplyStrategy func(map[Element]rating.Rating, map[Element]float64) error

//Match is a model that represents multiple Team / Player battles
type Match struct {
	scores        map[Element]float64
	applyStrategy ApplyStrategy
}

//Add function adds a Score.
func (m *Match) Add(element Element, score float64) error {
	if _, ok := m.scores[element]; !ok {
		return errors.New("this element not join match")
	}
	m.scores[element] += score
	return nil
}

//Reset returns to the zero score
func (m *Match) Reset() {
	for element := range m.scores {
		m.scores[element] = 0.0
	}
}

//Scores return copy internal scores
func (m *Match) Scores() map[Element]float64 {
	ret := make(map[Element]float64, len(m.scores))
	for elem, score := range m.scores {
		ret[elem] = score
	}
	return ret
}

//Ratings return match joined Team/Players current Rating
func (m *Match) Ratings() map[Element]rating.Rating {
	ratings := make(map[Element]rating.Rating, len(m.scores))
	for target := range m.scores {
		ratings[target] = target.Rating()
	}
	return ratings
}

//Apply function determines the current score and reflects it on Team / Player's Rating.
func (m *Match) Apply(scoresAt time.Time, config *Config) error {
	scores := make(map[Element]float64, len(m.scores))
	for target, score := range m.scores {
		if err := target.Prepare(scoresAt, config); err != nil {
			return errors.Wrapf(err, "failed prepare %v", target.Name())
		}
		scores[target] = score
	}
	ratings := m.Ratings()
	m.Reset()
	return m.applyStrategy(ratings, scores)
}

//WinProbs returns the probability that each Team / Player will be first
func (m *Match) WinProbs() map[Element]float64 {
	probs := make(map[Element]float64, len(m.scores))
	ratings := m.Ratings()
	for target := range m.scores {
		probs[target] = 1.0
		for opponent, r := range ratings {
			if target == opponent {
				continue
			}
			probs[target] *= ratings[target].WinProb(r)
		}
	}
	return probs
}

func (m *Match) String() string {
	probs := m.WinProbs()
	sortedKey := make([]Element, 0, len(probs))
	for elem := range probs {
		sortedKey = append(sortedKey, elem)
	}
	sort.Slice(sortedKey, func(i, j int) bool { return sortedKey[i].Name() < sortedKey[j].Name() })
	str := "["
	for _, elem := range sortedKey {
		str += fmt.Sprintf(" %s(%0.2f) ", elem, probs[elem])
	}
	return str + "]"
}

//AsRoundrobin considers Multiplayer Matches to be a round-trip tournament ApplyStrategy
func AsRoundrobin(ratings map[Element]rating.Rating, scores map[Element]float64) error {
	for target, score1 := range scores {
		for opponent, score2 := range scores {
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
			if err := target.ApplyMatch(ratings[opponent], score); err != nil {
				return errors.Wrapf(err, "failed apply %v vs %v", target.Name(), opponent.Name())
			}
		}
	}
	return nil
}
