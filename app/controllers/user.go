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

func (user *User) Index(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	userInfo, e := manager.GetUserByNickName(nickName)
	if e != nil {
		return user.Redirect((*ErrorPages).Page404)
	}

	articles := userInfo.Articles
	articlesCount := len(articles)

	var pageCount int
	if (articlesCount % models.ArticlesInSinglePage) == 0 {
		pageCount = articlesCount / models.ArticlesInSinglePage
	} else {
		pageCount = articlesCount/models.ArticlesInSinglePage + 1
	}

	pageSlice := make([]int, 0)
	for i := 1; i <= pageCount; i++ {
		pageSlice = append(pageSlice, i)
	}

	articlesOnOnePage := []models.UserArticle{}
	if articlesCount > models.ArticlesInUserHomePanel {
		articlesOnOnePage = articles[(articlesCount - models.ArticlesInUserHomePanel):]
	} else {
		articlesOnOnePage = articles
	}

	// 判断访问该页面的用户是否是本人
	var isAuthor bool
	if user.Session["nickName"] == nickName {
		isAuthor = true
	}

	user.RenderArgs["isAuthor"] = isAuthor
	if len(userInfo.Message) != 0 {
		user.RenderArgs["lastMessage"] = userInfo.Message[len(userInfo.Message)-1]
	}

	isWatched := false
	for _, v := range userInfo.Fans {
		// 访问者是被访问者的粉丝，设置标记isWatched
		if v.NickName == user.Session["nickName"] {
			isWatched = true
		}
	}

	user.RenderArgs["ownerNickName"] = nickName
	user.RenderArgs["isWatched"] = isWatched
	user.RenderArgs["fansCount"] = len(userInfo.Fans)
	user.RenderArgs["wathCount"] = len(userInfo.Watch)
	user.RenderArgs["articleCollectionCount"] = len(userInfo.ArticleCollection)
	user.RenderArgs["messageCount"] = len(userInfo.Message)
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["allArticles"] = articles
	user.RenderArgs["articleCount"] = len(articles)
	user.RenderArgs["articlesOnOnePage"] = articlesOnOnePage
	user.RenderArgs["pageCount"] = pageCount
	user.RenderArgs["pageSlice"] = pageSlice

	return user.Render()
}

func (user *User) OriginalArticleList(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	articles, _ := manager.GetAllArticlesByUserNickName(nickName)
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
	var isAuthor bool
	if user.Session["nickName"] == nickName {
		isAuthor = true
	}

	user.RenderArgs["isAuthor"] = isAuthor
	user.RenderArgs["ownerNickName"] = nickName
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["allArticles"] = articles
	user.RenderArgs["articleCount"] = len(articles)
	user.RenderArgs["articlesOnOnePage"] = articlesOnOnePage
	user.RenderArgs["pageCount"] = pageCount
	user.RenderArgs["pageSlice"] = pageSlice

	return user.Render()
}

func (user *User) PageList(authorNickName string, pageNumber string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	articles, err := manager.GetAllArticlesByUserNickName(authorNickName)
	count := len(articles)

	var pageCount int
	if (count % models.ArticlesInSinglePage) == 0 {
		pageCount = count / models.ArticlesInSinglePage
	} else {
		pageCount = count/models.ArticlesInSinglePage + 1
	}

	var iPageNumber int
	iPageNumber, err = strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Println(err)
	}

	artilcesOnOnePage := []models.UserArticle{}
	if count >= iPageNumber*models.ArticlesInSinglePage {
		artilcesOnOnePage = articles[(iPageNumber-1)*models.ArticlesInSinglePage : iPageNumber*models.ArticlesInSinglePage]
	} else if (iPageNumber-1)*models.ArticlesInSinglePage < count && count < iPageNumber*models.ArticlesInSinglePage {
		artilcesOnOnePage = articles[(iPageNumber-1)*models.ArticlesInSinglePage:]
	} else if (iPageNumber-1)*models.ArticlesInSinglePage > count {
		fmt.Println("已超过最大页码")
	}
	fmt.Println("pageCount:", pageCount)
	fmt.Println("pageNumber:", pageNumber)

	// 判断访问该页面的用户是否是作者本人
	var isAuthor bool
	if user.Session["nickName"] == authorNickName {
		isAuthor = true
	}

	user.RenderArgs["isAuthor"] = isAuthor
	user.RenderArgs["allArticles"] = articles
	user.RenderArgs["artilcesOnOnePage"] = artilcesOnOnePage
	user.RenderArgs["pageCount"] = pageCount
	user.RenderArgs["pageNumber"] = pageNumber

	return user.Render()
}

func (user *User) Info(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	userInfo, _ := manager.GetUserByNickName(nickName)

	// 判断访问该页面的用户是否是本人
	var isAuthor bool
	if user.Session["nickName"] == nickName {
		isAuthor = true
	}

	user.RenderArgs["isAuthor"] = isAuthor
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["userInfo"] = userInfo

	return user.Render()
}

func (user *User) Message(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	messages, _ := manager.GetAllMessageByUserNickName(nickName)

	user.RenderArgs["ownerNickName"] = nickName
	user.RenderArgs["messages"] = messages
	user.RenderArgs["messageCount"] = len(messages)
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) Fans(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	fans, _ := manager.GetAllFansByUserNickName(nickName)

	user.RenderArgs["ownerNickName"] = nickName
	user.RenderArgs["fans"] = fans
	user.RenderArgs["fansCount"] = len(fans)
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) Watch(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	watches, _ := manager.GetAllWatchByUserNickName(nickName)

	user.RenderArgs["ownerNickName"] = nickName
	user.RenderArgs["watches"] = watches
	user.RenderArgs["watchesCount"] = len(watches)
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) AddWatch(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	// 添加到自己的关注
	watcherNickName := user.Session["nickName"]
	watchedUserNickName := nickName
	fmt.Println(watcherNickName, watchedUserNickName)
	err = manager.AddWatch(watcherNickName, watchedUserNickName)

	// 添加到对方的粉丝
	fansNickName := user.Session["nickName"]
	ownerNickName := nickName
	err = manager.AddFans(fansNickName, ownerNickName)

	user.RenderArgs["ownerNickName"] = nickName
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Redirect("/user/%s", nickName)
}

func (user *User) DeleteWatch(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	// 从访问者的关注里删除数据
	watcherNickName := user.Session["nickName"]
	watchedUserNickName := nickName
	err = manager.DeleteWatch(watcherNickName, watchedUserNickName)

	// 从访问者的粉丝里删除数据
	fansNickName := user.Session["nickName"]
	ownerNickName := nickName
	err = manager.DeleteFans(fansNickName, ownerNickName)

	user.RenderArgs["ownerNickName"] = nickName
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Redirect("/user/%s", nickName)
}

func (user *User) PostMessage(nickName string, message models.Comment) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	// 增加评论人的昵称
	if user.Session["nickName"] == "" {
		message.Author.NickName = "匿名网友"
	} else {
		message.Author.NickName = user.Session["nickName"]
	}

	err = manager.UpdateMessageByNickName(nickName, message)

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Redirect("/user/%s/message", nickName)
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

	// 使用的是article的指针， AddUserArticle里面给article ID赋值后
	// 可调用下面controller(ShowArticle) 访问新增的article
	return user.Redirect("/user/%s/article/%s", article.AuthorNickName, article.Id)
}

func (user *User) ArticleCollection(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	articlCollection, _ := manager.GetArticleCollectionByNickName(nickName)

	user.RenderArgs["ownerNickName"] = nickName
	user.RenderArgs["articlCollection"] = articlCollection
	user.RenderArgs["articlCollectionCount"] = len(articlCollection)
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) AddToArticleCollection(articleAuthorNickName string, aticleTitle string, articleId string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	userNickName := user.Session["nickName"]

	// 添加文章到收藏
	var article models.ArticleInCollection
	article.Id = articleId
	article.Title = aticleTitle
	article.AuthorNickName = articleAuthorNickName

	err = manager.AddAticleToArticleCollection(userNickName, article)

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Redirect("/user/%s/articlecollection", userNickName)
}

func (user *User) DeleteFromArticleCollection(articleAuthorNickName string, aticleTitle string, articleId string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	userNickName := user.Session["nickName"]

	// 从收藏里删除该文章
	var article models.ArticleInCollection
	article.Id = articleId
	article.Title = aticleTitle
	article.AuthorNickName = articleAuthorNickName

	err = manager.DeleteAticleFromArticleCollection(userNickName, article)

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Redirect("/user/%s/articlecollection", userNickName)
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

	authorNickName := user.Session["nickName"]

	newArticle.Id = oldArticleId
	err = manager.UpdateUserArticle(authorNickName, newArticle)
	if err != nil {
		fmt.Println("更新文章失败")
	}

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Redirect("/user/%s/article/%s", authorNickName, newArticle.Id)
}

func (user *User) ShowArticle(nickName string, articleid string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	article, _ := manager.GetArticleByUserNickNameAndArticleId(nickName, articleid)

	// 判断访问该页面的用户是否是本人
	var isCurrentUser bool
	if user.Session["userid"] == strconv.Itoa(article.AuthorId) {
		isCurrentUser = true
	}

	// 从访问者本人收藏里查找，确认是否已经收藏
	isInArticleCollection := false
	userNickName := user.Session["nickName"]
	articlCollection, _ := manager.GetArticleCollectionByNickName(userNickName)
	for _, v := range articlCollection {
		if v.Id == articleid {
			isInArticleCollection = true
			break
		}
	}

	user.RenderArgs["isInArticleCollection"] = isInArticleCollection
	user.RenderArgs["isCurrentUser"] = isCurrentUser
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["article"] = article
	user.RenderArgs["commentCount"] = len(article.Comments)

	return user.Render()
}

func (user *User) PostArticleComment(authorNickName string, articleid string, comment *models.Comment) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	// 根据作者ID和文章ID查找到该文章
	article, _ := manager.GetArticleByUserNickNameAndArticleId(authorNickName, articleid)

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

	return user.Redirect("/user/%s/article/%s", authorNickName, articleid)
}

func (user *User) EditInfo(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	userInfo, _ := manager.GetUserByNickName(nickName)

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
	err = manager.UpdateUserInfoById(userInfo.Id, userInfo)

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Redirect("/user/%s/info", user.Session["nickName"])
}

func (user *User) DeleteArticle(articleid string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()
	// 能调用此方法的用户就是作者本人，所以nickName等于Session["nickName"]
	err = manager.DeleteUserArticle(user.Session["nickName"], articleid)

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Redirect("/user/%s", user.Session["nickName"])
}
