package golicko_test

import (
	"fmt"
	"sort"
	"time"

	"github.com/mashiike/golicko"
)

func ExampleMatch() {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	baseTime := time.Date(2019, 05, 01, 0, 0, 0, 0, loc)
	config := golicko.NewConfig()
	players := map[string]golicko.IPlayer{
		"sheep":   golicko.NewPlayer(baseTime, config),
		"goat":    golicko.NewPlayer(baseTime, config),
		"donkey":  golicko.NewPlayer(baseTime, config),
		"horse":   golicko.NewPlayer(baseTime, config),
		"giraffe": golicko.NewPlayer(baseTime, config),
	}

	matches := []golicko.Match{
		//1st week
		{players["sheep"], players["goat"], nil, baseTime.Add(time.Hour)},
		{players["sheep"], players["donkey"], players["sheep"], baseTime.Add(time.Hour)},
		{players["donkey"], players["goat"], players["goat"], baseTime.Add(time.Hour)},
		//2nd week
		{players["giraffe"], players["goat"], players["giraffe"], baseTime.AddDate(0, 0, 7).Add(time.Hour)},
		{players["giraffe"], players["horse"], players["giraffe"], baseTime.AddDate(0, 0, 7).Add(2.0 * time.Hour)},
		{players["goat"], players["horse"], players["horse"], baseTime.AddDate(0, 0, 7).Add(3.0 * time.Hour)},
		//3rd week
		{players["giraffe"], players["sheep"], players["giraffe"], baseTime.AddDate(0, 0, 14).Add(time.Hour)},
		{players["giraffe"], players["donkey"], players["giraffe"], baseTime.AddDate(0, 0, 14).Add(2.0 * time.Hour)},
		{players["sheep"], players["horse"], players["sheep"], baseTime.AddDate(0, 0, 14).Add(3.0 * time.Hour)},
		{players["donkey"], players["horse"], players["horse"], baseTime.AddDate(0, 0, 14).Add(3.0 * time.Hour)},
	}
	for _, m := range matches {
		if err := m.Apply(config); err != nil {
			fmt.Println(err)
		}
	}

	//sort by rating
	pairs := []struct {
		name string
		p    golicko.IPlayer
	}{
		{"sheep", players["sheep"]},
		{"donkey", players["donkey"]},
		{"goat", players["goat"]},
		{"horse", players["horse"]},
		{"giraffe", players["giraffe"]},
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].p.Rating().Strength() > (pairs[j].p.Rating()).Strength() })
	for i, pair := range pairs {
		fmt.Printf("%d. %s   \t:%s\n", i+1, pair.name, pair.p)
	}
	//Output:
	//1. giraffe   	:2046.3p-357.0(4/0/0)
	//2. sheep   	:1793.9p-323.9(2/1/1)
	//3. horse   	:1603.5p-349.7(2/2/0)
	//4. goat   	:1465.6p-362.9(1/2/1)
	//5. donkey   	:949.6p-340.1(0/4/0)
}
