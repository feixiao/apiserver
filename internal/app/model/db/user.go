package db

import (
	"apiserver/internal/app/model"
	"apiserver/pkg/auth"
	"apiserver/pkg/constvar"
	"fmt"

	validator "gopkg.in/go-playground/validator.v9"
)

type User struct {
	Id       int    `json:"id" gorm:"column:id"`
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
	Nickname string `json:"nickname" gorm:"column:nickname"`
	Address  string `json:"address"  gorm:"column:address"`
}

func (user *User) TableName() string {
	return "user"
}

// User represents a registered user.
type UserModel struct {
	BaseModel
	Username string `json:"username" gorm:"column:username;not null" binding:"required" validate:"min=1,max=32"`
	Password string `json:"password" gorm:"column:password;not null" binding:"required" validate:"min=5,max=128"`
}

func (c *UserModel) TableName() string {
	return "tb_users"
}

// Create creates a new user account.
func (u *UserModel) Create() error {
	return model.MysqlHandler.Create(&u).Error
}

// DeleteUser deletes the user by the user identifier.
func DeleteUser(id uint64) error {
	user := UserModel{}
	user.BaseModel.Id = id
	return model.MysqlHandler.Delete(&user).Error
}

// Update updates an user account information.
func (u *UserModel) Update() error {
	return model.MysqlHandler.Save(u).Error
}

// GetUser gets an user by the user identifier.
func GetUser(username string) (*UserModel, error) {
	u := &UserModel{}
	d := model.MysqlHandler.Where("username = ?", username).First(&u)
	return u, d.Error
}

// ListUser List all users
func ListUser(username string, offset, limit int) ([]*UserModel, uint64, error) {
	if limit == 0 {
		limit = constvar.DefaultLimit
	}

	users := make([]*UserModel, 0)
	var count uint64

	where := fmt.Sprintf("username like '%%%s%%'", username)
	if err := model.MysqlHandler.Model(&UserModel{}).Where(where).Count(&count).Error; err != nil {
		return users, count, err
	}

	if err := model.MysqlHandler.Where(where).Offset(offset).Limit(limit).Order("id desc").Find(&users).Error; err != nil {
		return users, count, err
	}

	return users, count, nil
}

// Compare with the plain text password. Returns true if it's the same as the encrypted one (in the `User` struct).
func (u *UserModel) Compare(pwd string) (err error) {
	err = auth.Compare(u.Password, pwd)
	return
}

// Encrypt the user password.
func (u *UserModel) Encrypt() (err error) {
	u.Password, err = auth.Encrypt(u.Password)
	return
}

// Validate the fields.
func (u *UserModel) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
