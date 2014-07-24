package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/revel/revel"
	"strconv"
)

type Account struct {
	*revel.Controller
}

func (this *Account) Register() revel.Result {
	return this.Render()
}

func (this *Account) RegisterSuccessful() revel.Result {
	return this.Render()
}

func (this *Account) Login() revel.Result {
	return this.Render()
}

func (this *Account) Logout() revel.Result {
	fmt.Println("登出用户昵称: ", this.Session["nickName"])
	for k := range this.Session {
		delete(this.Session, k)
	}
	return this.Redirect((*App).Index)
}

func (this *Account) PostRegister(user *models.MockUser) revel.Result {
	//this.Validation.Email(user.Email).Message("电子邮件格式无效")
	this.Validation.Required(user.NickName).Message("用户昵称不能为空")
	this.Validation.Required(user.Password).Message("密码不能为空")
	this.Validation.Required(user.ConfirmPassword).Message("确认密码不能为空")
	this.Validation.MinSize(user.Password, 6).Message("密码长度不短于6位")
	this.Validation.Required(user.ConfirmPassword == user.Password).Message("两次输入的密码不一致")

	if this.Validation.HasErrors() {
		this.Validation.Keep()
		this.FlashParams()
		return this.Redirect((*Account).Register)
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

		// 添加错误信息，显示在页面的用户名下面
		var e revel.ValidationError
		e.Message = err.Error()
		e.Key = "user.NickName"
		this.Validation.Errors = append(this.Validation.Errors, &e)

		this.Validation.Keep()
		this.FlashParams()
		return this.Redirect((*Account).Register)
	}

	return this.Redirect((*Account).RegisterSuccessful)
}

func (this *Account) PostLogin(loginUser *models.LoginUser) revel.Result {
	this.Validation.Required(loginUser.NickName).Message("请输入昵称")
	this.Validation.Required(loginUser.Password).Message("请输入密码")

	if this.Validation.HasErrors() {
		this.Validation.Keep()
		this.FlashParams()
		return this.Redirect((*Account).Login)
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
		if err.Error() == "该用户不存在" {
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

	fmt.Println("使用昵称登陆: ", loginUser.NickName)

	return this.Redirect((*App).Index)
}
