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

func (c *Account) Register() revel.Result {
	return c.Render()
}

func (c *Account) RegisterSuccessful() revel.Result {
	return c.Render()
}

func (c *Account) Login() revel.Result {
	return c.Render()
}

func (c *Account) Logout() revel.Result {
	fmt.Println("登出用户昵称: ", c.Session["nickName"])
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect((*App).Index)
}

func (c *Account) PostRegister(user *models.MockUser) revel.Result {
	//c.Validation.Email(user.Email).Message("电子邮件格式无效")
	c.Validation.Required(user.NickName).Message("用户昵称不能为空")
	c.Validation.Required(user.Password).Message("密码不能为空")
	c.Validation.MinSize(user.Password, 6).Message("密码长度不短于6位")
	c.Validation.Required(user.ConfirmPassword == user.Password).Message("两次输入的密码不一致")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect((*Account).Register)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		c.Response.Status = 500
		return c.RenderError(err)
	}
	defer manager.Close()

	err = manager.RegisterUser(user)
	if err != nil {
		// c.Validation.Keep()
		// c.FlashParams()
		c.Flash.Error(err.Error())
		return c.Redirect((*Account).Register)
	}

	return c.Redirect((*Account).RegisterSuccessful)
}

func (c *Account) PostLogin(loginUser *models.LoginUser) revel.Result {
	c.Validation.Required(loginUser.NickName).Message("请输入昵称")
	c.Validation.Required(loginUser.Password).Message("请输入密码")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect((*Account).Login)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		c.Response.Status = 500
		return c.RenderError(err)
	}
	defer manager.Close()

	var u *models.User
	u, err = manager.LoginUser(loginUser)
	if err != nil {
		//c.Validation.Keep()
		// c.FlashParams()
		c.Flash.Error(err.Error())
		return c.Redirect((*Account).Login)
	}
	c.Session["userid"] = strconv.Itoa(u.Id)
	c.Session["nickName"] = u.NickName
	fmt.Println("使用昵称登陆: ", loginUser.NickName)

	return c.Redirect((*App).Index)
}
