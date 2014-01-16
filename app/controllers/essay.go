package controllers

import (
	"ZhaiLuBaiKe/app/models"
	"fmt"
	"github.com/robfig/revel"
)

type Essay struct {
	*revel.Controller
}

func (e *Essay) Add() revel.Result {
	email := e.Session["email"]
	nickName := e.Session["nickName"]
	return e.Render(email, nickName)
}

func (e *Essay) PostAdd(essay *models.Essay) revel.Result {
	e.Validation.Required(essay.Tag).Message("请选择一个标签")
	e.Validation.Required(essay.Content).Message("摘录内容不能为空")
	e.Validation.Required(essay.Author).Message("作者不能为空")

	fmt.Println("诗歌标签： ", essay.Tag)
	fmt.Println("诗歌标题： ", essay.Title)
	fmt.Println("诗歌内容： ", essay.Content)
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
