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

func (this *User) Index(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	userInfo, e := manager.GetUserByNickName(nickName)
	if e != nil {
		return this.Redirect((*ErrorPages).Page404)
	}

	articles := userInfo.Articles
	articlesCount := len(articles)

	var pageCount int
	if (articlesCount % models.ArticlesInUserHomePanel) == 0 {
		pageCount = articlesCount / models.ArticlesInUserHomePanel
	} else {
		pageCount = articlesCount/models.ArticlesInUserHomePanel + 1
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
	if this.Session["nickName"] == nickName {
		isAuthor = true
	}

	this.RenderArgs["isAuthor"] = isAuthor
	if len(userInfo.Message) != 0 {
		this.RenderArgs["lastMessage"] = userInfo.Message[len(userInfo.Message)-1]
	} else {
		c := new(models.Comment)
		this.RenderArgs["lastMessage"] = c
	}

	isWatched := false
	for _, v := range userInfo.Fans {
		// 访问者是被访问者的粉丝，设置标记isWatched
		if v.NickName == this.Session["nickName"] {
			isWatched = true
		}
	}

	this.RenderArgs["ownerNickName"] = nickName
	this.RenderArgs["isWatched"] = isWatched
	this.RenderArgs["fansCount"] = len(userInfo.Fans)
	this.RenderArgs["wathCount"] = len(userInfo.Watch)
	this.RenderArgs["userAvatarUrl"] = userInfo.AvatarUrl
	this.RenderArgs["articleCollectionCount"] = len(userInfo.ArticleCollection)
	this.RenderArgs["messageCount"] = len(userInfo.Message)
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["allArticles"] = articles
	this.RenderArgs["articleCount"] = len(articles)
	this.RenderArgs["articlesOnOnePage"] = articlesOnOnePage
	this.RenderArgs["pageCount"] = pageCount
	this.RenderArgs["pageSlice"] = pageSlice

	return this.Render()
}

func (this *User) OriginalArticleList(nickName string) revel.Result {
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
	if this.Session["nickName"] == nickName {
		isAuthor = true
	}

	this.RenderArgs["isAuthor"] = isAuthor
	this.RenderArgs["ownerNickName"] = nickName
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["allArticles"] = articles
	this.RenderArgs["articleCount"] = len(articles)
	this.RenderArgs["articlesOnOnePage"] = articlesOnOnePage
	this.RenderArgs["pageCount"] = pageCount
	this.RenderArgs["pageSlice"] = pageSlice

	return this.Render()
}

func (this *User) PageList(authorNickName string, pageNumber string) revel.Result {
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
	if this.Session["nickName"] == authorNickName {
		isAuthor = true
	}

	this.RenderArgs["isAuthor"] = isAuthor
	this.RenderArgs["allArticles"] = articles
	this.RenderArgs["artilcesOnOnePage"] = artilcesOnOnePage
	this.RenderArgs["pageCount"] = pageCount
	this.RenderArgs["pageNumber"] = pageNumber

	return this.Render()
}

func (this *User) PageListOnHomePage(authorNickName string, pageNumber string) revel.Result {
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
	if this.Session["nickName"] == authorNickName {
		isAuthor = true
	}

	this.RenderArgs["isAuthor"] = isAuthor
	this.RenderArgs["allArticles"] = articles
	this.RenderArgs["artilcesOnOnePage"] = artilcesOnOnePage
	this.RenderArgs["pageCount"] = pageCount
	this.RenderArgs["pageNumber"] = pageNumber

	return this.Render()
}

func (this *User) Info(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	userInfo, _ := manager.GetUserByNickName(nickName)

	// 判断访问该页面的用户是否是本人
	var isAuthor bool
	if this.Session["nickName"] == nickName {
		isAuthor = true
	}

	this.RenderArgs["isAuthor"] = isAuthor
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["userInfo"] = userInfo

	return this.Render()
}

// 缺省显示，加载的数据跟controller Info 一样
func (this *User) GetBasicInfo(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	userInfo, _ := manager.GetUserByNickName(nickName)

	// 判断访问该页面的用户是否是本人
	var isAuthor bool
	if this.Session["nickName"] == nickName {
		isAuthor = true
	}

	this.RenderArgs["isAuthor"] = isAuthor
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["userInfo"] = userInfo

	return this.Render()
}

func (this *User) SetPassword(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	// 游客身份转到没授权页面
	if this.Session["nickName"] == "" {
		return this.Redirect((*App).NotAuthorized)
	}

	// 已注册用户，但不是本人，跳到其主页面
	if nickName != this.Session["nickName"] {
		return this.Redirect("/user/%s", this.Session["nickName"])
	}

	userInfo, _ := manager.GetUserByNickName(nickName)

	this.RenderArgs["user"] = userInfo
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Render()
}

func (this *User) UpdateSuccess() revel.Result {
	return this.Render()
}

func (this *User) PostSetPassword(oringinalPassword string, userInfo models.MockUser) revel.Result {
	this.Validation.Required(oringinalPassword).Message("原始密码不能为空")
	this.Validation.Required(userInfo.Password).Message("新密码不能为空")
	this.Validation.Required(userInfo.ConfirmPassword).Message("确认密码不能为空")
	this.Validation.MinSize(userInfo.Password, 6).Message("密码长度不短于6位")
	this.Validation.Required(userInfo.ConfirmPassword == userInfo.Password).Message("两次输入的密码不一致")
	//this.Validation.Required(userInfo.ConfirmPassword != oringinalPassword).Message("新密码不能跟原始密码相同")

	nickName := this.Session["nickName"]
	if this.Validation.HasErrors() {
		this.Validation.Keep()
		this.FlashParams()
		return this.Redirect("/user/profile/password?nickName=%s", nickName)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()

	var u *models.User
	u, err = manager.VerifyPasswordByNickName(nickName, oringinalPassword)
	if err != nil {
		this.Validation.Clear()

		// 添加错误信息，显示在页面的原始密码下面
		var e revel.ValidationError
		e.Message = err.Error()
		e.Key = "oringinalPassword"
		this.Validation.Errors = append(this.Validation.Errors, &e)

		this.Validation.Keep()
		this.FlashParams()
		return this.Redirect("/user/profile/password?nickName=%s", nickName)
	}

	err = manager.UpdateUserPasswordById(u.Id, userInfo.Password)
	if err != nil {
		this.Validation.Clear()

		// 添加错误信息，显示在表单下面
		var e revel.ValidationError
		e.Message = err.Error()
		e.Key = "userInfo.ConfirmPassword"
		this.Validation.Errors = append(this.Validation.Errors, &e)

		this.Validation.Keep()
		this.FlashParams()
		return this.Redirect("/user/profile/password?nickName=%s", nickName)
	}

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Redirect((*User).UpdateSuccess)
}

func (this *User) SetProfile(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	// 游客身份转到没授权页面
	if this.Session["nickName"] == "" {
		return this.Redirect((*App).NotAuthorized)
	}

	// 已注册用户，但不是本人，跳到其主页面
	if nickName != this.Session["nickName"] {
		return this.Redirect("/user/%s", this.Session["nickName"])
	}

	userInfo, _ := manager.GetUserByNickName(nickName)

	this.RenderArgs["user"] = userInfo
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Render()
}

func (this *User) PostSetProfile(uploadFile *os.File, userInfo models.User) revel.Result {
	// 使用revel requst formfile获取文件数据
	file, handler, err := this.Request.FormFile("uploadFile")
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
	userInfo.Id, _ = strconv.Atoi(this.Session["userid"])

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

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Redirect("/user/%s/info", this.Session["nickName"])
}

func (this *User) GetExtraProfile(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	userInfo, _ := manager.GetUserByNickName(nickName)

	// 判断访问该页面的用户是否是本人
	var isAuthor bool
	if this.Session["nickName"] == nickName {
		isAuthor = true
	}

	this.RenderArgs["isAuthor"] = isAuthor
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["userInfo"] = userInfo

	return this.Render()
}

func (this *User) Message(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	messages, _ := manager.GetAllMessageByUserNickName(nickName)

	// 判断访问该页面的用户是否是作者本人
	var isAuthor bool
	if this.Session["nickName"] == nickName {
		isAuthor = true
	}

	this.RenderArgs["isAuthor"] = isAuthor
	this.RenderArgs["ownerNickName"] = nickName
	this.RenderArgs["messages"] = messages
	this.RenderArgs["messageCount"] = len(messages)
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Render()
}

func (this *User) Fans(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	fans, _ := manager.GetAllFansByUserNickName(nickName)

	// 判断访问该页面的用户是否是作者本人
	var isAuthor bool
	if this.Session["nickName"] == nickName {
		isAuthor = true
	}

	this.RenderArgs["isAuthor"] = isAuthor
	this.RenderArgs["ownerNickName"] = nickName
	this.RenderArgs["fans"] = fans
	this.RenderArgs["fansCount"] = len(fans)
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Render()
}

func (this *User) Watch(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	watches, _ := manager.GetAllWatchByUserNickName(nickName)

	// 判断访问该页面的用户是否是本人
	var isAuthor bool
	if this.Session["nickName"] == nickName {
		isAuthor = true
	}

	this.RenderArgs["isAuthor"] = isAuthor
	this.RenderArgs["ownerNickName"] = nickName
	this.RenderArgs["watches"] = watches
	this.RenderArgs["watchesCount"] = len(watches)
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Render()
}

func (this *User) AddWatch(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	// 添加到自己的关注
	watcherNickName := this.Session["nickName"]
	watchedUserNickName := nickName
	fmt.Println(watcherNickName, watchedUserNickName)
	err = manager.AddWatch(watcherNickName, watchedUserNickName)

	// 添加到对方的粉丝
	fansNickName := this.Session["nickName"]
	ownerNickName := nickName
	err = manager.AddFans(fansNickName, ownerNickName)

	this.RenderArgs["ownerNickName"] = nickName
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Redirect("/user/%s", nickName)
}

func (this *User) DeleteWatch(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	// 从访问者的关注里删除数据
	watcherNickName := this.Session["nickName"]
	watchedUserNickName := nickName
	err = manager.DeleteWatch(watcherNickName, watchedUserNickName)

	// 从访问者的粉丝里删除数据
	fansNickName := this.Session["nickName"]
	ownerNickName := nickName
	err = manager.DeleteFans(fansNickName, ownerNickName)

	this.RenderArgs["ownerNickName"] = nickName
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Redirect("/user/%s", nickName)
}

func (this *User) PostMessage(nickName string, message models.Comment) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	// 增加评论人的昵称
	if this.Session["nickName"] == "" {
		message.Author.NickName = "匿名网友"
	} else {
		message.Author.NickName = this.Session["nickName"]
	}

	// 增加头像地址
	message.Author.AvatarUrl = this.Session["avatarUrl"]

	err = manager.UpdateMessageByNickName(nickName, message)

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Redirect("/user/%s/message", nickName)
}

func (this *User) Friend(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Render()
}

func (this *User) AddArticle() revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Render()
}

func (this *User) PostAddArticle(article *models.UserArticle) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	// 给文章附加上作者id和名字（昵称）
	article.AuthorId, err = strconv.Atoi(this.Session["userid"])
	if err != nil {
		fmt.Println("转换用户id失败")
		return this.Redirect((*User).AddArticle)
	}
	article.AuthorNickName = this.Session["nickName"]

	err = manager.AddUserArticle(article)
	if err != nil {
		this.Validation.Keep()
		this.FlashParams()
		this.Flash.Error(err.Error())
		return this.Redirect((*User).AddArticle)
	}

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	// 使用的是article的指针， AddUserArticle里面给article ID赋值后
	// 可调用下面controller(ShowArticle) 访问新增的article
	return this.Redirect("/user/%s/article/%s", article.AuthorNickName, article.Id)
}

func (this *User) ArticleCollection(nickName string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	articlCollection, _ := manager.GetArticleCollectionByNickName(nickName)

	// 判断访问该页面的用户是否是本人
	var isAuthor bool
	if this.Session["nickName"] == nickName {
		isAuthor = true
	}

	this.RenderArgs["isAuthor"] = isAuthor
	this.RenderArgs["ownerNickName"] = nickName
	this.RenderArgs["articlCollection"] = articlCollection
	this.RenderArgs["articlCollectionCount"] = len(articlCollection)
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Render()
}

func (this *User) AddToArticleCollection(articleAuthorNickName string, articleTitle string, articleId string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	userNickName := this.Session["nickName"]

	// 添加文章到收藏
	var article models.ArticleInCollection
	article.Id = articleId
	article.Title = articleTitle
	article.AuthorNickName = articleAuthorNickName

	err = manager.AddAticleToArticleCollection(userNickName, article)

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Redirect("/user/%s/article/%s", articleAuthorNickName, articleId)
}

func (this *User) DeleteFromArticleCollection(articleAuthorNickName string, articleTitle string, articleId string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	userNickName := this.Session["nickName"]

	// 从收藏里删除该文章
	var article models.ArticleInCollection
	article.Id = articleId
	article.Title = articleTitle
	article.AuthorNickName = articleAuthorNickName

	err = manager.DeleteAticleFromArticleCollection(userNickName, article)

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Redirect("/user/%s/article/%s", articleAuthorNickName, articleId)
}

func (this *User) EditArticle(articleid string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	authorid, err := strconv.Atoi(this.Session["userid"])
	// 根据作者ID和文章ID查找到该文章
	article, _ := manager.GetArticleByUserIdAndArticleId(authorid, articleid)

	this.RenderArgs["oldArticle"] = article
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Render()
}

func (this *User) PostEditArticle(oldArticleId string, newArticle models.UserArticle) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	authorNickName := this.Session["nickName"]

	newArticle.Id = oldArticleId
	err = manager.UpdateUserArticle(authorNickName, newArticle)
	if err != nil {
		fmt.Println("更新文章失败")
	}

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Redirect("/user/%s/article/%s", authorNickName, newArticle.Id)
}

func (this *User) ShowArticle(nickName string, articleid string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	article, _ := manager.GetArticleByUserNickNameAndArticleId(nickName, articleid)

	// 判断访问该页面的用户是否是本人
	var isCurrentUser bool
	if this.Session["userid"] == strconv.Itoa(article.AuthorId) {
		isCurrentUser = true
	}

	// 从访问者本人收藏里查找，确认是否已经收藏
	isInArticleCollection := false
	userNickName := this.Session["nickName"]
	articlCollection, _ := manager.GetArticleCollectionByNickName(userNickName)
	for _, v := range articlCollection {
		if v.Id == articleid {
			isInArticleCollection = true
			break
		}
	}

	// session 里保存的是字符串，所以需要类型转换
	var admin = false
	var role models.Role
	roleStr, isExists := this.Session[models.CSessionRole]

	if isExists {
		value, _ := strconv.Atoi(roleStr)
		role = models.Role(value)

		if role == models.AdminRole {
			admin = true
		}
	}

	// 查询该文章是否已经被管理员推荐到首页
	isRecommonded := manager.IsArticleRecommend(articleid)

	this.RenderArgs["admin"] = admin
	this.RenderArgs["isRecommonded"] = isRecommonded
	this.RenderArgs["isInArticleCollection"] = isInArticleCollection
	this.RenderArgs["isCurrentUser"] = isCurrentUser
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["article"] = article
	this.RenderArgs["commentCount"] = len(article.Comments)

	return this.Render()
}

func (this *User) PostArticleComment(authorNickName string, articleid string, comment *models.Comment) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	// 根据作者ID和文章ID查找到该文章
	article, _ := manager.GetArticleByUserNickNameAndArticleId(authorNickName, articleid)

	// 增加评论人的昵称
	if this.Session["nickName"] == "" {
		comment.Author.NickName = "匿名网友"
	} else {
		comment.Author.NickName = this.Session["nickName"]
	}

	// 增加头像地址
	comment.Author.AvatarUrl = this.Session["avatarUrl"]

	// 添加新的评论
	err = manager.AddArticleComment(article, comment)

	this.RenderArgs["article"] = article
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Redirect("/user/%s/article/%s", authorNickName, articleid)
}

func (this *User) DeleteArticle(articleid string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()
	// 能调用此方法的用户就是作者本人，所以nickName等于Session["nickName"]
	err = manager.DeleteUserArticle(this.Session["nickName"], articleid)

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Redirect("/user/%s", this.Session["nickName"])
}
