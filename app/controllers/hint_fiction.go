package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/robfig/revel"
)

type HintFiction struct {
	*revel.Controller
}

func (hf *HintFiction) Index() revel.Result {
	email := hf.Session["email"]
	nickName := hf.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	hintFictions, err := manager.GetAllHintFiction()

	return hf.Render(email, nickName, hintFictions)
}

func (hf *HintFiction) Add() revel.Result {
	email := hf.Session["email"]
	nickName := hf.Session["nickName"]
	return hf.Render(email, nickName)
}

func (hf *HintFiction) PostAdd(hintFiction *models.HintFiction) revel.Result {
	hf.Validation.Required(hintFiction.Tag).Message("请选择一个标签")
	hf.Validation.Required(hintFiction.Content).Message("内容不能为空")
	hf.Validation.Required(hintFiction.Author).Message("作者不能为空")

	fmt.Println("微小说标签： ", hintFiction.Tag)
	fmt.Println("微小说标题： ", hintFiction.Title)
	fmt.Println("微小说内容： ", hintFiction.Content)
	fmt.Println("微小说作者： ", hintFiction.Author)

	if hf.Validation.HasErrors() {
		hf.Validation.Keep()
		hf.FlashParams()
		return hf.Redirect((*HintFiction).Add)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		hf.Response.Status = 500
		return hf.RenderError(err)
	}
	defer manager.Close()

	err = manager.AddHintFiction(hintFiction)
	if err != nil {
		hf.Flash.Error(err.Error())
		return hf.Redirect((*HintFiction).Add)
	}

	return hf.Redirect((*App).Add)
}

func (hf *HintFiction) Show(id string) revel.Result {
	email := hf.Session["email"]
	nickName := hf.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		hf.Response.Status = 500
		return hf.RenderError(err)
	}
	defer manager.Close()
	hintFiction, _ := manager.GetHintFictionById(id)
	// if err != nil {
	// 	hf.Flash.Error(err.Error())
	// 	//return hf.Redirect((*Essay).Add)
	// }
	return hf.Render(email, nickName, hintFiction)
}
