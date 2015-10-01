package models

import (
	"os"

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
	// TODO: Fix this workaround with appconfig (mainly for testing)
	authRequired := os.Getenv("sessionAuthRequired")
	if authRequired == "false" {
		return
	}

	sessionIntf := c.Input.Session(SESSION_NAME)
	if sessionIntf == nil || !sessionIntf.(MagiSession).Authenticated || sessionIntf.(MagiSession).UserId == 0 {
		c.Redirect(302, "/v1/auth/unauthorised")
	}
}
