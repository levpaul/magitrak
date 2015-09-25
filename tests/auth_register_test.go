package test

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	_ "github.com/astaxie/beego/session/mysql"
	_ "github.com/levilovelock/magitrak/routers"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestInvalidJSONReturns400
// TestTooSmallPasswordReturns400
// TestTooLongPasswordReturns400
// TestInvalidEmailPasswordReturns400
// TestDuplicateEmailReturns400

func TestGETReturns404(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/auth/register", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestGet", "Code[%d]\n%s", w.Code, w.Body.String())

	assert.Equal(t, 404, w.Code)
}