package main

import (
	_ "github.com/astaxie/beego/session/mysql"
	_ "github.com/levilovelock/magitrak/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func main() {
	dbAddress := beego.AppConfig.String("modelORMaddress")
	dbType := beego.AppConfig.String("modelORMdb")
	if dbAddress == "" {
		beego.Error("Cannot find config line for modelORMaddress - please set it!")
	}
	dbErr := orm.RegisterDataBase("default", dbType, dbAddress, 30)
	if dbErr != nil {
		beego.Error(dbErr)
	}

	beego.Run()
}
