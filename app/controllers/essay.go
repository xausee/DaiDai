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
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	allEssays, err := manager.GetAllEssay()
	count := len(allEssays)

	var pageCount int
	if (count % 30) == 0 {
		pageCount = count / 30
	} else {
		pageCount = count/30 + 1
	}

	// var more bool
	// if count > 30 {
	// 	more = true
	// 	essaysOnOnePage := allEssays[(count - 30):]
	// } else {
	// 	more = false
	// 	essaysOnOnePage := allEssays
	// }

	e.RenderArgs["email"] = e.Session["email"]
	e.RenderArgs["nickName"] = e.Session["nickName"]
	e.RenderArgs["allEssays"] = allEssays
	e.RenderArgs["pageCount"] = pageCount

	return e.Render()
}

func (e *Essay) TypeIndex(tag string) revel.Result {
	email := e.Session["email"]
	nickName := e.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	allEssays, err := manager.GetEssayByTag(tag)
	count := len(allEssays)

	var pageCount int
	if (count % 30) == 0 {
		pageCount = count / 30
	} else {
		pageCount = count/30 + 1
	}

	e.RenderArgs["email"] = e.Session["email"]
	e.RenderArgs["nickName"] = e.Session["nickName"]
	e.RenderArgs["allEssays"] = allEssays
	e.RenderArgs["pageCount"] = pageCount

	return e.Render(email, nickName, allEssays)
}

func (e *Essay) Add() revel.Result {
	email := e.Session["email"]
	nickName := e.Session["nickName"]
	return e.Render(email, nickName)
}

func (e *Essay) Edit(id string) revel.Result {
	email := e.Session["email"]
	nickName := e.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		e.Response.Status = 500
		return e.RenderError(err)
	}
	defer manager.Close()
	oringinalEssay, _ := manager.GetEssayById(id)

	return e.Render(email, nickName, oringinalEssay)
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

func (e *Essay) PostEdit(originalEssayID string, newEssay *models.Essay) revel.Result {
	e.Validation.Required(newEssay.Tag).Message("请选择一个标签")
	e.Validation.Required(newEssay.Content).Message("摘录内容不能为空")
	e.Validation.Required(newEssay.Author).Message("作者不能为空")

	fmt.Println("散文标签： ", newEssay.Tag)
	fmt.Println("散文标题： ", newEssay.Title)
	fmt.Println("散文内容： ", newEssay.Content)
	fmt.Println("作者： ", newEssay.Author)

	if e.Validation.HasErrors() {
		e.Validation.Keep()
		e.FlashParams()
		return e.Redirect((*Essay).Edit)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		e.Response.Status = 500
		return e.RenderError(err)
	}
	defer manager.Close()

	err = manager.UpdateEssay(originalEssayID, newEssay)
	if err != nil {
		e.Flash.Error(err.Error())
		return e.Redirect((*Essay).Edit)
	}

	return e.Redirect((*Essay).Index)
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

func (e *Essay) PageList(pageNumber string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	allEssays, err := manager.GetAllEssay()
	count := len(allEssays)

	var pageCount int
	if (count % 30) == 0 {
		pageCount = count / 30
	} else {
		pageCount = count/30 + 1
	}

	e.RenderArgs["allEssays"] = allEssays
	e.RenderArgs["pageCount"] = pageCount
	e.RenderArgs["pageNumber"] = pageNumber

	return e.Render()
}

func (e *Essay) Delete(id string) revel.Result {
	email := e.Session["email"]
	nickName := e.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		e.Response.Status = 500
		return e.RenderError(err)
	}
	defer manager.Close()
	err = manager.DeleteEssayById(id)

	e.Render(email, nickName)
	return e.Redirect((*Essay).Index)
}
