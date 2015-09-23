package controllers

import "github.com/astaxie/beego"

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
	//	create user
	//	create session
}
