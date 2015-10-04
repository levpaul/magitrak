package controllers

import (
	"encoding/json"

	"gopkg.in/olivere/elastic.v2"

	"github.com/levilovelock/magitrak/models"

	"github.com/astaxie/beego"
)

// Operations about match
type MatchController struct {
	beego.Controller
}

// @router / [post]
func (m *MatchController) Create() {
	var match models.Match
	unmarshalErr := json.Unmarshal(m.Ctx.Input.RequestBody, &match)
	if unmarshalErr != nil {
		beego.Debug("Error unmarshalling request for Match POST:", unmarshalErr.Error())
		m.Abort("400")
	}

	validationErr := match.Validate()
	if validationErr != nil {
		beego.Debug("Error validating match for POST:", validationErr.Error())
		m.Abort("400")
	}

	// Create a client
	client, elasticClientErr := elastic.NewClient()
	if elasticClientErr != nil {
		beego.Debug("Error connecting to ElasticSearch:", elasticClientErr.Error())
		m.Abort("500")
	}

	magiSession := m.GetSession(models.SESSION_NAME)
	if magiSession == nil {
		beego.Debug("Failed to find valid session for match creation request")
		m.Abort("500")
	}

	if match.UserId != magiSession.(models.MagiSession).UserId {
		beego.Debug("Session userid does not match the userId in match data for match creation request")
		m.Abort("400")
	}

	matchData, marshalErr := json.Marshal(match)
	if marshalErr != nil {
		beego.Debug("Failed to marsahl match object:", marshalErr)
	}

	_, elasticInsertErr := client.Index().
		Index(models.ELASTIC_INDEX).
		Type(models.ELASTIC_MATCH_TYPE).
		BodyJson(string(matchData)).
		Do()

	if elasticInsertErr != nil {
		beego.Debug("Failed to insert document into ElasticSearch:", elasticInsertErr.Error())
		m.Abort("501")
	}

	m.Data["json"] = map[string]string{"MatchId": "some-uuid"}
	m.ServeJson()
}

// @router /:matchId [get]
func (m *MatchController) GetSingle() {
	matchId := m.Ctx.Input.Params[":matchId"]
	if matchId != "" {
		match, err := models.GetOne(matchId)
		if err != nil {
			m.Data["json"] = err
		} else {
			m.Data["json"] = match
		}
	}
	m.ServeJson()
}

// @router /:matchId [delete]
func (m *MatchController) Delete() {
	matchId := m.Ctx.Input.Params[":matchId"]
	models.Delete(matchId)
	m.Data["json"] = "delete success!"
	m.ServeJson()
}

// @router / [get]
func (m *MatchController) GetAll() {
	// validate session
	// return all matches for user
}
