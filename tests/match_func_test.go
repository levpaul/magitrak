package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/levilovelock/magitrak/models"
	_ "github.com/levilovelock/magitrak/routers"
	"github.com/levilovelock/magitrak/tests/common"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"

	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type MatchFuncTestSuite struct {
	suite.Suite
}

func TestMatchTestSuite(t *testing.T) {
	suite.Run(t, new(MatchFuncTestSuite))
}

func (s *MatchFuncTestSuite) SetupSuite() {
	beego.TestBeegoInit("../../../levilovelock/magitrak")

	dbAddress := beego.AppConfig.String("modelORMPrepopulatedAdress")
	dbType := beego.AppConfig.String("modelORMdb")

	dbErr := orm.RegisterDataBase("default", dbType, dbAddress, 30)
	if dbErr != nil {
		beego.Error(dbErr)
	}
}

func (s *MatchFuncTestSuite) TestMatchPOSTWithInvalidJSONReturns400() {
	body := []byte(`{"m"___,,L"'...aalidpassword"}`)

	r, _ := http.NewRequest("POST", "/v1/match", bytes.NewBuffer(body))
	r.AddCookie(common.GetValidLoggedInSessionCookie())
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(400, w.Code)
}

func (s *MatchFuncTestSuite) TestMatchPOSTWithDifferentUserIdInMatchThanSessionReturns400() {
	body := []byte(`{"userid": 4}`)

	r, _ := http.NewRequest("POST", "/v1/match", bytes.NewBuffer(body))
	r.AddCookie(common.GetValidLoggedInSessionCookie())
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(400, w.Code)
}

func (s *MatchFuncTestSuite) TestMatchPOSTValidUserIdAndMatchReturns200() {
	bodyObject := GetValidMatch()
	body, _ := json.Marshal(bodyObject)

	r, _ := http.NewRequest("POST", "/v1/match", bytes.NewBuffer(body))
	r.AddCookie(common.GetValidLoggedInSessionCookie())
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(200, w.Code)
}

func (s *MatchFuncTestSuite) TestMatchGETNoLoginReturns401() {
	r, _ := http.NewRequest("GET", "/v1/match/1", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(302, w.Code)
	s.Assert().Equal("/v1/auth/unauthorised", w.Header().Get("Location"))
}

func (s *MatchFuncTestSuite) TestMatchGETWithLoginReturns200() {
	r, _ := http.NewRequest("GET", "/v1/match/1", nil)
	w := httptest.NewRecorder()
	r.AddCookie(common.GetValidLoggedInSessionCookie())

	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(200, w.Code)
}

func (s *MatchFuncTestSuite) TestMatchPOSTNoPlayerDeckReturns400() {
	bodyObject := GetValidMatch()
	bodyObject.OpponentDeck = ""

	body, _ := json.Marshal(bodyObject)

	r, _ := http.NewRequest("POST", "/v1/match", bytes.NewBuffer(body))
	r.AddCookie(common.GetValidLoggedInSessionCookie())
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	s.Assert().Equal(400, w.Code)
}

func GetValidMatch() models.Match {
	return models.Match{UserId: 1, PlayerDeck: "burn", OpponentDeck: "bloom"}
}

// Match Addition Tests
// TestMatchPOSTNoOpponentDeckReturns400
// TestMatchPOSTNoDateReturns400

// Match Get Tests
// TestMatchGETInvalidIdReturns404

// Match Delete Tests
// TestMatchDELETEInvalidIDReturns404
// TestMatchDELETEDifferentUserIDReturns400
// TestMatchDELETEInvalidJSONReturns400

// Match Complete Tests
// TestMatchInsertRetrievalASuccess
// TestMatchInsertRetrievalBSuccess
// TestMatchInsertDeleteRetrieveSuccess

// Match Retrieval Tests
// TestMatchesGETInvalidJSON400
// TestMatchesGETIncorrectUserID400
// TestMatchesGETReturnsArray
// TestMatchesInsertTwoThenGetReturnsAtLeastTwoSuccess

// Then add some unit tests for match model funcs
