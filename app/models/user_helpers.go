package models

import (
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
	"time"
)

func (manager *DbManager) GetUserById(userid int) (userInfo User, err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	err = uc.Find(bson.M{"id": userid}).One(&userInfo)
	if err != nil {
		fmt.Println("查询用户信息失败")
	}

	return userInfo, err
}

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

func (manager *DbManager) GetArticleByUserIdAndArticleId(userid int, articleid string) (article UserArticle, err error) {
	articles, _ := manager.GetAllArticlesByUserId(userid)

	for _, art := range articles {
		fmt.Println(art)
		if art.Id == articleid {
			article = art
			fmt.Println("找到指定的文章")
			return article, err
		} else {
			errors.New("查找文章失败")
			fmt.Println("找不到指定的文章")
		}
	}

	return article, err
}
