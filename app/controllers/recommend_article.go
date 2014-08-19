package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/revel/revel"
)

type RecommondArticle struct {
	*revel.Controller
}

func (this *RecommondArticle) Recommond(articleAuthorNickName string, articleTitle string, articleId string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	var article models.ArticleInCollection
	article.Id = articleId
	article.Title = articleTitle
	article.AuthorNickName = articleAuthorNickName

	err = manager.UpdateRecommendArticle(article)
	if err != nil {
		fmt.Println("推荐文章失败")
	}

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Redirect("/user/%s/article/%s", articleAuthorNickName, articleId)
}

func (this *RecommondArticle) DeleteRecommond(articleAuthorNickName string, articleTitle string, articleId string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	var article models.ArticleInCollection
	article.Id = articleId
	article.Title = articleTitle
	article.AuthorNickName = articleAuthorNickName

	err = manager.DeleteRecommendArticle(article)
	if err != nil {
		fmt.Println("推荐文章失败")
	}

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Redirect("/user/%s/article/%s", articleAuthorNickName, articleId)
}
