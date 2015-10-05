package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/astaxie/beego"

	"gopkg.in/olivere/elastic.v2"
)

const (
	ELASTIC_MATCH_TYPE = "match"
	ELASTIC_INDEX      = "magitrak"

	NO_MATCH_FOUND_ERROR = "No match found"
)

var (
	Matches map[string]*Match
)

type Match struct {
	UserId           int
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

func InsertMatch(m Match) (string, error) {
	client, elasticClientErr := elastic.NewClient()
	if elasticClientErr != nil {
		return "", elasticClientErr
	}

	matchData, marshalErr := json.Marshal(m)
	if marshalErr != nil {
		return "", marshalErr
	}

	matchId, elasticInsertErr := client.Index().
		Index(ELASTIC_INDEX).
		Type(ELASTIC_MATCH_TYPE).
		BodyJson(string(matchData)).
		Do()

	if elasticInsertErr != nil {
		return "", elasticInsertErr
	}

	if matchId.Created == false || matchId.Id == "" {
		return "", errors.New("Match was not created")
	}

	return matchId.Id, nil
}

func GetOne(matchId string) (*Match, error) {
	client, elasticClientErr := elastic.NewClient()
	if elasticClientErr != nil {
		return nil, elasticClientErr
	}

	matchResult, elasticSearchErr := client.Get().
		Index(ELASTIC_INDEX).
		Type(ELASTIC_MATCH_TYPE).
		Id(matchId).
		Do()

	if elasticSearchErr != nil {
		return nil, elasticSearchErr
	}

	if !matchResult.Found {
		return nil, errors.New(NO_MATCH_FOUND_ERROR)
	}

	match := &Match{}

	sourceData, unmarshalErr := matchResult.Source.MarshalJSON()
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	unmarshalErr = json.Unmarshal(sourceData, match)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return match, nil
}

func Delete(matchId string) bool {
	client, elasticClientErr := elastic.NewClient()
	if elasticClientErr != nil {
		beego.Debug("Error deleting match: ", matchId, " - error:", elasticClientErr.Error())
		return false
	}

	deleteResult, elasticSearchErr := client.Delete().
		Index(ELASTIC_INDEX).
		Type(ELASTIC_MATCH_TYPE).
		Id(matchId).
		Do()

	if elasticSearchErr != nil {
		beego.Debug("Error deleting match: ", matchId, " - error:", elasticSearchErr.Error())
		return false
	}

	if !deleteResult.Found {
		beego.Debug("Error deleting match: ", matchId, " - error:", elasticClientErr.Error())
		return false
	}

	return true
}

func (m *Match) Validate() error {
	if m.OpponentDeck == "" {
		return errors.New("No opponent deck supplied")
	}
	if m.PlayerDeck == "" {
		return errors.New("No player deck supplied")
	}
	if m.Date.Equal(time.Time{}) {
		return errors.New("No valid date given for match")
	}
	return nil
}
