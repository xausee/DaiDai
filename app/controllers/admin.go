package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/revel/revel"
	"os"
	"strconv"
	"strings"
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

func (this *Admin) Logout() revel.Result {
	fmt.Println("登出用户昵称: ", this.Session["nickName"])
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
	this.Validation.Required(loginUser.NickName).Message("请输入昵称")
	this.Validation.Required(loginUser.Password).Message("请输入密码")

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
		return this.Redirect((*Account).Login)
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
	this.Session[models.CSessionRole] = models.AdminRole

	// 登陆成功，跳转到管理员首页
	return this.Redirect((*Admin).Index)
}

func (this *Admin) UploadPicture() revel.Result {
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

func initAuthMap() {
	authMap := make(map[string]string)
	authMap["account.login"] = models.AnonymousRole
	authMap["account.logout"] = models.UserRole
	authMap["admin.index"] = models.AdminRole
}

func checkAuthentication(c *revel.Controller) revel.Result {
	authMap := make(map[string]string)
	authMap["account.login"] = models.AnonymousRole
	authMap["account.logout"] = models.UserRole
	authMap["admin.index"] = models.AdminRole
	//获取当前登陆用户的角色信息
	userRole, isExists := c.Session[models.CSessionRole]
	if !isExists {
		userRole = models.AnonymousRole
	}

	//获取紧接着要调用的Action的名称
	action := strings.ToLower(c.Action)
	//获取相关action的权限定义
	if requiredRole, isExists := authMap[action]; isExists {
		//判断权限，如果权限要求不相符
		if requiredRole != userRole {
			//跳转到首页
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
