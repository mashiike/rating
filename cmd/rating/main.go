package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/jszwec/csvutil"
	"github.com/mashiike/golicko"
	"github.com/pkg/errors"
)

type player struct {
	ID string `csv:"player_id"`
	*golicko.Player
}

type match struct {
	LeftID    string    `csv:"left_id"`
	RightID   string    `csv:"right_id"`
	WinnerID  string    `csv:"winner_id"`
	OutcomeAt time.Time `csv:"outcome_at"`
}

func main() {
	var (
		playersFile = flag.String("p", "./example_players.json", "players list")
		matchesFile = flag.String("m", "./example_matches.csv", "matches list")
		outputFile  = flag.String("o", "./output.json", "after updated players list")
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
	config := golicko.NewConfig()
	loc, _ := time.LoadLocation("Asia/Tokyo")
	baseTime := time.Date(2019, 05, 01, 0, 0, 0, 0, loc)
	config.Tau = tau
	playersMap := make(map[string]*golicko.Player, len(players))
	for _, p := range players {
		playersMap[p.ID] = p.Player
	}
	for _, m := range matches {
		left, ok := playersMap[m.LeftID]
		if !ok {
			left = golicko.NewPlayer(baseTime, config)
			playersMap[m.LeftID] = left

		}
		right, ok := playersMap[m.RightID]
		if !ok {
			right = golicko.NewPlayer(baseTime, config)
			playersMap[m.RightID] = right

		}
		winner, ok := playersMap[m.WinnerID]
		if !ok {
			winner = nil
		}
		match := &golicko.Match{
			Left:      left,
			Right:     right,
			Winner:    winner,
			OutcomeAt: m.OutcomeAt,
		}
		if err := match.Apply(config); err != nil {
			panic(err)
		}
	}

	ret := make([]player, 0, len(playersMap))
	for id, p := range playersMap {
		ret = append(ret, player{
			ID:     id,
			Player: p,
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

	scanner := bufio.NewScanner(file)
	var buf bytes.Buffer
	buf.WriteByte('[')
	lineCount := 0

	for scanner.Scan() {
		b := scanner.Bytes()

		buf.Write(b)
		buf.WriteByte(',')

		lineCount += 1
	}
	if lineCount == 0 {
		buf.WriteByte(',')
	}

	data := buf.Bytes()
	data[len(data)-1] = ']'
	players := make([]player, lineCount)
	if err := json.Unmarshal(data, &players); err != nil {
		return nil, err
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
	enc := json.NewEncoder(file)
	for _, p := range players {
		enc.Encode(p)
	}
	return nil
}
