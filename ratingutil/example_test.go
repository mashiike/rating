package ratingutil_test

import (
	"fmt"

	"github.com/mashiike/rating"
	"github.com/mashiike/rating/ratingutil"
)

func ExampleService() {

	//for example, Tag battle
	svc := ratingutil.New(ratingutil.NewConfig())
	team1 := svc.NewTeam(
		"bovidae",
		ratingutil.Players{
			svc.NewPlayer(
				"sheep",
				rating.New(1700.0, 50.0, svc.Config.InitialVolatility()),
				svc.Config.Now(),
			),
			svc.NewDefaultPlayer("goat"),
		},
	)
	team2 := svc.NewTeam(
		"equidae",
		ratingutil.Players{
			svc.NewPlayer(
				"donkey",
				rating.New(1400.0, 50.0, svc.Config.InitialVolatility()),
				svc.Config.Now(),
			),
			svc.NewDefaultPlayer("zebra"),
		},
	)
	match, _ := svc.NewMatch(team1, team2)

	fmt.Println(match)
	match.Add(team1, 1.0)
	match.Add(team2, 0.0)
	err := svc.Apply(match)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("=== after match ===")
	fmt.Println(match)

	//Output:
	//[ bovidae:{ sheep:1700.0p-138.6 goat:1500.0p-700.0 }(0.66)  equidae:{ donkey:1400.0p-138.6 zebra:1500.0p-700.0 }(0.34) ]
	//=== after match ===
	//[ bovidae:{ sheep:1705.2p-137.2 goat:1654.5p-530.9 }(0.81)  equidae:{ donkey:1393.7p-137.0 zebra:1364.1p-536.3 }(0.19) ]

}
