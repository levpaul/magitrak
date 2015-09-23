package controllers

import (
	"encoding/json"

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
	json.Unmarshal(m.Ctx.Input.RequestBody, &match)
	// TODO: Generate a Match ID (uuid)
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
