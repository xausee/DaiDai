package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/revel/revel"
	"strconv"
)

type HintFiction struct {
	*revel.Controller
}

func (this *HintFiction) Index() revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	hintFictions, err := manager.GetAllHintFiction()

	// 倒序处理
	count := len(hintFictions)
	for i := 0; i < count/2; i++ {
		hintFictions[i], hintFictions[count-i-1] = hintFictions[count-i-1], hintFictions[i]
	}

	var pageCount int
	if (count % models.ArticlesInSinglePage) == 0 {
		pageCount = count / models.ArticlesInSinglePage
	} else {
		pageCount = count/models.ArticlesInSinglePage + 1
	}

	pages := make([]int, 0)
	for i := 1; i <= pageCount; i++ {
		pages = append(pages, i)
	}

	hintFictionsOnOnePage := []models.HintFiction{}
	if count > models.ArticlesInSinglePage {
		hintFictionsOnOnePage = hintFictions[:models.ArticlesInSinglePage]
	} else {
		hintFictionsOnOnePage = hintFictions
	}

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["allHintFinctions"] = hintFictions
	this.RenderArgs["hintFictionsOnOnePage"] = hintFictionsOnOnePage
	this.RenderArgs["pageCount"] = pageCount
	this.RenderArgs["pages"] = pages

	return this.Render()
}

func (this *HintFiction) TypeIndex(tag string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	hintFictions, err := manager.GetHintFictionByTag(tag)

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["hintFictions"] = hintFictions

	return this.Render()
}

func (this *HintFiction) Add() revel.Result {
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Render()
}

func (this *HintFiction) Edit(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()
	oringinalHintFiction, _ := manager.GetHintFictionById(id)

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["oringinalHintFiction"] = oringinalHintFiction

	return this.Render()
}

func (this *HintFiction) PostAdd(hintFiction *models.HintFiction) revel.Result {
	this.Validation.Required(hintFiction.Tag).Message("请选择一个标签")
	this.Validation.Required(hintFiction.Content).Message("内容不能为空")
	this.Validation.Required(hintFiction.Author).Message("作者不能为空")

	fmt.Println("微小说标签： ", hintFiction.Tag)
	fmt.Println("微小说标题： ", hintFiction.Title)
	fmt.Println("微小说内容： ", hintFiction.Content)
	fmt.Println("微小说作者： ", hintFiction.Author)

	if this.Validation.HasErrors() {
		this.Validation.Keep()
		this.FlashParams()
		return this.Redirect((*HintFiction).Add)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()

	err = manager.AddHintFiction(hintFiction)
	if err != nil {
		this.Flash.Error(err.Error())
		return this.Redirect((*HintFiction).Add)
	}

	return this.Redirect((*App).Add)
}

func (this *HintFiction) PostEdit(originalHintFictionID string, newHintFiction *models.HintFiction) revel.Result {
	this.Validation.Required(newHintFiction.Tag).Message("请选择一个标签")
	this.Validation.Required(newHintFiction.Content).Message("内容不能为空")
	this.Validation.Required(newHintFiction.Author).Message("作者不能为空")

	fmt.Println("微小说标签： ", newHintFiction.Tag)
	fmt.Println("微小说标题： ", newHintFiction.Title)
	fmt.Println("微小说内容： ", newHintFiction.Content)
	fmt.Println("微小说作者： ", newHintFiction.Author)

	if this.Validation.HasErrors() {
		this.Validation.Keep()
		this.FlashParams()
		return this.Redirect((*HintFiction).Edit)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()

	err = manager.UpdateHintFiction(originalHintFictionID, newHintFiction)
	if err != nil {
		this.Flash.Error(err.Error())
		return this.Redirect((*HintFiction).Edit)
	}

	return this.Redirect((*HintFiction).Index)
}

func (this *HintFiction) Show(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()
	hintFiction, _ := manager.GetHintFictionById(id)
	// if err != nil {
	// 	this.Flash.Error(err.Error())
	// 	//return this.Redirect((*Essay).Add)
	// }

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["hintFiction"] = hintFiction

	return this.Render()
}

func (this *HintFiction) PageList(pageNumber string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	allHintfictions, err := manager.GetAllHintFiction()

	// 倒序处理
	count := len(allHintfictions)
	for i := 0; i < count/2; i++ {
		allHintfictions[i], allHintfictions[count-i-1] = allHintfictions[count-i-1], allHintfictions[i]
	}

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

	hintFictionsOnOnePage := []models.HintFiction{}
	if count >= iPageNumber*models.ArticlesInSinglePage {
		hintFictionsOnOnePage = allHintfictions[(iPageNumber-1)*models.ArticlesInSinglePage : iPageNumber*models.ArticlesInSinglePage]
	} else if (iPageNumber-1)*models.ArticlesInSinglePage < count && count < iPageNumber*models.ArticlesInSinglePage {
		hintFictionsOnOnePage = allHintfictions[(iPageNumber-1)*models.ArticlesInSinglePage:]
	} else if (iPageNumber-1)*models.ArticlesInSinglePage > count {
		fmt.Println("已超过最大页码")
	}
	fmt.Println("pageCount:", pageCount)
	fmt.Println("pageNumber:", pageNumber)

	this.RenderArgs["allHintfictions"] = allHintfictions
	this.RenderArgs["hintFictionsOnOnePage"] = hintFictionsOnOnePage
	this.RenderArgs["pageCount"] = pageCount
	this.RenderArgs["pageNumber"] = pageNumber

	return this.Render()
}

func (this *HintFiction) Delete(id string) revel.Result {
	userid := this.Session["userid"]
	nickName := this.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()
	err = manager.DeleteHintFictionById(id)

	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	this.Render(userid, nickName)
	return this.Redirect((*HintFiction).Index)
}
