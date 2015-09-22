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

// @Title create
// @Description create match
// @Param	body		body 	models.Match	true		"The match content"
// @Success 200 {string} models.Match.Id
// @Failure 403 body is empty
// @router / [post]
func (m *MatchController) Post() {
	var match models.Match
	json.Unmarshal(m.Ctx.Input.RequestBody, &match)
	// TODO: Generate a Match ID (uuid)
	m.Data["json"] = map[string]string{"MatchId": "some-uuid"}
	m.ServeJson()
}

// @Title Get
// @Description find match by matchid
// @Param	matchId		path 	string	true		"the matchid you want to get"
// @Success 200 {match} models.Match
// @Failure 403 :matchId is empty
// @router /:matchId [get]
func (m *MatchController) Get() {
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

// @Title delete
// @Description delete the object
// @Param	objectId		path 	string	true		"The objectId you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 objectId is empty
// @router /:objectId [delete]
func (m *MatchController) Delete() {
	matchId := m.Ctx.Input.Params[":matchId"]
	models.Delete(matchId)
	m.Data["json"] = "delete success!"
	m.ServeJson()
}
