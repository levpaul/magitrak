package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/astaxie/beego/orm"
	_ "github.com/levilovelock/magitrak/routers"
	_ "github.com/mattn/go-sqlite3"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)

	dbAddress := beego.AppConfig.String("modelORMPrepopulatedAdress")
	dbType := beego.AppConfig.String("modelORMdb")

	dbErr := orm.RegisterDataBase("default", dbType, dbAddress, 30)
	if dbErr != nil {
		beego.Error(dbErr)
	}
}

func TestAuthLogoutPOSTReturns404(t *testing.T) {
	r, _ := http.NewRequest("POST", "/v1/auth/logout", bytes.NewBuffer([]byte("")))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 404, w.Code)
}

func TestAuthLogoutValidLoginThenLogoutThenLogin200(t *testing.T) {
	body := []byte(`{"email": "some@email.com", "password":"validpassword"}`)
	loginRequest, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, loginRequest)

	// Add session cookie from first login to request
	resp := http.Response{Header: w.HeaderMap}
	session := resp.Cookies()[0]
	logoutRequest, _ := http.NewRequest("GET", "/v1/auth/logout", nil)
	logoutRequest.AddCookie(session)

	// Logout
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, logoutRequest)

	assert.Equal(t, 200, w.Code)

	// Login for the second time
	loginRequest.AddCookie(session)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, loginRequest)

	resp = http.Response{Header: w.HeaderMap}

	assert.Equal(t, 0, len(resp.Cookies()))
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Body.String())
}
