package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego/orm"
	_ "github.com/levilovelock/magitrak/routers"
	_ "github.com/mattn/go-sqlite3"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/suite"
)

type AuthLogoutTestSuite struct {
	suite.Suite
}

func TestAuthLogoutTestSuite(t *testing.T) {
	suite.Run(t, new(AuthLogoutTestSuite))
}

func (s *AuthLogoutTestSuite) SetupSuite() {
	beego.TestBeegoInit("../../../levilovelock/magitrak")

	dbAddress := beego.AppConfig.String("modelORMPrepopulatedAdress")
	dbType := beego.AppConfig.String("modelORMdb")

	dbErr := orm.RegisterDataBase("default", dbType, dbAddress, 30)
	if dbErr != nil {
		beego.Error(dbErr)
	}
}

func (s *AuthLogoutTestSuite) TestAuthLogoutPOSTReturns404() {
	r, _ := http.NewRequest("POST", "/v1/auth/logout", bytes.NewBuffer([]byte("")))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(404, w.Code)
}

func (s *AuthLogoutTestSuite) TestAuthLogoutValidLoginThenLogoutThenLogin200() {
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

	s.Assert().Equal(200, w.Code)

	// Login for the second time
	loginRequest.AddCookie(session)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, loginRequest)

	resp = http.Response{Header: w.HeaderMap}

	s.Assert().Equal(0, len(resp.Cookies()))
	s.Assert().Equal(200, w.Code)
	s.Assert().Equal("", w.Body.String())
}
