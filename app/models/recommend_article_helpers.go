package models

import (
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
)

func (manager *DbManager) UpdateRecommendArticle(article ArticleInCollection) error {
	uc := manager.session.DB(DbName).C(RecommendArticleCollection)

	i, _ := uc.Find(bson.M{"Content": article.Title}).Count()

	j, _ := uc.Find(bson.M{"AuthorNickName": article.AuthorNickName}).Count()

	// 文章标题和作者都一致可认为已经被推荐过
	if i != 0 && j != 0 {
		return errors.New("此篇文章已被推荐")
	}

	articles, _ := manager.GetAllRecommendArticle()
	count := len(articles)

	err := errors.New("")
	for _, v := range articles {
		err = uc.Remove(v)
	}

	// 如果已经超过15篇，则删掉最早的一篇
	if count == ArticlesInHomePanel {
		articles = articles[1:]
	}

	articles = append(articles, article)

	for _, v := range articles {
		err = uc.Insert(v)
	}

	return err
}

func (manager *DbManager) DeleteRecommendArticle(article ArticleInCollection) error {
	uc := manager.session.DB(DbName).C(RecommendArticleCollection)

	err := uc.Remove(bson.M{"id": article.Id})

	return err
}

func (manager *DbManager) GetAllRecommendArticle() ([]ArticleInCollection, error) {
	uc := manager.session.DB(DbName).C(RecommendArticleCollection)

	count, err := uc.Count()
	fmt.Println("共有推荐的文章： ", count, "篇")
	allArticle := []ArticleInCollection{}
	err = uc.Find(nil).All(&allArticle)

	return allArticle, err
}

func (manager *DbManager) GetRecommendArticleById(articleId string) (ArticleInCollection, error) {
	uc := manager.session.DB(DbName).C(RecommendArticleCollection)

	var article ArticleInCollection
	err := uc.Find(bson.M{"id": articleId}).One(&article)
	if err != nil {
		fmt.Println("推荐的文章中不包含该篇文章")
	}

	return article, err
}

func (manager *DbManager) IsArticleRecommend(articleId string) (isRecommonded bool) {
	uc := manager.session.DB(DbName).C(RecommendArticleCollection)

	count, _ := uc.Find(bson.M{"id": articleId}).Count()
	if count != 0 {
		isRecommonded = true
		fmt.Println("推荐的文章包含该篇文章")
	} else {
		fmt.Println("推荐的文章不包含该篇文章")
		isRecommonded = false
	}

	return isRecommonded
}
