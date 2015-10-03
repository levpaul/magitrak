package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/levilovelock/magitrak/routers"
	"github.com/levilovelock/magitrak/tests/common"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	SESSION_USER_ID = 1
)

type MatchTestSuite struct {
	suite.Suite
}

func TestMatchTestSuite(t *testing.T) {
	suite.Run(t, new(MatchTestSuite))
}

func (s *MatchTestSuite) SetupSuite() {
	beego.TestBeegoInit("../../../levilovelock/magitrak")

	dbAddress := beego.AppConfig.String("modelORMPrepopulatedAdress")
	dbType := beego.AppConfig.String("modelORMdb")

	dbErr := orm.RegisterDataBase("default", dbType, dbAddress, 30)
	if dbErr != nil {
		beego.Error(dbErr)
	}
}

func (s *MatchTestSuite) TestMatchPOSTWithInvalidMatchReturns400() {
	body := []byte(`{"m"___,,L"'...aalidpassword"}`)

	r, _ := http.NewRequest("POST", "/v1/match", bytes.NewBuffer(body))
	r.AddCookie(common.GetValidLoggedInSessionCookie())
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(400, w.Code)
}

func (s *MatchTestSuite) TestMatchPOSTWithDifferentUserIdInMatchThanSessionReturns400() {
	body := []byte(`{"userid": 4}`)

	r, _ := http.NewRequest("POST", "/v1/match", bytes.NewBuffer(body))
	r.AddCookie(common.GetValidLoggedInSessionCookie())
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(400, w.Code)
}

func (s *MatchTestSuite) TestMatchPOSTValidUserIdAndMatchReturns200() {
	body := []byte(`{"userid": 1}`)

	r, _ := http.NewRequest("POST", "/v1/match", bytes.NewBuffer(body))
	r.AddCookie(common.GetValidLoggedInSessionCookie())
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(200, w.Code)
}

func (s *MatchTestSuite) TestMatchGETNoLoginReturns401() {
	r, _ := http.NewRequest("GET", "/v1/match/1", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(302, w.Code)
	s.Assert().Equal("/v1/auth/unauthorised", w.Header().Get("Location"))
}

func (s *MatchTestSuite) TestMatchGETWithLoginReturns200() {
	r, _ := http.NewRequest("GET", "/v1/match/1", nil)
	w := httptest.NewRecorder()
	r.AddCookie(common.GetValidLoggedInSessionCookie())

	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(200, w.Code)
}

// TestMatchPOSTNoPlayerDeckReturns400
// TestMatchPOSTNoOpponentDeckReturns400
// TestMatchPOSTNoDateReturns400
// TestMatchPOSTValidSampleMatchAReturn200AndGetAReturn200
// TestMatchPOSTValidSampleMatchBReturn200AndGetBReturn200
