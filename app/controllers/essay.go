package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/robfig/revel"
)

type Essay struct {
	*revel.Controller
}

func (e *Essay) Index() revel.Result {
	email := e.Session["email"]
	nickName := e.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	essays, err := manager.GetAllEssay()

	return e.Render(email, nickName, essays)
}

func (e *Essay) Add() revel.Result {
	email := e.Session["email"]
	nickName := e.Session["nickName"]
	return e.Render(email, nickName)
}

func (e *Essay) MinGuoEssay() revel.Result {
	email := e.Session["email"]
	nickName := e.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	essays, err := manager.GetAllEssay()

	return e.Render(email, nickName, essays)
}

func (e *Essay) DangDaiEssay() revel.Result {
	email := e.Session["email"]
	nickName := e.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	essays, err := manager.GetAllEssay()

	return e.Render(email, nickName, essays)
}

func (e *Essay) PostAdd(essay *models.Essay) revel.Result {
	e.Validation.Required(essay.Tag).Message("请选择一个标签")
	e.Validation.Required(essay.Content).Message("摘录内容不能为空")
	e.Validation.Required(essay.Author).Message("作者不能为空")

	fmt.Println("散文标签： ", essay.Tag)
	fmt.Println("散文标题： ", essay.Title)
	fmt.Println("散文内容： ", essay.Content)
	fmt.Println("作者： ", essay.Author)

	if e.Validation.HasErrors() {
		e.Validation.Keep()
		e.FlashParams()
		return e.Redirect((*Essay).Add)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		e.Response.Status = 500
		return e.RenderError(err)
	}
	defer manager.Close()

	err = manager.AddEssay(essay)
	if err != nil {
		e.Flash.Error(err.Error())
		return e.Redirect((*Essay).Add)
	}

	return e.Redirect((*App).Add)
}

func (e *Essay) Show(id string) revel.Result {
	email := e.Session["email"]
	nickName := e.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		e.Response.Status = 500
		return e.RenderError(err)
	}
	defer manager.Close()
	essay, _ := manager.GetEssayById(id)
	// if err != nil {
	// 	e.Flash.Error(err.Error())
	// 	//return e.Redirect((*Essay).Add)
	// }
	return e.Render(email, nickName, essay)
}
