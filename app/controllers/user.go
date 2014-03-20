package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/robfig/revel"
)

type User struct {
	*revel.Controller
}

func (user *User) Index() revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) Info() revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) Message() revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) Friend() revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

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

	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) EditArticle() revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) ShowArticle(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) PostEditArticle(article *models.UserArticle) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) EditInfo() revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) PostEditInfo() revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}
