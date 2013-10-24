package models


import (   
  "errors" 
  "labix.org/v2/mgo/bson" 
  "code.google.com/p/go.crypto/bcrypt"
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
  u.Password, _ = bcrypt.GenerateFromPassword([]byte(mu.Password), bcrypt.DefaultCost)

  err := uc.Insert(u)

  return err 
}