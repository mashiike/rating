package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"

	"github.com/jszwec/csvutil"
	"github.com/mashiike/golicko"
	"github.com/pkg/errors"
)

type player struct {
	ID string `csv:"player_id"`
	golicko.Rating
}

type match struct {
	PlayerID         string        `csv:"player_id"`
	OpponentPlayerID string        `csv:"opponent_player_id"`
	Score            golicko.Score `csv:"score"`
}

func main() {
	var (
		playersFile = flag.String("p", "./example_players.csv", "players list")
		matchesFile = flag.String("m", "./example_matches.csv", "matches list")
		outputFile  = flag.String("o", "./output.csv", "after updated players list")
		tau         = flag.Float64("tau", 0.5, "tau value (0.3 ~ 1.6) default 0.5")
	)
	flag.Parse()

	players, err := readPlayers(*playersFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	matches, err := readMatches(*matchesFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	updated := updateRatings(players, matches, *tau)
	if err := writePlayers(*outputFile, updated); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func updateRatings(players []player, matches []match, tau float64) []player {
	ratings := make(map[string]golicko.Rating, len(players))
	results := make(map[string][]golicko.Result, len(players))
	for _, p := range players {
		ratings[p.ID] = p.Rating
		results[p.ID] = make([]golicko.Result, 0, 10)
	}

	for _, m := range matches {
		if _, ok := ratings[m.PlayerID]; !ok {
			ratings[m.PlayerID] = golicko.DefaultRating
			results[m.PlayerID] = make([]golicko.Result, 0, 10)
		}
		if _, ok := ratings[m.OpponentPlayerID]; !ok {
			ratings[m.OpponentPlayerID] = golicko.DefaultRating
			results[m.OpponentPlayerID] = make([]golicko.Result, 0, 10)
		}
		results[m.PlayerID] = append(results[m.PlayerID], golicko.Result{
			Opponent: ratings[m.OpponentPlayerID],
			Score:    m.Score,
		})
		results[m.OpponentPlayerID] = append(results[m.OpponentPlayerID], golicko.Result{
			Opponent: ratings[m.PlayerID],
			Score:    m.Score.Opponent(),
		})
	}

	setting := golicko.Setting{Tau: tau}
	ret := make([]player, 0, len(ratings))
	for id, r := range ratings {
		ret = append(ret, player{
			ID:     id,
			Rating: r.Update(results[id], setting),
		})
	}
	return ret
}

func readPlayers(filename string) ([]player, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "file open failed")
	}
	defer file.Close()

	dec, err := csvutil.NewDecoder(csv.NewReader(file))
	if err != nil {
		return nil, errors.Wrap(err, "decode player csv failed")
	}

	var players []player
	var p player
	err = dec.Decode(&p)
	for err != io.EOF {
		if err != nil {
			return nil, errors.Wrap(err, "player decode line failed")
		}
		players = append(players, p)
		err = dec.Decode(&p)
	}
	return players, nil
}

func readMatches(filename string) ([]match, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "file open failed")
	}
	defer file.Close()

	dec, err := csvutil.NewDecoder(csv.NewReader(file))
	if err != nil {
		return nil, errors.Wrap(err, "decode match csv failed")
	}

	var matches []match
	var m match
	err = dec.Decode(&m)
	for err != io.EOF {
		if err != nil {
			return nil, errors.Wrap(err, "match decode line failed")
		}
		matches = append(matches, m)
		err = dec.Decode(&m)
	}
	return matches, nil
}

func writePlayers(filename string, players []player) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	enc := csvutil.NewEncoder(writer)
	for _, p := range players {
		enc.Encode(p)
	}
	writer.Flush()
	return nil
}
