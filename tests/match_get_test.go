package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/levilovelock/magitrak/routers"
	_ "github.com/mattn/go-sqlite3"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/stretchr/testify/suite"
)

type MatchGetTestSuite struct {
	suite.Suite
}

func TestMatchGetTestSuite(t *testing.T) {
	suite.Run(t, new(MatchGetTestSuite))
}

func (s *MatchGetTestSuite) SetupSuite() {
	beego.TestBeegoInit("../../../levilovelock/magitrak")

	dbAddress := beego.AppConfig.String("modelORMPrepopulatedAdress")
	dbType := beego.AppConfig.String("modelORMdb")

	dbErr := orm.RegisterDataBase("default", dbType, dbAddress, 30)
	if dbErr != nil {
		beego.Error(dbErr)
	}
}

func (s *MatchGetTestSuite) TestMatchGETNoLoginReturns401() {
	r, _ := http.NewRequest("GET", "/v1/match/1", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(302, w.Code)
	s.Assert().Equal("/v1/auth/unauthorised", w.Header().Get("Location"))
}

func (s *MatchGetTestSuite) TestMatchGETWithLoginReturns200() {
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

	s.Assert().Equal(200, w.Code)
}

func (s *MatchGetTestSuite) TestAuthUnauthorised401() {
	r, _ := http.NewRequest("GET", "/v1/auth/unauthorised", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(401, w.Code)
	s.Assert().Equal("Unauthorised", w.Body.String())
}
