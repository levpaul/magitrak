package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	_ "github.com/levilovelock/magitrak/routers"
	_ "github.com/mattn/go-sqlite3"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
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

func TestMatchGETNoLoginReturns401(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/match/1", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 302, w.Code)
	assert.Equal(t, "/v1/auth/unauthorised", w.Header().Get("Location"))
}

func TestMatchGETWithLoginReturns200(t *testing.T) {
	body := []byte(`{"email": "some@email.com", "password":"validpassword"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	// Get session cookie from first login
	resp := http.Response{Header: w.HeaderMap}
	cookies := resp.Cookies()

	r, _ = http.NewRequest("GET", "/v1/match/1", nil)
	w = httptest.NewRecorder()
	r.AddCookie(cookies[0])

	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}

func TestAuthUnauthorised401(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/auth/unauthorised", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "Unauthorised", w.Body.String())
}
