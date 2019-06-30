package ratingutil

import (
	"fmt"
	"time"

	"github.com/mashiike/rating"
	"github.com/pkg/errors"
)

//Element is an interface of Team/Player
type Element interface {
	Name() string
	Rating() rating.Rating
	ApplyMatch(rating.Rating, float64) error
	Prepare(time.Time, *Config) error
}

//Player is rating resouse
type Player struct {
	name      string
	estimated *rating.Estimated
	fixedAt   time.Time
}

//Name is player name
func (p *Player) Name() string {
	return p.name
}

//ApplyMatch reflects match results between players.
func (p *Player) ApplyMatch(opponent rating.Rating, score float64) error {
	return p.estimated.ApplyMatch(opponent, score)
}

//Prepare does the work before operating the player according to the configs
func (p *Player) Prepare(outcomeAt time.Time, config *Config) error {
	for outcomeAt.Sub(p.fixedAt) > config.RatingPeriod {
		if err := p.estimated.Fix(config.Tau); err != nil {
			return err
		}
		p.fixedAt = p.fixedAt.Add(config.RatingPeriod)
	}
	return nil
}

//Rating returns the estimated strength of the current player
func (p *Player) Rating() rating.Rating {
	return p.estimated.Rating()
}

//String is implements fmt.Stringer
func (p *Player) String() string {
	return fmt.Sprintf("%s:%s", p.Name(), p.Rating().Format(rating.PlusMinusFormat))
}

//Team is multiple player
// http://rhetoricstudios.com/downloads/AbstractingGlicko2ForTeamGames.pdf
type Team struct {
	name    string
	members []*Player
}

//Name is team name
func (t *Team) Name() string {
	return t.name
}

//ApplyMatch reflects match results between teams.
func (t *Team) ApplyMatch(opponent rating.Rating, score float64) error {
	for _, member := range t.members {
		if err := member.ApplyMatch(opponent, score); err != nil {
			return errors.Wrapf(err, "apply match %v", member)
		}
	}
	return nil
}

//Prepare does the work before operating the team according to the configs
func (t *Team) Prepare(outcomeAt time.Time, config *Config) error {
	for _, member := range t.members {
		if err := member.Prepare(outcomeAt, config); err != nil {
			return errors.Wrapf(err, "prepare %v", member)
		}
	}
	return nil
}

//Rating return estimated team rating
func (t *Team) Rating() rating.Rating {
	ratings := make([]rating.Rating, 0, len(t.members))
	for _, member := range t.members {
		ratings = append(ratings, member.Rating())
	}
	return rating.Average(ratings)
}

//String is implements fmt.Stringer
func (t *Team) String() string {
	str := fmt.Sprintf("%s:{", t.Name())
	for _, member := range t.members {
		str += " " + member.String()
	}
	str += " }"
	return str
}
