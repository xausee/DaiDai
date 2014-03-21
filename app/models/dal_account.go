package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
)

func (manager *DbManager) RegisterUser(mu *MockUser) error {
	uc := manager.session.DB(DbName).C(UserCollection)

	i, _ := uc.Find(bson.M{"nickname": mu.Nickname}).Count()
	if i != 0 {
		return errors.New("用户昵称已经被使用")
	}

	i, _ = uc.Find(bson.M{"email": mu.Email}).Count()
	if i != 0 {
		return errors.New("邮件地址已经被使用")
	}

	count, _ := uc.Count()

	var u User
	// Id start from 10000
	u.Id = 10000 + count
	u.Email = mu.Email
	u.Nickname = mu.Nickname
	u.Gender = mu.Gender
	u.Password, _ = bcrypt.GenerateFromPassword([]byte(mu.Password), bcrypt.DefaultCost)

	err := uc.Insert(u)

	return err
}

func (manager *DbManager) LoginUser(lu *LoginUser) (user *User, err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	i, _ := uc.Find(bson.M{"email": lu.Email}).Count()
	if i == 0 {
		fmt.Println("此账号不存在")
		err = errors.New("此账号不存在")
	}

	uc.Find(bson.M{"email": lu.Email}).One(&user)

	if user.Password == nil {
		err = errors.New("获取密码错误")
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(lu.Password))
	if err != nil {
		fmt.Println("密码不正确")
		err = errors.New("密码不正确")
	}
	return
}
