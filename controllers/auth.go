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
	// Check if already logged in
	existingSession := a.GetSession(models.SESSION_NAME)
	beego.Debug("Existing session:", existingSession)
	if existingSession != nil && existingSession.(models.MagiSession).Authenticated {
		beego.Debug("Tried to login when already logged in")
		a.Ctx.ResponseWriter.Write([]byte("Already logged in!"))
		return
	}

	// parse json data
	user := models.User{}
	body := a.Ctx.Input.RequestBody
	jsonParseError := json.Unmarshal(body, &user)
	if jsonParseError != nil {
		beego.Debug("Failed to parse registration JSON load")
		a.Abort("400")
	}

	userId, authErr := user.Authenticate()
	if authErr != nil {
		beego.Debug("Failed to authenticate user for login")
		a.Abort("401")
	}

	a.SetSession(models.SESSION_NAME, models.MagiSession{UserId: userId, Authenticated: true})
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

	type response struct {
		Id int
	}
	a.Data["json"] = response{Id: newUser.Id}
	a.ServeJson()
}
