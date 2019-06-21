package rating

import "github.com/pkg/errors"

//TeamMember is an interface of elements included by Team
type TeamMember interface {
	Rating() Rating
	ApplyMatch(Rating, float64) error
	Fix() error
}

//Team is multiple player
// http://rhetoricstudios.com/downloads/AbstractingGlicko2ForTeamGames.pdf
type Team struct {
	Members []TeamMember
}

//ApplyMatch reflects match results between teams.
func (t *Team) ApplyMatch(opponent *Team, score float64) error {
	opponentRating := opponent.Rating()
	for _, member := range t.Members {
		if err := member.ApplyMatch(opponentRating, score); err != nil {
			return errors.Wrapf(err, "apply match %v", member)
		}
	}
	return nil
}

//Rating return estimated team rating
func (t *Team) Rating() Rating {
	ratings := make([]Rating, 0, len(t.Members))
	for _, member := range t.Members {
		ratings = append(ratings, member.Rating())
	}
	return Average(ratings)
}

//Fix close rating period.
func (t *Team) Fix() error {
	for _, member := range t.Members {
		if err := member.Fix(); err != nil {
			return errors.Wrapf(err, "fix %v", member)
		}
	}
	return nil
}
