package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/revel/revel"
	"strconv"
)

type Essay struct {
	*revel.Controller
}

func (this *Essay) Index() revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	allEssays, err := manager.GetAllEssay()

	// 倒序处理
	count := len(allEssays)
	for i := 0; i < count/2; i++ {
		allEssays[i], allEssays[count-i-1] = allEssays[count-i-1], allEssays[i]
	}

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
		essaysOnOnePage = allEssays[:models.ArticlesInSinglePage]
	} else {
		essaysOnOnePage = allEssays
	}

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["allEssays"] = allEssays
	this.RenderArgs["essaysOnOnePage"] = essaysOnOnePage
	this.RenderArgs["pageCount"] = pageCount
	this.RenderArgs["pageSlice"] = pageSlice

	return this.Render()
}

func (this *Essay) TypeIndex(tag string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	allEssays, err := manager.GetEssayByTag(tag)

	// 倒序处理
	count := len(allEssays)
	for i := 0; i < count/2; i++ {
		allEssays[i], allEssays[count-i-1] = allEssays[count-i-1], allEssays[i]
	}

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
		essaysOnOnePage = allEssays[:models.ArticlesInSinglePage]
	} else {
		essaysOnOnePage = allEssays
	}

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["allEssays"] = allEssays
	this.RenderArgs["essaysOnOnePage"] = essaysOnOnePage
	this.RenderArgs["pageCount"] = pageCount
	this.RenderArgs["type"] = tag
	this.RenderArgs["pageSlice"] = pageSlice

	return this.Render()
}

func (this *Essay) Add() revel.Result {
	userid := this.Session["userid"]
	nickName := this.Session["nickName"]

	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Render(userid, nickName)
}

func (this *Essay) Edit(id string) revel.Result {
	userid := this.Session["userid"]
	nickName := this.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()
	oringinalEssay, _ := manager.GetEssayById(id)

	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	return this.Render(userid, nickName, oringinalEssay)
}

func (this *Essay) PostAdd(essay *models.Essay) revel.Result {
	this.Validation.Required(essay.Tag).Message("请选择一个标签")
	this.Validation.Required(essay.Content).Message("散文内容不能为空")
	this.Validation.Required(essay.Author).Message("作者不能为空")

	fmt.Println("散文标签： ", essay.Tag)
	fmt.Println("散文标题： ", essay.Title)
	fmt.Println("散文内容： ", essay.Content)
	fmt.Println("作者： ", essay.Author)

	if this.Validation.HasErrors() {
		this.Validation.Keep()
		this.FlashParams()
		return this.Redirect((*Essay).Add)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()

	err = manager.AddEssay(essay)
	if err != nil {
		this.Flash.Error(err.Error())
		return this.Redirect((*Essay).Add)
	}

	// 返回到管理员首页
	return this.Redirect((*App).Index)
}

func (this *Essay) PostEdit(originalEssayID string, newEssay *models.Essay) revel.Result {
	this.Validation.Required(newEssay.Tag).Message("请选择一个标签")
	this.Validation.Required(newEssay.Content).Message("摘录内容不能为空")
	this.Validation.Required(newEssay.Author).Message("作者不能为空")

	fmt.Println("散文标签： ", newEssay.Tag)
	fmt.Println("散文标题： ", newEssay.Title)
	fmt.Println("散文内容： ", newEssay.Content)
	fmt.Println("作者： ", newEssay.Author)

	if this.Validation.HasErrors() {
		this.Validation.Keep()
		this.FlashParams()
		return this.Redirect((*Essay).Edit)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()

	err = manager.UpdateEssay(originalEssayID, newEssay)
	if err != nil {
		this.Flash.Error(err.Error())
		return this.Redirect((*Essay).Edit)
	}

	return this.Redirect((*Essay).Index)
}

func (this *Essay) Show(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()
	essay, _ := manager.GetEssayById(id)
	// if err != nil {
	// 	this.Flash.Error(err.Error())
	// 	//return this.Redirect((*Essay).Add)
	// }
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["essay"] = essay

	return this.Render()
}

func (this *Essay) PageList(pageNumber string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	allEssays, err := manager.GetAllEssay()

	// 倒序处理
	count := len(allEssays)
	for i := 0; i < count/2; i++ {
		allEssays[i], allEssays[count-i-1] = allEssays[count-i-1], allEssays[i]
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

	this.RenderArgs["allEssays"] = allEssays
	this.RenderArgs["essaysOnOnePage"] = essaysOnOnePage
	this.RenderArgs["pageCount"] = pageCount
	this.RenderArgs["pageNumber"] = pageNumber

	return this.Render()
}

func (this *Essay) PageListWithTag(uPageNumber string, tag string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	allEssays, err := manager.GetEssayByTag(tag)

	// 倒序处理
	count := len(allEssays)
	for i := 0; i < count/2; i++ {
		allEssays[i], allEssays[count-i-1] = allEssays[count-i-1], allEssays[i]
	}

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

	this.RenderArgs["allEssays"] = allEssays
	this.RenderArgs["essaysOnOnePage"] = essaysOnOnePage
	this.RenderArgs["pageCount"] = pageCount
	this.RenderArgs["uPageNumber"] = uPageNumber

	return this.Render()
}

func (this *Essay) Delete(id string) revel.Result {
	userid := this.Session["userid"]
	nickName := this.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()
	err = manager.DeleteEssayById(id)

	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	this.Render(userid, nickName)
	return this.Redirect((*Essay).Index)
}
