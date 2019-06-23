package ratingutil

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
		estimated: rating.NewEstimated(fixed),
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
func (s *Service) NewTeam(name string, members Players) *Team {
	return &Team{
		name:    name,
		members: members,
	}
}

//NewMatch creates a new Match model from a given Team / Player
func (s *Service) NewMatch(elements ...Element) (*Match, error) {
	if len(elements) < 2 {
		return nil, errors.New("two or more elements are required for the match")
	}
	scores := make(map[Element]float64, len(elements))
	for _, element := range elements {
		scores[element] = 0.0
	}
	return &Match{
		scores: scores,
	}, nil
}

//ApplyWithTime is a function for apply Match for Team/Player outcome.
func (s *Service) ApplyWithTime(match *Match, outcomeAt time.Time) error {
	return match.Apply(outcomeAt, s.Config)
}

//Apply is a function for apply Match for Team/Player outcome.
func (s *Service) Apply(match *Match) error {
	return s.ApplyWithTime(match, s.Config.Now())
}
