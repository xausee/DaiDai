package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/revel/revel"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

	// 倒序处理
	for i := 0; i < articlesCount/2; i++ {
		articles[i], articles[articlesCount-i-1] = articles[articlesCount-i-1], articles[i]
	}

	articlesOnOnePage := []models.UserArticle{}
	if articlesCount > models.ArticlesInUserHomePanel {
		articlesOnOnePage = articles[:models.ArticlesInUserHomePanel]
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
	user.RenderArgs["userAvatarUrl"] = userInfo.AvatarUrl
	user.RenderArgs["articleCollectionCount"] = len(userInfo.ArticleCollection)
	user.RenderArgs["messageCount"] = len(userInfo.Message)
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]
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

	// 倒序处理
	count := len(articles)
	for i := 0; i < count/2; i++ {
		articles[i], articles[count-i-1] = articles[count-i-1], articles[i]
	}

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
		// 获取前面最新的文章
		articlesOnOnePage = articles[:models.ArticlesInSinglePage]
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
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]
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

	// 倒序处理
	count := len(articles)
	for i := 0; i < count/2; i++ {
		articles[i], articles[count-i-1] = articles[count-i-1], articles[i]
	}

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

func (user *User) PageListOnHomePage(authorNickName string, pageNumber string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	articles, err := manager.GetAllArticlesByUserNickName(authorNickName)

	// 倒序处理
	count := len(articles)
	for i := 0; i < count/2; i++ {
		articles[i], articles[count-i-1] = articles[count-i-1], articles[i]
	}

	var pageCount int
	if (count % models.ArticlesInUserHomePanel) == 0 {
		pageCount = count / models.ArticlesInUserHomePanel
	} else {
		pageCount = count/models.ArticlesInUserHomePanel + 1
	}

	var iPageNumber int
	iPageNumber, err = strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Println(err)
	}

	artilcesOnOnePage := []models.UserArticle{}
	if count >= iPageNumber*models.ArticlesInUserHomePanel {
		artilcesOnOnePage = articles[(iPageNumber-1)*models.ArticlesInUserHomePanel : iPageNumber*models.ArticlesInUserHomePanel]
	} else if (iPageNumber-1)*models.ArticlesInUserHomePanel < count && count < iPageNumber*models.ArticlesInUserHomePanel {
		artilcesOnOnePage = articles[(iPageNumber-1)*models.ArticlesInUserHomePanel:]
	} else if (iPageNumber-1)*models.ArticlesInUserHomePanel > count {
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
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]
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

	// 判断访问该页面的用户是否是作者本人
	var isAuthor bool
	if user.Session["nickName"] == nickName {
		isAuthor = true
	}

	user.RenderArgs["isAuthor"] = isAuthor
	user.RenderArgs["ownerNickName"] = nickName
	user.RenderArgs["messages"] = messages
	user.RenderArgs["messageCount"] = len(messages)
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]

	return user.Render()
}

func (user *User) Fans(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	fans, _ := manager.GetAllFansByUserNickName(nickName)

	// 判断访问该页面的用户是否是作者本人
	var isAuthor bool
	if user.Session["nickName"] == nickName {
		isAuthor = true
	}

	user.RenderArgs["isAuthor"] = isAuthor
	user.RenderArgs["ownerNickName"] = nickName
	user.RenderArgs["fans"] = fans
	user.RenderArgs["fansCount"] = len(fans)
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]

	return user.Render()
}

func (user *User) Watch(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	watches, _ := manager.GetAllWatchByUserNickName(nickName)

	// 判断访问该页面的用户是否是本人
	var isAuthor bool
	if user.Session["nickName"] == nickName {
		isAuthor = true
	}

	user.RenderArgs["isAuthor"] = isAuthor
	user.RenderArgs["ownerNickName"] = nickName
	user.RenderArgs["watches"] = watches
	user.RenderArgs["watchesCount"] = len(watches)
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]

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
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]

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
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]

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
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]

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
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]

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
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]

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
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]

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

	// 判断访问该页面的用户是否是本人
	var isAuthor bool
	if user.Session["nickName"] == nickName {
		isAuthor = true
	}

	user.RenderArgs["isAuthor"] = isAuthor
	user.RenderArgs["ownerNickName"] = nickName
	user.RenderArgs["articlCollection"] = articlCollection
	user.RenderArgs["articlCollectionCount"] = len(articlCollection)
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]

	return user.Render()
}

func (user *User) AddToArticleCollection(articleAuthorNickName string, articleTitle string, articleId string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	userNickName := user.Session["nickName"]

	// 添加文章到收藏
	var article models.ArticleInCollection
	article.Id = articleId
	article.Title = articleTitle
	article.AuthorNickName = articleAuthorNickName

	err = manager.AddAticleToArticleCollection(userNickName, article)

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]

	return user.Redirect("/user/%s/article/%s", articleAuthorNickName, articleId)
}

func (user *User) DeleteFromArticleCollection(articleAuthorNickName string, articleTitle string, articleId string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	userNickName := user.Session["nickName"]

	// 从收藏里删除该文章
	var article models.ArticleInCollection
	article.Id = articleId
	article.Title = articleTitle
	article.AuthorNickName = articleAuthorNickName

	err = manager.DeleteAticleFromArticleCollection(userNickName, article)

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]

	return user.Redirect("/user/%s/article/%s", articleAuthorNickName, articleId)
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
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]

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
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]

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

	// 查询该文章是否已经被管理员推荐到首页
	isRecommonded := manager.IsArticleRecommend(articleid)

	user.RenderArgs["isRecommonded"] = isRecommonded
	user.RenderArgs["isInArticleCollection"] = isInArticleCollection
	user.RenderArgs["isCurrentUser"] = isCurrentUser
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]
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
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]

	return user.Redirect("/user/%s/article/%s", authorNickName, articleid)
}

func (user *User) EditInfo(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	// 如果不是本人，则不能访问此页面，跳到其主页面
	if nickName != user.Session["nickName"] {
		return user.Redirect("/user/%s", user.Session["userid"])
	}

	userInfo, _ := manager.GetUserByNickName(nickName)

	user.RenderArgs["user"] = userInfo
	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]

	return user.Render()
}

func (user *User) PostEditInfo(uploadFile *os.File, userInfo models.User) revel.Result {
	// 使用revel requst formfile获取文件数据
	file, handler, err := user.Request.FormFile("uploadFile")
	if err != nil {
		fmt.Println(err)
	}
	// 读取所有数据
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	// 获取当前路径
	dir, patherr := filepath.Abs(filepath.Dir(os.Args[0]))
	if patherr != nil {
		log.Fatal(patherr)
	}

	// 文件路径
	filePath := dir + "/" + handler.Filename

	// 保存到文件
	err = ioutil.WriteFile(filePath, data, 0777)
	if err != nil {
		fmt.Println(err)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()
	userInfo.Id, _ = strconv.Atoi(user.Session["userid"])

	// 上传到七牛云储存
	key, e := models.UploadToQiniu(filePath)
	if e != nil {
		fmt.Println("修改头像失败:", e)
	} else {
		// 删除原来的头像文件
		oldUserInfo, _ := manager.GetUserById(userInfo.Id)
		avatarUrl := oldUserInfo.AvatarUrl
		urlArray := strings.SplitN(avatarUrl, "/", -1)
		oldKey := urlArray[len(urlArray)-1]
		derr := models.DeleteFileOnQiNiu(oldKey)
		if derr != nil {
			fmt.Println("删除原始头像文件失败:", derr)
		}
		// 保存新的头像地址
		userInfo.AvatarUrl = models.QiNiuSpace + key
		fmt.Println("头像地址：", userInfo.AvatarUrl)
	}

	// 删除临时头像文件
	err = os.Remove(filePath)
	if err != nil {
		fmt.Println("删除临时头像文件失败:", err)
	}

	// TODO： 直接使用 uploadFile 进行操作，放弃上面的方案

	err = manager.UpdateUserInfoById(userInfo.Id, userInfo)

	user.RenderArgs["userid"] = user.Session["userid"]
	user.RenderArgs["nickName"] = user.Session["nickName"]
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]

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
	user.RenderArgs["avatarUrl"] = user.Session["avatarUrl"]

	return user.Redirect("/user/%s", user.Session["nickName"])
}
