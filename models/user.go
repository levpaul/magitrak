package models

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int
	Email    string `orm:"size(150);unique"`
	Password string `orm:"size(60)" json:"-"`
}

const (
	MIN_PASSWORD_LEN = 8
	MAX_PASSWORD_LEN = 80
)

func init() {
	orm.RegisterModel(&User{})
}

func AddNewUser(email, password string) (*User, error) {
	// Validate password
	if len(password) < MIN_PASSWORD_LEN || len(password) > MAX_PASSWORD_LEN {
		return nil, errors.New("Password does not meet length requirements")
	}

	// Validate email
	if !govalidator.IsEmail(email) || len(email) > 150 {
		return nil, errors.New("Email does not meet email address requirements")
	}

	// Bcrypt password
	encryptedPass, bcryptErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if bcryptErr != nil {
		return nil, bcryptErr
	}

	beego.Debug("Length of bcryupt:", len(encryptedPass))

	// Create user and save to DB
	newUser := User{Email: email, Password: string(encryptedPass)}
	o := orm.NewOrm()
	_, dbErr := o.Insert(&newUser)
	if dbErr != nil {
		return nil, dbErr
	}

	// Everything went well
	return &newUser, nil
}
