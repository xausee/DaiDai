package models

import (
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
	"time"
)

func (manager *DbManager) AddUserArticle(article *UserArticle) error {
	uc := manager.session.DB(DbName).C(UserCollection)

	articles := []UserArticle{}
	//err := uc.Find(nil).All(&articles)

	err := uc.Find(nil).All(&articles)
	if err != nil {
		return errors.New("此条慧语已经存在")
	}

	fmt.Println(articles)

	article.Id = bson.NewObjectId().Hex()
	err = uc.Insert(article)
	var commonuser User
	commonuser.Article.AuthorId = 1
	commonuser.Article.Content = "aaaaaaaaaaaaaaaa"
	commonuser.Id = 1111111111
	commonuser.Email = "asd@asd.com"
	commonuser.Article.Comments.Author.Nickname = "Phiso"
	commonuser.Article.Comments.Author.Email = "Phiso"
	commonuser.Article.Comments.Time = time.Now()
	commonuser.Article.Comments.Score = 10
	uc.Insert(commonuser)

	return err
}
