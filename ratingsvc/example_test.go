package ratingsvc_test

import (
	"fmt"

	"github.com/mashiike/rating"
	"github.com/mashiike/rating/ratingsvc"
)

func ExampleService() {

	//for example, Tag battle
	svc := ratingsvc.New(ratingsvc.NewConfig())
	team1 := svc.NewTeam(
		"bovidae",
		ratingsvc.Players{
			svc.NewPlayer("sheep", rating.New(1700.0, 50.0, svc.Config.InitialVolatility()), svc.Config.Now()),
			svc.NewDefaultPlayer("goat"),
		},
	)
	team2 := svc.NewTeam(
		"equidae",
		ratingsvc.Players{
			svc.NewPlayer("donkey", rating.New(1400.0, 50.0, svc.Config.InitialVolatility()), svc.Config.Now()),
			svc.NewDefaultPlayer("zebra"),
		},
	)
	fmt.Println(team1)
	fmt.Println(team2)
	fmt.Printf("== %s win %% = %f ==\n", team1.Name(), team1.Rating().WinProb(team2.Rating()))
	err := svc.ApplyMatch(ratingsvc.Outcome{team1: 1.0, team2: 0.0})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(team1)
	fmt.Println(team2)
	fmt.Printf("== %s win %% = %f ==\n", team1.Name(), team1.Rating().WinProb(team2.Rating()))

	//Output:
	//bovidae:{ sheep:1700.0p-138.6 goat:1500.0p-700.0 }
	//equidae:{ donkey:1400.0p-138.6 zebra:1500.0p-700.0 }
	//== bovidae win % = 0.662400 ==
	//bovidae:{ sheep:1705.2p-137.2 goat:1654.5p-530.9 }
	//equidae:{ donkey:1393.7p-137.0 zebra:1364.1p-536.3 }
	//== bovidae win % = 0.813600 ==

}
