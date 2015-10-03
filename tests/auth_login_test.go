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

type AuthLoginTestSuite struct {
	suite.Suite
}

func TestAuthLoginTestSuite(t *testing.T) {
	suite.Run(t, new(AuthLoginTestSuite))
}

func (s *AuthLoginTestSuite) SetupSuite() {
	beego.TestBeegoInit("../../../levilovelock/magitrak")

	dbAddress := beego.AppConfig.String("modelORMPrepopulatedAdress")
	dbType := beego.AppConfig.String("modelORMdb")

	dbErr := orm.RegisterDataBase("default", dbType, dbAddress, 30)
	if dbErr != nil {
		beego.Error(dbErr)
	}
}

func (s *AuthLoginTestSuite) TestAuthLoginGETReturns404() {
	r, _ := http.NewRequest("GET", "/v1/auth/login", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(404, w.Code)
}

func (s *AuthLoginTestSuite) TestAuthLoginInvalidJSON400() {
	body := []byte(`{"email":"asfd@gmail.as",,,'':':sword":"small"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(400, w.Code)
}

func (s *AuthLoginTestSuite) TestAuthLoginNoPassword401() {
	body := []byte(`{"email":"someemail@gmail.com"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(401, w.Code)
}

func (s *AuthLoginTestSuite) TestAuthLoginNoEmail401() {
	body := []byte(`{"password":"somepass"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(401, w.Code)
}

func (s *AuthLoginTestSuite) TestAuthLoginUnregisteredEmail401() {
	body := []byte(`{"email": "randomemailthatnotexist@email.com", "password":"somepass"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(401, w.Code)
}

func (s *AuthLoginTestSuite) TestAuthLoginValid200() {
	body := []byte(`{"email": "some@email.com", "password":"validpassword"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(200, w.Code)
}

func (s *AuthLoginTestSuite) TestAuthLoginInvalidPassword401() {
	body := []byte(`{"email": "some@email.com", "password":"invalidpassword"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(401, w.Code)
}

func (s *AuthLoginTestSuite) TestAuthLoginAlreadyLoggedIn200() {
	body := []byte(`{"email": "some@email.com", "password":"validpassword"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	// Add session cookie from first login to request
	resp := http.Response{Header: w.HeaderMap}
	cookies := resp.Cookies()
	r.AddCookie(cookies[0])

	// Rerun request (with session cookie)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(200, w.Code)
	s.Assert().Equal("Already logged in!", w.Body.String())
}
