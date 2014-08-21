package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/revel/revel"
	"os"
	"strconv"
)

type Admin struct {
	*revel.Controller
}

func (this *Admin) Index() revel.Result {
	userid := this.Session["userid"]
	nickName := this.Session["nickName"]
	avatarUrl := this.Session["avatarUrl"]

	return this.Render(userid, nickName, avatarUrl)
}

func (this *Admin) Register() revel.Result {
	return this.Render()
}

func (this *Admin) RegisterSuccess() revel.Result {
	return this.Render()
}

func (this *Admin) Login() revel.Result {
	return this.Render()
}

// 该action不实际使用，共用account的logout
func (this *Admin) Logout() revel.Result {
	fmt.Println("登出管理员账号: ", this.Session["nickName"])
	for k := range this.Session {
		delete(this.Session, k)
	}
	return this.Redirect((*App).Index)
}

func (this *Admin) PostRegister(user *models.MockUser) revel.Result {
	//this.Validation.Email(user.Email).Message("电子邮件格式无效")
	this.Validation.Required(user.NickName).Message("管理员账号名称不能为空")
	this.Validation.Required(user.Password).Message("密码不能为空")
	this.Validation.Required(user.ConfirmPassword).Message("确认密码不能为空")
	this.Validation.MinSize(user.Password, 6).Message("密码长度不短于6位")
	this.Validation.Required(user.ConfirmPassword == user.Password).Message("两次输入的密码不一致")

	if this.Validation.HasErrors() {
		this.Validation.Keep()
		this.FlashParams()
		return this.Redirect((*Admin).Register)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()

	err = manager.RegisterUser(user)
	if err != nil {
		this.Validation.Clear()

		// 添加错误信息，显示在页面的账号名下面
		var e revel.ValidationError
		e.Message = err.Error()
		e.Key = "user.NickName"
		this.Validation.Errors = append(this.Validation.Errors, &e)

		this.Validation.Keep()
		this.FlashParams()
		return this.Redirect((*Admin).Register)
	}

	return this.Redirect((*Admin).RegisterSuccess)
}

func (this *Admin) PostLogin(loginUser *models.LoginUser) revel.Result {
	this.Validation.Required(loginUser.NickName).Message("请输入管理员账号")
	this.Validation.Required(loginUser.Password).Message("请输入密码")
	this.Validation.Required(loginUser.NickName == "水蓝色").Message("管理员账号不正确")

	if this.Validation.HasErrors() {
		this.Validation.Keep()
		this.FlashParams()
		return this.Redirect((*Admin).Login)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()

	var u *models.User
	u, err = manager.LoginUser(loginUser)

	if err != nil {
		this.Validation.Clear()

		// 添加错误提示信息，显示在页面的用户名/密码下面
		var e revel.ValidationError
		if err.Error() == "该管理员账号不存在" {
			e.Key = "loginUser.NickName"
		} else {
			e.Key = "loginUser.Password"
		}
		e.Message = err.Error()
		this.Validation.Errors = append(this.Validation.Errors, &e)

		this.Validation.Keep()
		this.FlashParams()
		return this.Redirect((*Admin).Login)
	}

	this.Session["userid"] = strconv.Itoa(u.Id)
	this.Session["nickName"] = u.NickName

	userInfo, e := manager.GetUserByNickName(u.NickName)
	fmt.Println("头像地址: ", userInfo.AvatarUrl)
	if e != nil {
		return this.Redirect((*ErrorPages).Page404)
	}

	if userInfo.AvatarUrl == "" {
		if userInfo.Gender == "男" {
			this.Session["avatarUrl"] = models.DefaultBoyAvatarUrl
		} else {
			this.Session["avatarUrl"] = models.DefaultGirlAvatarUrl
		}
	} else {
		this.Session["avatarUrl"] = userInfo.AvatarUrl
	}

	fmt.Println("使用管理员账号登陆: ", loginUser.NickName)

	// 使用session保存管理员权限
	// session 里保存的只能是字符串，所以需要先将其转换字符串
	this.Session[models.CSessionRole] = strconv.Itoa(int(models.AdminRole))

	// 登陆成功，跳转到管理员首页
	return this.Redirect((*Admin).Index)
}

func (this *Admin) UploadPicture() revel.Result {
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Render()
}

func (this *Admin) PostUploadPicture(pictures []*os.File) revel.Result {
	err := models.BackupPicture()
	if err != nil {
		fmt.Println("备份图片失败")
	}

	models.SavePicture(pictures)
	if err != nil {
		fmt.Println("保存新图片失败")
	}

	return this.Redirect((*Admin).UploadPicture)
}

func (this *Admin) UserManage() revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()

	var users []models.User
	users, err = manager.GetAllUser()

	if err != nil {
		return this.RenderError(err)
	}

	this.RenderArgs["users"] = users
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Render()
}

func (this *Admin) SearchUserByNickName(nickName string) revel.Result {

	return this.Render()
}

func (this *Admin) PostSearchUserByNickName(keywords string) revel.Result {
	fmt.Println("搜索用户：", keywords)
	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()

	var users []models.User
	users, err = manager.SearchUser(keywords)

	if err != nil {
		return this.RenderError(err)
	}

	this.RenderArgs["users"] = users
	return this.Render()
}

func initAuthMap() map[string]models.Role {
	authMap := make(map[string]models.Role)
	authMap["account.login"] = models.AnonymousRole
	authMap["account.logout"] = models.UserRole

	authMap["App.Index"] = models.AnonymousRole
	authMap["App.AboutUs"] = models.AnonymousRole
	authMap["App.Donate"] = models.AnonymousRole
	authMap["App.NotAuthorized"] = models.AnonymousRole

	authMap["Admin.Index"] = models.AdminRole
	authMap["Admin.Register"] = models.AdminRole
	authMap["Admin.PostRegister"] = models.AdminRole
	authMap["Admin.RegisterSuccess"] = models.AdminRole
	authMap["Admin.Login"] = models.AnonymousRole
	authMap["Admin.PostLogin"] = models.AnonymousRole
	authMap["Admin.Logout"] = models.AdminRole
	authMap["Admin.UploadPicture"] = models.AdminRole
	authMap["Admin.PostUploadPicture"] = models.AdminRole
	authMap["Admin.UserManage"] = models.AdminRole
	authMap["Admin.SearchUserByNickName"] = models.AdminRole
	authMap["Admin.PostSearchUserByNickName"] = models.AdminRole

	authMap["ErrorPages.Page404"] = models.AnonymousRole
	authMap["ErrorPages.Page500"] = models.AnonymousRole

	authMap["Account.Register"] = models.AnonymousRole
	authMap["Account.RegisterSuccessful"] = models.AnonymousRole
	authMap["Account.Login"] = models.AnonymousRole
	authMap["Account.PostLogin"] = models.AnonymousRole
	authMap["Account.Logout"] = models.AnonymousRole

	authMap["User.Index"] = models.AnonymousRole
	authMap["User.Info"] = models.AnonymousRole
	authMap["User.GetBasicInfo"] = models.AnonymousRole
	authMap["User.SetPassword"] = models.UserRole
	authMap["User.PostSetPassword"] = models.UserRole
	authMap["User.SetProfile"] = models.UserRole
	authMap["User.PostSetProfile"] = models.UserRole
	authMap["User.GetExtraProfile"] = models.AnonymousRole
	authMap["User.AddArticle "] = models.UserRole
	authMap["User.PostAddArticle"] = models.UserRole
	authMap["User.OriginalArticleList"] = models.AnonymousRole
	authMap["User.ShowArticle"] = models.AnonymousRole
	authMap["User.PostArticleComment"] = models.AnonymousRole
	authMap["User.EditArticle"] = models.UserRole
	authMap["User.PostEditArticle"] = models.UserRole
	authMap["User.Message"] = models.AnonymousRole
	authMap["User.Fans"] = models.AnonymousRole
	authMap["User.Watch"] = models.AnonymousRole
	authMap["User.ArticleCollection"] = models.AnonymousRole
	authMap["User.PageList"] = models.AnonymousRole
	authMap["User.PageList"] = models.AnonymousRole
	authMap["User.Friend"] = models.AnonymousRole

	authMap["Quotation.Index"] = models.AnonymousRole
	authMap["Quotation.TypeIndex"] = models.AnonymousRole
	authMap["Quotation.Add"] = models.AdminRole
	authMap["Quotation.PostAdd"] = models.AdminRole
	authMap["Quotation.Edit"] = models.AdminRole
	authMap["Quotation.PostEdit"] = models.AdminRole
	authMap["Quotation.Show"] = models.AnonymousRole
	authMap["Quotation.PageList"] = models.AnonymousRole
	authMap["Quotation.PageListWithTag"] = models.AnonymousRole

	authMap["Witticism.Index"] = models.AnonymousRole
	authMap["Witticism.Add"] = models.AdminRole
	authMap["Witticism.PostAdd"] = models.AdminRole
	authMap["Witticism.Edit"] = models.AdminRole
	authMap["Witticism.PostEdit"] = models.AdminRole
	authMap["Witticism.Show"] = models.AnonymousRole

	authMap["AncientPoem.Add"] = models.AdminRole
	authMap["AncientPoem.PostAdd"] = models.AdminRole

	authMap["ModernPoem.Index"] = models.AnonymousRole
	authMap["ModernPoem.TypeIndex"] = models.AnonymousRole
	authMap["ModernPoem.Add"] = models.AdminRole
	authMap["ModernPoem.PostAdd"] = models.AdminRole
	authMap["ModernPoem.Edit"] = models.AdminRole
	authMap["ModernPoem.PostEdit"] = models.AdminRole
	authMap["ModernPoem.Show"] = models.AnonymousRole
	authMap["ModernPoem.PageList"] = models.AnonymousRole
	authMap["ModernPoem.PageListWithTag"] = models.AnonymousRole

	authMap["Essay.Index"] = models.AnonymousRole
	authMap["Essay.TypeIndex"] = models.AnonymousRole
	authMap["Essay.Add"] = models.AdminRole
	authMap["Essay.PostAdd"] = models.AdminRole
	authMap["Essay.Edit"] = models.AdminRole
	authMap["Essay.PostEdit"] = models.AdminRole
	authMap["Essay.Show"] = models.AnonymousRole
	authMap["Essay.PageList"] = models.AnonymousRole
	authMap["Essay.PageListWithTag"] = models.AnonymousRole

	authMap["HintFiction.Index"] = models.AnonymousRole
	authMap["HintFiction.TypeIndex"] = models.AnonymousRole
	authMap["HintFiction.Add"] = models.AdminRole
	authMap["HintFiction.PostAdd"] = models.AdminRole
	authMap["HintFiction.Edit"] = models.AdminRole
	authMap["HintFiction.PostEdit"] = models.AdminRole
	authMap["HintFiction.Show"] = models.AnonymousRole
	authMap["HintFiction.PageList"] = models.AnonymousRole

	return authMap
}

func checkAuthentication(c *revel.Controller) revel.Result {
	fmt.Println(c.Action)
	authMap := initAuthMap()
	//获取当前登陆用户的角色信息

	// session 里保存的只能是字符串，所以需要类型转换
	var userRole models.Role
	userRoleStr, isExists := c.Session[models.CSessionRole]

	if !isExists {
		userRole = models.AnonymousRole
	} else {
		value, _ := strconv.Atoi(userRoleStr)
		userRole = models.Role(value)
	}

	//获取相关action的权限定义
	if requiredRole, isExists := authMap[c.Action]; isExists {
		//判断权限，如果权限不够, 跳转到首页
		if userRole < requiredRole {
			return c.Redirect((*App).Index)
		}
	}

	//返回nil表示可以接着调用后面的Action，在这里就代表有权限访问
	return nil
}

func init() {
	//将checkAuthentication注册到revel的处理链中
	revel.InterceptFunc(checkAuthentication, revel.BEFORE, revel.ALL_CONTROLLERS)
}
