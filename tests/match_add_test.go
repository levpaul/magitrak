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

	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
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

	os.Setenv("SessionAuthRequired", "false")
}

func TestMatchPOSTWithValidMatchReturns200(t *testing.T) {
	body := []byte(`{"match": "somematchcom", "password":"validpassword"}`)

	r, _ := http.NewRequest("POST", "/v1/match", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}
