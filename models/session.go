package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

const (
	SESSION_NAME = "magitrak"
)

type MagiSession struct {
	Authenticated bool
	UserId        int
}

func IsLoggedInFilter(c *context.Context) {
	authRequired, confErr := beego.AppConfig.Bool("SessionAuthRequired")
	if confErr != nil {
		c.Abort(501, "Internal Server Error")
	}
	if !authRequired {
		return
	}

	sessionIntf := c.Input.Session(SESSION_NAME)
	if sessionIntf == nil || !sessionIntf.(MagiSession).Authenticated || sessionIntf.(MagiSession).UserId == 0 {
		c.Redirect(302, "/v1/auth/unauthorised")
	}
}
