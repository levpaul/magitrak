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
	// parse json data
	user := models.User{}
	body := a.Ctx.Input.RequestBody
	jsonParseError := json.Unmarshal(body, &user)
	if jsonParseError != nil {
		beego.Debug("Failed to parse registration JSON load")
		a.Abort("400")
	}

	//	Check login details
	//	if good, create session data
}

// @router /logout [get]
func (a *AuthController) Logout() {
	//	GET auth/logout
	//	remove session data
}

// @router /register [post]
func (a *AuthController) Register() {
	//	validate params
	type registration struct {
		Email    string
		Password string
	}
	form := registration{}
	body := a.Ctx.Input.RequestBody
	jsonParseErr := json.Unmarshal(body, &form)
	if jsonParseErr != nil {
		beego.Debug("Failed to parse registration JSON load")
		a.Abort("400")
	}

	//	create user
	newUser, addUserErr := models.AddNewUser(form.Email, form.Password)
	if addUserErr != nil {
		beego.Debug("Error occured during registration of new user:", addUserErr.Error())
		a.Abort("400")
	}

	a.Data["json"] = newUser
	a.ServeJson()
}
