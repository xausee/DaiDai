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

func (user *User) Index(userid int) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

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

	// 判断访问该页面的用户是否是本人
	var isCurrentUser bool
	if user.Session["userid"] == strconv.Itoa(userid) {
		isCurrentUser = true
	}

	user.RenderArgs["isCurrentUser"] = isCurrentUser
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["allArticles"] = articles
	user.RenderArgs["articlesOnOnePage"] = articlesOnOnePage
	user.RenderArgs["pageCount"] = pageCount
	user.RenderArgs["pageSlice"] = pageSlice

	return user.Render()
}

func (user *User) Info(id int) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	userInfo, _ := manager.GetUserById(id)

	// 判断访问该页面的用户是否是本人
	var isCurrentUser bool
	if user.Session["userid"] == strconv.Itoa(id) {
		isCurrentUser = true
	}

	user.RenderArgs["isCurrentUser"] = isCurrentUser
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["userInfo"] = userInfo

	return user.Render()
}

func (user *User) Message(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	user.RenderArgs["userid"] = user.Session["userid"]
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
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) PostAddArticle(article *models.UserArticle) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	// 给文章附加上作者id和名字（昵称）
	article.AuthorId, err = strconv.Atoi(user.Session["userid"])
	if err != nil {
		fmt.Println("转换用户id失败")
		return user.Redirect((*User).AddArticle)
	}
	article.AuthorNickName = user.Session["nickName"]

	err = manager.AddUserArticle(article)
	if err != nil {
		user.Validation.Keep()
		user.FlashParams()
		user.Flash.Error(err.Error())
		return user.Redirect((*User).AddArticle)
	}

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Redirect("/user/article/show/%d/%s", article.AuthorId, article.Id)
}

func (user *User) EditArticle(articleid string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	authorid, err := strconv.Atoi(user.Session["userid"])
	// 根据作者ID和文章ID查找到该文章
	article, _ := manager.GetArticleByUserIdAndArticleId(authorid, articleid)

	user.RenderArgs["oldArticle"] = article
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) PostEditArticle(oldArticleId string, newArticle models.UserArticle) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	authorid, err := strconv.Atoi(user.Session["userid"])

	newArticle.Id = oldArticleId
	err = manager.UpdateUserArticle(authorid, newArticle)
	if err != nil {
		fmt.Println("更新文章失败")
	}

	//user.RenderArgs["oldArticle"] = article
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Redirect("/user/article/show/%d/%s", authorid, newArticle.Id)
}

func (user *User) ShowArticle(userid int, articleid string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()
	article, _ := manager.GetArticleByUserIdAndArticleId(userid, articleid)

	// 判断访问该页面的用户是否是本人
	var isCurrentUser bool
	if user.Session["userid"] == strconv.Itoa(article.AuthorId) {
		isCurrentUser = true
	}

	user.RenderArgs["isCurrentUser"] = isCurrentUser
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["article"] = article
	user.RenderArgs["commentCount"] = len(article.Comments)

	return user.Render()
}

func (user *User) PostArticleComment(authorid int, articleid string, comment *models.Comment) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	// 根据作者ID和文章ID查找到该文章
	article, _ := manager.GetArticleByUserIdAndArticleId(authorid, articleid)

	// 增加评论人的昵称
	if user.Session["nickName"] == "" {
		comment.Author.NickName = "匿名网友"
	} else {
		comment.Author.NickName = user.Session["nickName"]
	}

	// 添加新的评论
	err = manager.AddArticleComment(article, comment)

	user.RenderArgs["article"] = article
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Redirect("/user/article/show/%d/%s", authorid, articleid)
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
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Redirect("/user/info/%d", userInfo.Id)
}

func (user *User) DeleteArticle(articleid string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()
	// 能调用此方法的用户就是作者本人，所以authorid等于Session["userid"]
	authorid, _ := strconv.Atoi(user.Session["userid"])
	err = manager.DeleteUserArticle(authorid, articleid)

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Redirect("/user/%d", authorid)
}
