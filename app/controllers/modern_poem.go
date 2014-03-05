package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/robfig/revel"
)

type ModernPoem struct {
	*revel.Controller
}

func (mp *ModernPoem) Index() revel.Result {
	email := mp.Session["email"]
	nickName := mp.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	modernPoems, err := manager.GetAllModernPoem()

	return mp.Render(email, nickName, modernPoems)
}

func (mp *ModernPoem) TypeIndex(tag string) revel.Result {
	email := mp.Session["email"]
	nickName := mp.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	//modernPoems, err := manager.GetAllModernPoem()
	modernPoems, err := manager.GetModernPoemByTag(tag)

	return mp.Render(email, nickName, modernPoems)
}

func (mp *ModernPoem) Add() revel.Result {
	email := mp.Session["email"]
	nickName := mp.Session["nickName"]
	return mp.Render(email, nickName)
}

func (mp *ModernPoem) Edit(id string) revel.Result {
	email := mp.Session["email"]
	nickName := mp.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		mp.Response.Status = 500
		return mp.RenderError(err)
	}
	defer manager.Close()
	oringinalModernPoem, _ := manager.GetModernPoemById(id)

	return mp.Render(email, nickName, oringinalModernPoem)
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

func (mp *ModernPoem) PostEdit(originalModernPoemID string, newModernPoem *models.ModernPoem) revel.Result {
	mp.Validation.Required(newModernPoem.Tag).Message("请选择一个标签")
	mp.Validation.Required(newModernPoem.Content).Message("摘录内容不能为空")
	mp.Validation.Required(newModernPoem.Author).Message("作者不能为空")

	fmt.Println("诗歌标签： ", newModernPoem.Tag)
	fmt.Println("诗歌标题： ", newModernPoem.Title)
	fmt.Println("诗歌内容： ", newModernPoem.Content)
	fmt.Println("作者： ", newModernPoem.Author)

	if mp.Validation.HasErrors() {
		mp.Validation.Keep()
		mp.FlashParams()
		return mp.Redirect((*ModernPoem).Edit)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		mp.Response.Status = 500
		return mp.RenderError(err)
	}
	defer manager.Close()

	err = manager.UpdateModernPeom(originalModernPoemID, newModernPoem)
	if err != nil {
		mp.Flash.Error(err.Error())
		return mp.Redirect((*ModernPoem).Edit)
	}

	return mp.Redirect((*ModernPoem).Index)
}

func (mp *ModernPoem) Show(id string) revel.Result {
	email := mp.Session["email"]
	nickName := mp.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		mp.Response.Status = 500
		return mp.RenderError(err)
	}
	defer manager.Close()
	modernPoem, _ := manager.GetModernPoemById(id)
	// if err != nil {
	// 	mp.Flash.Error(err.Error())
	// 	//return mp.Redirect((*Essay).Add)
	// }
	return mp.Render(email, nickName, modernPoem)
}

func (mp *ModernPoem) Delete(id string) revel.Result {
	email := mp.Session["email"]
	nickName := mp.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		mp.Response.Status = 500
		return mp.RenderError(err)
	}
	defer manager.Close()
	err = manager.DeleteModernPoemById(id)

	mp.Render(email, nickName)
	return mp.Redirect((*ModernPoem).Index)
}
