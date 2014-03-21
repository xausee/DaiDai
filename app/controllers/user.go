package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/robfig/revel"
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

	user.RenderArgs["id"] = user.Session["id"]
	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}

func (user *User) Info(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	user.RenderArgs["id"] = user.Session["id"]
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

	user.RenderArgs["id"] = user.Session["id"]
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

	user.RenderArgs["id"] = user.Session["id"]
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

	user.RenderArgs["id"] = user.Session["id"]
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

	err = manager.AddUserArticle(article)
	if err != nil {
		user.Validation.Keep()
		user.FlashParams()
		user.Flash.Error(err.Error())
		return user.Redirect((*User).AddArticle)
	}

	user.RenderArgs["id"] = user.Session["id"]
	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Redirect((*User).AddArticle)
}

func (user *User) EditArticle() revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	user.RenderArgs["id"] = user.Session["id"]
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

	user.RenderArgs["id"] = user.Session["id"]
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

	user.RenderArgs["id"] = user.Session["id"]
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

	user.RenderArgs["id"] = user.Session["id"]
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

	user.RenderArgs["id"] = user.Session["id"]
	user.RenderArgs["email"] = user.Session["email"]
	user.RenderArgs["nickName"] = user.Session["nickName"]

	return user.Render()
}
