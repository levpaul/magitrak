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
	"github.com/stretchr/testify/assert"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	SESSION_USER_ID = 1
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

func TestMatchPOSTWithInvalidMatchReturns400(t *testing.T) {
	body := []byte(`{"m"___,,L"'...aalidpassword"}`)

	r, _ := http.NewRequest("POST", "/v1/match", bytes.NewBuffer(body))
	r.AddCookie(getValidLoggedInSessionCookie())
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 400, w.Code)
}

func TestMatchPOSTWithDifferentUserIdInMatchThanSessionReturns400(t *testing.T) {
	body := []byte(`{"userid": 4}`)

	r, _ := http.NewRequest("POST", "/v1/match", bytes.NewBuffer(body))
	r.AddCookie(getValidLoggedInSessionCookie())
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 400, w.Code)
}

func TestMatchPOSTValidUserIdAndMatchReturns200(t *testing.T) {
	body := []byte(`{"userid": 1}`)

	r, _ := http.NewRequest("POST", "/v1/match", bytes.NewBuffer(body))
	r.AddCookie(getValidLoggedInSessionCookie())
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}

func getValidLoggedInSessionCookie() *http.Cookie {
	body := []byte(`{"email": "some@email.com", "password":"validpassword"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	resp := http.Response{Header: w.HeaderMap}
	cookies := resp.Cookies()
	return cookies[0]
}
