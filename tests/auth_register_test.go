package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego/orm"
	"github.com/levilovelock/magitrak/models"
	_ "github.com/levilovelock/magitrak/routers"
	_ "github.com/mattn/go-sqlite3"

	"bytes"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/suite"
)

type AuthRegisterTestSuite struct {
	suite.Suite
}

func TestAuthRegisterTestSuite(t *testing.T) {
	suite.Run(t, new(AuthRegisterTestSuite))
}

func (s *AuthRegisterTestSuite) SetupSuite() {
	beego.TestBeegoInit("../../../levilovelock/magitrak")

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

func (s *AuthRegisterTestSuite) TestAuthRegisterGETReturns404() {
	r, _ := http.NewRequest("GET", "/v1/auth/register", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(404, w.Code)
}

func (s *AuthRegisterTestSuite) TestAuthRegisterInvalidJSONReturns400() {
	body := []byte(`"{incomplete and not "valid JSON`)
	r, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(400, w.Code)
}

func (s *AuthRegisterTestSuite) TestAuthRegistertTooSmallPasswordReturns400() {
	body := []byte(`{"email":"asfd@gmail.com", "password":"small"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(400, w.Code)
}

func (s *AuthRegisterTestSuite) TestAuthRegisterTooLongPasswordReturns400() {
	body := []byte(`{"email":"asfd@gmail.com", "password":"reallyreallyreallyreallyreally
	reallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallyreallylongpassword"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(400, w.Code)
}

func (s *AuthRegisterTestSuite) TestAuthRegisterInvalidEmailPasswordReturns400() {
	body := []byte(`{"email":"asfdnotemail", "password":"validpassword"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(400, w.Code)
}

func (s *AuthRegisterTestSuite) TestAuthRegisterValidRegistrationReturns200() {
	body := []byte(`{"email":"some@email.com", "password":"validpassword"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(200, w.Code)
	cleanupRegisterTest()
}

func (s *AuthRegisterTestSuite) TestAuthRegisterRegisterSameEmailTwiceReturns400() {
	body := []byte(`{"email":"some@otheremail.com", "password":"validpassword"}`)
	r, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	w = httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(400, w.Code)
}

func cleanupRegisterTest() {
	user, err := models.GetUserByEmail("some@email.com")
	user.Delete()
}
