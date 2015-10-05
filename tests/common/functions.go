package common

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/astaxie/beego"
)

const (
	SESSION_USER_ID = 1
)

func GetValidLoggedInSessionCookie() *http.Cookie {
	body := []byte(`{"email": "some@email.com", "password":"validpassword"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	resp := http.Response{Header: w.HeaderMap}
	cookies := resp.Cookies()
	return cookies[0]
}

func GetValidLoggedInSessionCookieOtherUser() *http.Cookie {
	body := []byte(`{"email": "other@gmail.com", "password":"pepepepe"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	resp := http.Response{Header: w.HeaderMap}
	cookies := resp.Cookies()
	return cookies[0]
}
