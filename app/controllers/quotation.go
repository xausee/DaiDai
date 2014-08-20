package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/revel/revel"
	"strconv"
)

type Quotation struct {
	*revel.Controller
}

func (this *Quotation) Index() revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	quotations, err := manager.GetAllQuotation()

	// 倒序处理
	count := len(quotations)
	for i := 0; i < count/2; i++ {
		quotations[i], quotations[count-i-1] = quotations[count-i-1], quotations[i]
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

	quotationsOnOnePage := []models.Quotation{}
	if count > models.ArticlesInSinglePage {
		quotationsOnOnePage = quotations[:models.ArticlesInSinglePage]
	} else {
		quotationsOnOnePage = quotations
	}

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["allQuotations"] = quotations
	this.RenderArgs["quotationsOnOnePage"] = quotationsOnOnePage
	this.RenderArgs["pageCount"] = pageCount
	this.RenderArgs["pages"] = pages

	return this.Render()
}

func (this *Quotation) TypeIndex(tag string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	quotations, err := manager.GetQuotationByTag(tag)

	// 倒序处理
	count := len(quotations)
	for i := 0; i < count/2; i++ {
		quotations[i], quotations[count-i-1] = quotations[count-i-1], quotations[i]
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

	quotationsOnOnePage := []models.Quotation{}
	if count > models.ArticlesInSinglePage {
		quotationsOnOnePage = quotations[:models.ArticlesInSinglePage]
	} else {
		quotationsOnOnePage = quotations
	}

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["allQuotations"] = quotations
	this.RenderArgs["quotationsOnOnePage"] = quotationsOnOnePage
	this.RenderArgs["pageCount"] = pageCount
	this.RenderArgs["type"] = tag
	this.RenderArgs["pages"] = pages

	return this.Render()
}

func (this *Quotation) Add() revel.Result {
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]

	return this.Render()
}

func (this *Quotation) Edit(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()
	originalQuotation, _ := manager.GetQuotationById(id)

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["originalQuotation"] = originalQuotation

	return this.Render()
}

func (this *Quotation) PostAdd(quotation *models.Quotation) revel.Result {
	this.Validation.Required(quotation.Tag).Message("请选择一个标签")
	this.Validation.Required(quotation.Content).Message("摘录内容不能为空")
	this.Validation.Required(quotation.Author).Message("作者不能为空")

	fmt.Println("摘录标签： ", quotation.Tag)
	fmt.Println("摘录被容： ", quotation.Content)
	fmt.Println("原文： ", quotation.Original)
	fmt.Println("作者： ", quotation.Author)

	if this.Validation.HasErrors() {
		this.Validation.Keep()
		this.FlashParams()
		return this.Redirect((*Quotation).Add)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()

	err = manager.AddQuotation(quotation)
	if err != nil {
		// this.Validation.Keep()
		// this.FlashParams()
		this.Flash.Error(err.Error())
		return this.Redirect((*Quotation).Add)
	}

	// 返回到管理员首页
	return this.Redirect((*Admin).Index)
}

func (this *Quotation) PostEdit(originalQuotationID string, newQuotation *models.Quotation) revel.Result {
	this.Validation.Required(newQuotation.Tag).Message("请选择一个标签")
	this.Validation.Required(newQuotation.Content).Message("摘录内容不能为空")
	this.Validation.Required(newQuotation.Author).Message("作者不能为空")

	fmt.Println("摘录标签： ", newQuotation.Tag)
	fmt.Println("摘录被容： ", newQuotation.Content)
	fmt.Println("原文： ", newQuotation.Original)
	fmt.Println("作者： ", newQuotation.Author)

	if this.Validation.HasErrors() {
		this.Validation.Keep()
		this.FlashParams()
		fmt.Println("error in validation ")
		return this.Redirect((*Quotation).Edit)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()

	err = manager.UpdateQuotation(originalQuotationID, newQuotation)
	if err != nil {
		// this.Validation.Keep()
		// this.FlashParams()
		this.Flash.Error(err.Error())
		fmt.Println("error in update quotation ")
		return this.Redirect((*Quotation).Edit)
	}

	return this.Redirect((*Quotation).Index)
}

func (this *Quotation) Show(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()
	quotation, _ := manager.GetQuotationById(id)
	// if err != nil {
	// 	this.Flash.Error(err.Error())
	// 	//return this.Redirect((*Essay).Add)
	// }

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["quotation"] = quotation

	return this.Render()
}

func (this *Quotation) PageList(pageNumber string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	allQuotations, err := manager.GetAllQuotation()

	// 倒序处理
	count := len(allQuotations)
	for i := 0; i < count/2; i++ {
		allQuotations[i], allQuotations[count-i-1] = allQuotations[count-i-1], allQuotations[i]
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

	quotationsOnOnePage := []models.Quotation{}
	if count >= iPageNumber*models.ArticlesInSinglePage {
		quotationsOnOnePage = allQuotations[(iPageNumber-1)*models.ArticlesInSinglePage : iPageNumber*models.ArticlesInSinglePage]
	} else if (iPageNumber-1)*models.ArticlesInSinglePage < count && count < iPageNumber*models.ArticlesInSinglePage {
		quotationsOnOnePage = allQuotations[(iPageNumber-1)*models.ArticlesInSinglePage:]
	} else if (iPageNumber-1)*models.ArticlesInSinglePage > count {
		fmt.Println("已超过最大页码")
	}
	fmt.Println("pageCount:", pageCount)
	fmt.Println("pageNumber:", pageNumber)

	this.RenderArgs["allQuotations"] = allQuotations
	this.RenderArgs["quotationsOnOnePage"] = quotationsOnOnePage
	this.RenderArgs["pageCount"] = pageCount
	this.RenderArgs["pageNumber"] = pageNumber

	return this.Render()
}

func (this *Quotation) PageListWithTag(uPageNumber string, tag string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	allQuotations, err := manager.GetQuotationByTag(tag)

	// 倒序处理
	count := len(allQuotations)
	for i := 0; i < count/2; i++ {
		allQuotations[i], allQuotations[count-i-1] = allQuotations[count-i-1], allQuotations[i]
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

	quotationsOnOnePage := []models.Quotation{}
	if count >= iPageNumber*models.ArticlesInSinglePage {
		quotationsOnOnePage = allQuotations[(iPageNumber-1)*models.ArticlesInSinglePage : iPageNumber*models.ArticlesInSinglePage]
	} else if (iPageNumber-1)*models.ArticlesInSinglePage < count && count < iPageNumber*models.ArticlesInSinglePage {
		quotationsOnOnePage = allQuotations[(iPageNumber-1)*models.ArticlesInSinglePage:]
	} else if (iPageNumber-1)*models.ArticlesInSinglePage > count {
		fmt.Println("已超过最大页码")
	}
	fmt.Println("pageCount:", pageCount)
	fmt.Println("uPageNumber:", uPageNumber)

	this.RenderArgs["allQuotations"] = allQuotations
	this.RenderArgs["quotationsOnOnePage"] = quotationsOnOnePage
	this.RenderArgs["pageCount"] = pageCount
	this.RenderArgs["uPageNumber"] = uPageNumber

	return this.Render()
}

func (this *Quotation) Delete(id string) revel.Result {
	userid := this.Session["userid"]
	nickName := this.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()
	err = manager.DeleteQuotationById(id)

	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	this.Render(userid, nickName)
	return this.Redirect((*Quotation).Index)
}
