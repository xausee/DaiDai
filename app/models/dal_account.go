package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
)

func (d *Dal) RegisterUser(mu *MockUser) error {
	uc := d.session.DB(DbName).C(UserCollection)

	//先检查email和nickname是否已经被使用
	i, _ := uc.Find(bson.M{"nickname": mu.Nickname}).Count()
	if i != 0 {
		return errors.New("用户昵称已经被使用")
	}

	i, _ = uc.Find(bson.M{"email": mu.Email}).Count()
	if i != 0 {
		return errors.New("邮件地址已经被使用")
	}

	var u User
	u.Email = mu.Email
	u.Nickname = mu.Nickname
	u.Gender = mu.Gender
	u.Password, _ = bcrypt.GenerateFromPassword([]byte(mu.Password), bcrypt.DefaultCost)

	err := uc.Insert(u)

	return err
}

func (d *Dal) LoginUser(mu *MockUser) error {
	uc := d.session.DB(DbName).C(UserCollection)

	//先检查email和nickname是否已经被使用
	i, _ := uc.Find(bson.M{"nickname": mu.Nickname}).Count()
	if i == 0 {
		return errors.New("用户不存在")
	}

	i, _ = uc.Find(bson.M{"email": mu.Email}).Count()
	if i == 0 {
		return errors.New("邮件帐号不存在")
	}

	var user *User
	//var password []byte
	uc.Find(bson.M{"email": mu.Email}).One(&user)
	fmt.Println(user.Email)
	if mu.Email != user.Email {
		return errors.New("邮件不正确")
	}

	uc.Find(bson.M{"password": mu.Password}).One(&user)
	fmt.Println(user.Password)
	fmt.Println("user.Password")

	if user.Password == nil {
		return errors.New("获取密码错误")
	}

	var u User
	u.Email = mu.Email
	u.Nickname = mu.Nickname
	u.Gender = mu.Gender
	u.Password, _ = bcrypt.GenerateFromPassword([]byte(mu.Password), bcrypt.DefaultCost)

	// if u.Password != user.Password {
	// 	return errors.New("密码不正确")
	// }
	return nil
}
