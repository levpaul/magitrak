// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/levilovelock/magitrak/controllers"
	"github.com/levilovelock/magitrak/models"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/match",
			beego.NSInclude(
				&controllers.MatchController{},
			),
		),
		beego.NSNamespace("/auth",
			beego.NSInclude(
				&controllers.AuthController{},
			),
		),
	)
	beego.AddNamespace(ns)

	beego.InsertFilter("/v1/match", beego.BeforeExec, models.IsLoggedInFilter)
	beego.InsertFilter("/v1/match/*", beego.BeforeExec, models.IsLoggedInFilter)
}
