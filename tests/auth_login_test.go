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

func TestAuthLoginGETReturns404(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/auth/login", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 404, w.Code)
}

func TestAuthLoginInvalidJSON400(t *testing.T) {
	body := []byte(`{"email":"asfd@gmail.as",,,'':':sword":"small"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 400, w.Code)
}

func TestAuthLoginNoPassword401(t *testing.T) {
	body := []byte(`{"email":"asfd@gmail."small"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 400, w.Code)
}

// TestAuthLoginNoUsername401
// TestAuthLoginNoUser401
// TestAuthLoginInvalidPassword401
// TestAuthLoginValid200
// TestAuthLoginAlreadyLoggedIn200
