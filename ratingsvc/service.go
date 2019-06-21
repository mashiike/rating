package ratingsvc

import (
	"time"

	"github.com/mashiike/rating"
	"github.com/pkg/errors"
)

//Service is usecase
type Service struct {
	*Config
}

func New(config *Config) *Service {
	return &Service{
		Config: config,
	}
}

//NewPlayer is constractor of *Player
func (s *Service) NewPlayer(name string, fixed rating.Rating, fixedAt time.Time) *Player {
	return &Player{
		name:      name,
		estimated: rating.NewEstimated(fixed, s.Config.Tau),
		fixedAt:   fixedAt.Truncate(s.Config.RatingPeriod),
	}
}

//NewDefaultPlayer is factory of default player
func (s *Service) NewDefaultPlayer(name string) *Player {
	return s.NewPlayer(
		name,
		rating.Default(s.Config.InitialVolatility()),
		s.Config.Now(),
	)
}

type Players []*Player

//NewTeam is constractor of *Team
func (s *Service) NewTeam(name string, members []*Player) *Team {
	return &Team{
		name:    name,
		members: members,
	}
}

type Outcome map[Element]float64

//ApplyMatchWithTime is a function for apply Match for Team/Player outcome.
func (s *Service) ApplyMatchWithTime(outcome Outcome, outcomeAt time.Time) error {

	opponents := make(map[Element]rating.Rating, len(outcome))
	for target, _ := range outcome {
		if err := target.Prepare(s.Config.RatingPeriod, outcomeAt); err != nil {
			return errors.Wrapf(err, "failed prepare %v", target.Name())
		}
		opponents[target] = target.Rating()
	}

	for target, result1 := range outcome {
		for opponent, result2 := range outcome {
			if target == opponent {
				continue
			}

			score := rating.ScoreLose
			if result1 > result2 {
				score = rating.ScoreWin
			}
			if result1 == result2 {
				score = rating.ScoreDraw
			}
			if err := target.ApplyMatch(opponents[opponent], score); err != nil {
				return errors.Wrapf(err, "failed apply %v vs %v", target.Name(), opponent.Name())
			}
		}
	}
	return nil
}

//ApplyMatch is a function for apply Match for Team/Player outcome.
func (s *Service) ApplyMatch(outcome Outcome) error {
	return s.ApplyMatchWithTime(outcome, s.Config.Now())
}
