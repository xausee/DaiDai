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

func (q *Quotation) Index() revel.Result {
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

	q.RenderArgs["userid"] = q.Session["userid"]
	q.RenderArgs["nickName"] = q.Session["nickName"]
	q.RenderArgs["avatarUrl"] = q.Session["avatarUrl"]
	q.RenderArgs["allQuotations"] = quotations
	q.RenderArgs["quotationsOnOnePage"] = quotationsOnOnePage
	q.RenderArgs["pageCount"] = pageCount
	q.RenderArgs["pages"] = pages

	return q.Render()
}

func (q *Quotation) TypeIndex(tag string) revel.Result {
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

	q.RenderArgs["userid"] = q.Session["userid"]
	q.RenderArgs["nickName"] = q.Session["nickName"]
	q.RenderArgs["avatarUrl"] = q.Session["avatarUrl"]
	q.RenderArgs["allQuotations"] = quotations
	q.RenderArgs["quotationsOnOnePage"] = quotationsOnOnePage
	q.RenderArgs["pageCount"] = pageCount
	q.RenderArgs["type"] = tag
	q.RenderArgs["pages"] = pages

	return q.Render()
}

func (q *Quotation) Add() revel.Result {
	q.RenderArgs["userid"] = q.Session["userid"]
	q.RenderArgs["nickName"] = q.Session["nickName"]

	return q.Render()
}

func (q *Quotation) Edit(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		q.Response.Status = 500
		return q.RenderError(err)
	}
	defer manager.Close()
	originalQuotation, _ := manager.GetQuotationById(id)

	q.RenderArgs["userid"] = q.Session["userid"]
	q.RenderArgs["nickName"] = q.Session["nickName"]
	q.RenderArgs["avatarUrl"] = q.Session["avatarUrl"]
	q.RenderArgs["originalQuotation"] = originalQuotation

	return q.Render()
}

func (q *Quotation) PostAdd(quotation *models.Quotation) revel.Result {
	q.Validation.Required(quotation.Tag).Message("请选择一个标签")
	q.Validation.Required(quotation.Content).Message("摘录内容不能为空")
	q.Validation.Required(quotation.Author).Message("作者不能为空")

	fmt.Println("摘录标签： ", quotation.Tag)
	fmt.Println("摘录被容： ", quotation.Content)
	fmt.Println("原文： ", quotation.Original)
	fmt.Println("作者： ", quotation.Author)

	if q.Validation.HasErrors() {
		q.Validation.Keep()
		q.FlashParams()
		return q.Redirect((*Quotation).Add)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		q.Response.Status = 500
		return q.RenderError(err)
	}
	defer manager.Close()

	err = manager.AddQuotation(quotation)
	if err != nil {
		// q.Validation.Keep()
		// q.FlashParams()
		q.Flash.Error(err.Error())
		return q.Redirect((*Quotation).Add)
	}

	return q.Redirect((*App).Add)
}

func (q *Quotation) PostEdit(originalQuotationID string, newQuotation *models.Quotation) revel.Result {
	q.Validation.Required(newQuotation.Tag).Message("请选择一个标签")
	q.Validation.Required(newQuotation.Content).Message("摘录内容不能为空")
	q.Validation.Required(newQuotation.Author).Message("作者不能为空")

	fmt.Println("摘录标签： ", newQuotation.Tag)
	fmt.Println("摘录被容： ", newQuotation.Content)
	fmt.Println("原文： ", newQuotation.Original)
	fmt.Println("作者： ", newQuotation.Author)

	if q.Validation.HasErrors() {
		q.Validation.Keep()
		q.FlashParams()
		fmt.Println("error in validation ")
		return q.Redirect((*Quotation).Edit)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		q.Response.Status = 500
		return q.RenderError(err)
	}
	defer manager.Close()

	err = manager.UpdateQuotation(originalQuotationID, newQuotation)
	if err != nil {
		// q.Validation.Keep()
		// q.FlashParams()
		q.Flash.Error(err.Error())
		fmt.Println("error in update quotation ")
		return q.Redirect((*Quotation).Edit)
	}

	return q.Redirect((*Quotation).Index)
}

func (q *Quotation) Show(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		q.Response.Status = 500
		return q.RenderError(err)
	}
	defer manager.Close()
	quotation, _ := manager.GetQuotationById(id)
	// if err != nil {
	// 	q.Flash.Error(err.Error())
	// 	//return q.Redirect((*Essay).Add)
	// }

	q.RenderArgs["userid"] = q.Session["userid"]
	q.RenderArgs["nickName"] = q.Session["nickName"]
	q.RenderArgs["avatarUrl"] = q.Session["avatarUrl"]
	q.RenderArgs["quotation"] = quotation

	return q.Render()
}

func (q *Quotation) PageList(pageNumber string) revel.Result {
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

	q.RenderArgs["allQuotations"] = allQuotations
	q.RenderArgs["quotationsOnOnePage"] = quotationsOnOnePage
	q.RenderArgs["pageCount"] = pageCount
	q.RenderArgs["pageNumber"] = pageNumber

	return q.Render()
}

func (q *Quotation) PageListWithTag(uPageNumber string, tag string) revel.Result {
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

	q.RenderArgs["allQuotations"] = allQuotations
	q.RenderArgs["quotationsOnOnePage"] = quotationsOnOnePage
	q.RenderArgs["pageCount"] = pageCount
	q.RenderArgs["uPageNumber"] = uPageNumber

	return q.Render()
}

func (q *Quotation) Delete(id string) revel.Result {
	userid := q.Session["userid"]
	nickName := q.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		q.Response.Status = 500
		return q.RenderError(err)
	}
	defer manager.Close()
	err = manager.DeleteQuotationById(id)

	q.RenderArgs["avatarUrl"] = q.Session["avatarUrl"]

	q.Render(userid, nickName)
	return q.Redirect((*Quotation).Index)
}
