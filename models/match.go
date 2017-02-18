package models

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/astaxie/beego"

	elastic "gopkg.in/olivere/elastic.v5"
)

const (
	ELASTIC_MATCH_TYPE = "match"
	ELASTIC_INDEX      = "magitrak"

	NO_MATCH_FOUND_ERROR = "elastic: Error 404 (Not Found)"
)

var (
	eslogin, espass string
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

func initESCreds() {
	if eslogin == "" {
		eslogin = beego.AppConfig.String("ElasticAuthName")
	}

	if espass == "" {
		espass = beego.AppConfig.String("ElasticAuthPass")
	}
}

func InsertMatch(m Match) (string, error) {
	initESCreds()
	ctx := context.Background()
	client, elasticClientErr := elastic.NewClient(elastic.SetBasicAuth(eslogin, espass))
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
		Do(ctx)

	if elasticInsertErr != nil {
		return "", elasticInsertErr
	}

	if matchId.Created == false || matchId.Id == "" {
		return "", errors.New("Match was not created")
	}

	return matchId.Id, nil
}

func GetOne(matchId string) (*Match, error) {
	initESCreds()
	ctx := context.Background()
	client, elasticClientErr := elastic.NewClient(elastic.SetBasicAuth(eslogin, espass))

	if elasticClientErr != nil {
		return nil, elasticClientErr
	}

	matchResult, elasticSearchErr := client.Get().
		Index(ELASTIC_INDEX).
		Type(ELASTIC_MATCH_TYPE).
		Id(matchId).
		Do(ctx)

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

func GetAll(userId int) ([]*Match, error) {
	initESCreds()
	ctx := context.Background()
	client, elasticClientErr := elastic.NewClient(elastic.SetBasicAuth(eslogin, espass))

	if elasticClientErr != nil {
		return nil, elasticClientErr
	}

	termQuery := elastic.NewTermQuery("UserId", userId)
	matchResult, elasticSearchErr := client.Search().
		Index(ELASTIC_INDEX).
		Type(ELASTIC_MATCH_TYPE).
		Query(termQuery).
		Size(9999).
		Do(ctx)

	if elasticSearchErr != nil {
		return nil, elasticSearchErr
	}

	if matchResult.Hits.TotalHits == 0 {
		return []*Match{}, nil
	}

	matches := []*Match{}

	for _, hit := range matchResult.Hits.Hits {
		match := &Match{}
		sourceData, unmarshalErr := hit.Source.MarshalJSON()
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
		unmarshalErr = json.Unmarshal(sourceData, match)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
		matches = append(matches, match)
	}

	return matches, nil
}

func Delete(matchId string) bool {
	initESCreds()
	ctx := context.Background()
	client, elasticClientErr := elastic.NewClient(elastic.SetBasicAuth(eslogin, espass))

	if elasticClientErr != nil {
		beego.Debug("Error deleting match: ", matchId, " - error:", elasticClientErr.Error())
		return false
	}

	deleteResult, elasticSearchErr := client.Delete().
		Index(ELASTIC_INDEX).
		Type(ELASTIC_MATCH_TYPE).
		Id(matchId).
		Do(ctx)

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
