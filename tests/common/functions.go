package common

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/astaxie/beego"
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
