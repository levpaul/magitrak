package models

import (
	"errors"
	"time"
)

var (
	Matches map[string]*Match
)

type Match struct {
	ObjectId         string
	Date             time.Time
	PlayerDeck       string
	OpponentDeck     string
	Win              bool
	Reason           string
	Sideboard        bool
	PlayFirst        bool
	StartingHandSize int
	LandsInOpener    int
	OpponentName     string
	Notes            string
}

func init() {
	Matches = make(map[string]*Match)
}

func GetOne(matchId string) (*Match, error) {
	match := Matches[matchId]
	if match == nil {
		return nil, errors.New("No match found with id: " + matchId)
	}
	return match, nil
}

func Delete(objectId string) {
	delete(Matches, objectId)
}
