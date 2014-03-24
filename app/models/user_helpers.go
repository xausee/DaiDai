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

	fmt.Println(userBefore)
	//fmt.Println(article)

	article.Id = bson.NewObjectId().Hex()
	article.CreateTime = time.Now()

	userBefore = userAfter
	fmt.Println("ddddddddddddddddddddddddddddddddddddddddddddddddd")

	userAfter.Article = *article
	//fmt.Println(userAfter)
	fmt.Println(userAfter.Article)
	err = uc.Update(userBefore, userAfter)
	err = uc.Insert(userAfter)

	return err
}
