package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id       int
	Email    string `orm:"size(150)`
	Password string `orm:"size(50)`
}

func init() {
	orm.RegisterModel(&User{})
}
