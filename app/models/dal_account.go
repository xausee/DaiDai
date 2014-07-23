package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
)

func (manager *DbManager) RegisterUser(mu *MockUser) error {
	uc := manager.session.DB(DbName).C(UserCollection)

	i, _ := uc.Find(bson.M{"nickname": mu.NickName}).Count()
	if i != 0 {
		return errors.New("用户昵称已经被使用")
	}

	// i, _ = uc.Find(bson.M{"email": mu.Email}).Count()
	// if i != 0 {
	// 	return errors.New("邮件地址已经被使用")
	// }

	count, _ := uc.Count()

	var u User
	// 用户ID从10000开始
	u.Id = 10000 + count
	u.Email = mu.Email
	u.NickName = mu.NickName
	u.PenName = mu.PenName
	u.Gender = mu.Gender
	u.Password, _ = bcrypt.GenerateFromPassword([]byte(mu.Password), bcrypt.DefaultCost)

	if u.Gender == "男" {
		u.AvatarUrl = DefaultBoyAvatarUrl
	} else {
		u.AvatarUrl = DefaultGirlAvatarUrl
	}

	err := uc.Insert(u)

	return err
}

func (manager *DbManager) LoginUser(lu *LoginUser) (user *User, err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	i, _ := uc.Find(bson.M{"nickname": lu.NickName}).Count()
	if i == 0 {
		fmt.Println("此账号不存在")
		err = errors.New("此账号不存在")
		return
	}

	uc.Find(bson.M{"nickname": lu.NickName}).One(&user)

	if user.Password == nil {
		err = errors.New("获取密码错误")
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(lu.Password))
	if err != nil {
		fmt.Println("密码不正确")
		err = errors.New("密码不正确")
	}
	return
}
