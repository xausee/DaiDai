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
	article.CreateTime = time.Now()

	//userBefore = userAfter
	as := userAfter.Articles
	as = append(as, *article)
	userAfter.Articles = as

	fmt.Println(userAfter.Articles)
	err = uc.Update(userBefore, userAfter)

	return err
}
