package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/robfig/revel"
	"strconv"
)

type User struct {
	*revel.Controller
}

func (user *User) Index(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	userid, _ := strconv.Atoi(id)
	articles, _ := manager.GetAllArticlesByUserId(userid)
	count := len(articles)

	var pageCount int
	if (count % models.ArticlesInSinglePage) == 0 {
		pageCount = count / models.ArticlesInSinglePage
	} else {
		pageCount = count/models.ArticlesInSinglePage + 1
	}

	pageSlice := make([]int, 0)
	for i := 1; i <= pageCount; i++ {
		pageSlice = append(pageSlice, i)
	}

	articlesOnOnePage := []models.UserArticle{}
	if count > models.ArticlesInSinglePage {
		articlesOnOnePage = articles[(count - models.ArticlesInSinglePage):]
	} else {
		articlesOnOnePage = articles
	}

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["allArticles"] = articles
	user.RenderArgs["articlesOnOnePage"] = articlesOnOnePage
	user.RenderArgs["pageCount"] = pageCount
	user.RenderArgs["pageSlice"] = pageSlice

	return user.Render()
}

func (user *User) Info(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	user.RenderArgs["authorId"] = id
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) Message(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) Friend(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) AddArticle() revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) PostAddArticle(article *models.UserArticle) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	article.AuthorId, err = strconv.Atoi(user.Session["userid"])
	if err != nil {
		fmt.Println("转换用户id失败")
		return user.Redirect((*User).AddArticle)
	}

	err = manager.AddUserArticle(article)
	if err != nil {
		user.Validation.Keep()
		user.FlashParams()
		user.Flash.Error(err.Error())
		return user.Redirect((*User).AddArticle)
	}

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Redirect((*User).AddArticle)
}

func (user *User) EditArticle(articleid string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) ShowArticle(userid int, articleid string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()
	article, _ := manager.GetArticleByUserIdAndArticleId(userid, articleid)
	authorId := strconv.Itoa(article.AuthorId)

	user.RenderArgs["authorId"] = authorId
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["article"] = article

	return user.Render()
}

func (user *User) PostEditArticle(article *models.UserArticle) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) EditInfo(userid int) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	userInfo, _ := manager.GetUserById(userid)

	user.RenderArgs["user"] = userInfo
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) PostEditInfo(userInfo models.User) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()
	userInfo.Id, _ = strconv.Atoi(user.Session["userid"])
	err = manager.UpdateUserInfo(userInfo.Id, userInfo)

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Redirect((*App).Index)
}

func (user *User) DeleteArticle(articleid string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}
