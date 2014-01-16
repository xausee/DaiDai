package controllers

import (
	"ZhaiLuBaiKe/app/models"
	"fmt"
	"github.com/robfig/revel"
)

type ModernPoem struct {
	*revel.Controller
}

func (mp *ModernPoem) Add() revel.Result {
	email := mp.Session["email"]
	nickName := mp.Session["nickName"]
	return mp.Render(email, nickName)
}

func (mp *ModernPoem) PostAdd(modernPoem *models.ModernPoem) revel.Result {
	mp.Validation.Required(modernPoem.Tag).Message("请选择一个标签")
	mp.Validation.Required(modernPoem.Content).Message("摘录内容不能为空")
	mp.Validation.Required(modernPoem.Author).Message("作者不能为空")

	fmt.Println("诗歌标签： ", modernPoem.Tag)
	fmt.Println("诗歌标题： ", modernPoem.Title)
	fmt.Println("诗歌内容： ", modernPoem.Content)
	fmt.Println("作者： ", modernPoem.Author)

	if mp.Validation.HasErrors() {
		mp.Validation.Keep()
		mp.FlashParams()
		return mp.Redirect((*ModernPoem).Add)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		mp.Response.Status = 500
		return mp.RenderError(err)
	}
	defer manager.Close()

	err = manager.AddModernPeom(modernPoem)
	if err != nil {
		mp.Flash.Error(err.Error())
		return mp.Redirect((*ModernPoem).Add)
	}

	return mp.Redirect((*App).Add)
}
