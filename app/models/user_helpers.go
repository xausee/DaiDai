package models

import (
	"labix.org/v2/mgo/bson"
)

func (manager *DbManager) AddUserArticle(article *UserArticle) error {
	uc := manager.session.DB(DbName).C(CommonUserCollection)

	// i, _ := uc.Find(bson.M{"Content": cu.Id}).Count()
	// if i != 0 {
	// 	return errors.New("此条慧语已经存在")
	// }

	article.Id = bson.NewObjectId().Hex()
	err := uc.Insert(article)
	var commonuser CommonUser
	commonuser.Article.Author = "author"
	commonuser.Article.Content = "aaaaaaaaaaaaaaaa"
	commonuser.Id = "1111111111"
	commonuser.Email = "asd@asd.com"
	uc.Insert(commonuser)
	uc.Insert("{ _id: 10, type: \"misc\", item: \"card\", qty: 15 }")

	return err
}
