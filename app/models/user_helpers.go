package models

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"time"
)

func (manager *DbManager) AddUserArticle(article *UserArticle) error {
	uc := manager.session.DB(DbName).C(UserCollection)

	var userBefore, userAfter User
	err := uc.Find(bson.M{"id": article.AuthorId}).One(&userBefore)
	err = uc.Find(bson.M{"id": article.AuthorId}).One(&userAfter)

	article.Id = bson.NewObjectId().Hex()
	article.CreateTime = time.Now().Format("2006-01-02 15:04:05")

	//userBefore = userAfter
	as := userAfter.Articles
	as = append(as, *article)
	userAfter.Articles = as

	fmt.Println(userAfter.Articles)
	err = uc.Update(userBefore, userAfter)

	return err
}

func (manager *DbManager) GetAllArticlesByUserId(userid int) (articles []UserArticle, err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	var userInfo User
	err = uc.Find(bson.M{"id": userid}).One(&userInfo)
	if err != nil {
		fmt.Println("查询用户信息失败")
	}
	articles = userInfo.Articles

	return articles, err
}
