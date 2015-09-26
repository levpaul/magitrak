package test

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/astaxie/beego/orm"
	_ "github.com/levilovelock/magitrak/routers"
	_ "github.com/mattn/go-sqlite3"

	"bytes"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)

	dbAddress := beego.AppConfig.String("modelORMaddress")
	dbType := beego.AppConfig.String("modelORMdb")

	dbErr := orm.RegisterDataBase("default", dbType, dbAddress, 30)
	if dbErr != nil {
		beego.Error(dbErr)
	}

	err := orm.RunSyncdb("default", true, true)
	if err != nil {
		beego.Error(err)
	}
}

func TestAuthRegisterGETReturns404(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/auth/register", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 404, w.Code)
}

func TestAuthRegisterInvalidJSONReturns400(t *testing.T) {
	body := []byte(`"{incomplete and not "valid JSON`)
	r, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 400, w.Code)
}

func TesAuthRegistertTooSmallPasswordReturns400(t *testing.T) {
	body := []byte(`{"email":"asfd@gmail.com", "password":"small"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 400, w.Code)
}

func TestAuthRegisterTooLongPasswordReturns400(t *testing.T) {
	body := []byte(`{"email":"asfd@gmail.com", "password":"reallyreallyreallyreallyreally
	reallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallylongpassword"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 400, w.Code)
}

func TestAuthRegisterInvalidEmailPasswordReturns400(t *testing.T) {
	body := []byte(`{"email":"asfdnotemail", "password":"validpassword"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 400, w.Code)
}

func TestAuthRegisterValidRegistrationReturns200(t *testing.T) {
	body := []byte(`{"email":"some@email.com", "password":"validpassword"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}

func TestAuthRegisterRegisterSameEmailTwiceReturns400(t *testing.T) {
	body := []byte(`{"email":"some@otheremail.com", "password":"validpassword"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 400, w.Code)
}
