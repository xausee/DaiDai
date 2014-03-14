package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/robfig/revel"
	"strconv"
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
	if (count % models.ArticlesInSinglePage) == 0 {
		pageCount = count / models.ArticlesInSinglePage
	} else {
		pageCount = count/models.ArticlesInSinglePage + 1
	}

	essaysOnOnePage := []models.Essay{}
	if count > models.ArticlesInSinglePage {
		essaysOnOnePage = allEssays[(count - models.ArticlesInSinglePage):]
	} else {
		essaysOnOnePage = allEssays
	}

	e.RenderArgs["email"] = e.Session["email"]
	e.RenderArgs["nickName"] = e.Session["nickName"]
	e.RenderArgs["allEssays"] = allEssays
	e.RenderArgs["essaysOnOnePage"] = essaysOnOnePage
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
	if (count % models.ArticlesInSinglePage) == 0 {
		pageCount = count / models.ArticlesInSinglePage
	} else {
		pageCount = count/models.ArticlesInSinglePage + 1
	}

	pageSlice := make([]int, 0)
	for i := 1; i <= pageCount; i++ {
		pageSlice = append(pageSlice, i)
	}

	essaysOnOnePage := []models.Essay{}
	if count > models.ArticlesInSinglePage {
		essaysOnOnePage = allEssays[(count - models.ArticlesInSinglePage):]
	} else {
		essaysOnOnePage = allEssays
	}

	e.RenderArgs["email"] = e.Session["email"]
	e.RenderArgs["nickName"] = e.Session["nickName"]
	e.RenderArgs["allEssays"] = allEssays
	e.RenderArgs["essaysOnOnePage"] = essaysOnOnePage
	e.RenderArgs["pageCount"] = pageCount
	e.RenderArgs["type"] = tag
	e.RenderArgs["pageSlice"] = pageSlice

	return e.Render(email, nickName)
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

func (e *Essay) PostAdd(essay *models.Essay) revel.Result {
	e.Validation.Required(essay.Tag).Message("请选择一个标签")
	e.Validation.Required(essay.Content).Message("散文内容不能为空")
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
	if (count % models.ArticlesInSinglePage) == 0 {
		pageCount = count / models.ArticlesInSinglePage
	} else {
		pageCount = count/models.ArticlesInSinglePage + 1
	}

	var iPageNumber int
	iPageNumber, err = strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Println(err)
	}

	essaysOnOnePage := []models.Essay{}
	if count >= iPageNumber*models.ArticlesInSinglePage {
		essaysOnOnePage = allEssays[(iPageNumber-1)*models.ArticlesInSinglePage : iPageNumber*models.ArticlesInSinglePage]
	} else if (iPageNumber-1)*models.ArticlesInSinglePage < count && count < iPageNumber*models.ArticlesInSinglePage {
		essaysOnOnePage = allEssays[(iPageNumber-1)*models.ArticlesInSinglePage:]
	} else if (iPageNumber-1)*models.ArticlesInSinglePage > count {
		fmt.Println("已超过最大页码")
	}
	fmt.Println("pageCount:", pageCount)
	fmt.Println("pageNumber:", pageNumber)

	e.RenderArgs["allEssays"] = allEssays
	e.RenderArgs["essaysOnOnePage"] = essaysOnOnePage
	e.RenderArgs["pageCount"] = pageCount
	e.RenderArgs["pageNumber"] = pageNumber

	return e.Render()
}

func (e *Essay) PageListWithTag(uPageNumber string, tag string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	allEssays, err := manager.GetEssayByTag(tag)
	count := len(allEssays)

	var pageCount int
	if (count % models.ArticlesInSinglePage) == 0 {
		pageCount = count / models.ArticlesInSinglePage
	} else {
		pageCount = count/models.ArticlesInSinglePage + 1
	}

	var iPageNumber int
	iPageNumber, err = strconv.Atoi(uPageNumber)
	if err != nil {
		fmt.Println(err)
	}

	essaysOnOnePage := []models.Essay{}
	if count >= iPageNumber*models.ArticlesInSinglePage {
		essaysOnOnePage = allEssays[(iPageNumber-1)*models.ArticlesInSinglePage : iPageNumber*models.ArticlesInSinglePage]
	} else if (iPageNumber-1)*models.ArticlesInSinglePage < count && count < iPageNumber*models.ArticlesInSinglePage {
		essaysOnOnePage = allEssays[(iPageNumber-1)*models.ArticlesInSinglePage:]
	} else if (iPageNumber-1)*models.ArticlesInSinglePage > count {
		fmt.Println("已超过最大页码")
	}
	fmt.Println("pageCount:", pageCount)
	fmt.Println("uPageNumber:", uPageNumber)

	e.RenderArgs["allEssays"] = allEssays
	e.RenderArgs["essaysOnOnePage"] = essaysOnOnePage
	e.RenderArgs["pageCount"] = pageCount
	e.RenderArgs["uPageNumber"] = uPageNumber

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
