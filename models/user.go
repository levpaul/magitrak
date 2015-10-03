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
	Password string `orm:"size(60)"`
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

// Authenticates user by Email and on success reutrns the userID, on failure returns an error
func (u *User) Authenticate() (int, error) {
	if u.Email == "" || u.Password == "" {
		return 0, errors.New("Authentication failure, either email or password was empty")
	}
	suppliedPassword := u.Password

	o := orm.NewOrm()
	dbFindErr := o.Read(u, "Email")
	if dbFindErr != nil {
		return 0, errors.New("Failed to find User from DB with supplied Email")
	}

	passwordErr := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(suppliedPassword))
	if passwordErr != nil {
		return 0, errors.New("Failed to authenticate User with given email")
	}

	return u.Id, nil
}

func GetUserByEmail(email string) (*User, error) {
	u := &User{Email: email}
	o := orm.NewOrm()
	dbFindErr := o.Read(u, "Email")
	if dbFindErr != nil {
		return nil, dbFindErr
	}
	return u, nil
}

func (u *User) Delete() error {
	o := orm.NewOrm()
	beego.Debug("ASDFASDFASDFASDFSADFSADF:", u)
	_, err := o.Delete(u)
	return err
}
