package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/levilovelock/magitrak/routers"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	SESSION_USER_ID = 1
)

type MatchAddTestSuite struct {
	suite.Suite
}

func TestMatchAddTestSuite(t *testing.T) {
	suite.Run(t, new(MatchAddTestSuite))
}

func (s *MatchAddTestSuite) SetupSuite() {
	beego.TestBeegoInit("../../../levilovelock/magitrak")

	dbAddress := beego.AppConfig.String("modelORMPrepopulatedAdress")
	dbType := beego.AppConfig.String("modelORMdb")

	dbErr := orm.RegisterDataBase("default", dbType, dbAddress, 30)
	if dbErr != nil {
		beego.Error(dbErr)
	}
}

func (s *MatchAddTestSuite) TestMatchPOSTWithInvalidMatchReturns400() {
	body := []byte(`{"m"___,,L"'...aalidpassword"}`)

	r, _ := http.NewRequest("POST", "/v1/match", bytes.NewBuffer(body))
	r.AddCookie(getValidLoggedInSessionCookie())
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(400, w.Code)
}

func (s *MatchAddTestSuite) TestMatchPOSTWithDifferentUserIdInMatchThanSessionReturns400() {
	body := []byte(`{"userid": 4}`)

	r, _ := http.NewRequest("POST", "/v1/match", bytes.NewBuffer(body))
	r.AddCookie(getValidLoggedInSessionCookie())
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(400, w.Code)
}

func (s *MatchAddTestSuite) TestMatchPOSTValidUserIdAndMatchReturns200() {
	body := []byte(`{"userid": 1}`)

	r, _ := http.NewRequest("POST", "/v1/match", bytes.NewBuffer(body))
	r.AddCookie(getValidLoggedInSessionCookie())
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(200, w.Code)
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
