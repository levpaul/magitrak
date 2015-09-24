package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/levilovelock/magitrak/models"
)

// Operations about auth
type AuthController struct {
	beego.Controller
}

// @router /login [post]
func (a *AuthController) Login() {
	//	Check login details
	//	if good, create session data
}

// @router /logout [get]
func (m *AuthController) Logout() {
	//	GET auth/logout
	//	remove session data
}

// @router /register [post]
func (m *AuthController) Register() {
	//	validate params
	type registration struct {
		Email    string
		Password string
	}
	form := registration{}
	body := m.Ctx.Input.RequestBody
	jsonParseErr := json.Unmarshal(body, &form)
	if jsonParseErr != nil {
		beego.Debug("Failed to parse registration JSON load")
		m.Abort("400")
	}

	//	create user
	newUser, addUserErr := models.AddNewUser(form.Email, form.Password)
	if addUserErr != nil {
		beego.Debug("Error occured during registration of new user:", addUserErr.Error())
		m.Abort("400")
	}

	m.Data["json"] = newUser
	m.ServeJson()
}
