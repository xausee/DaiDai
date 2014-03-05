package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/robfig/revel"
)

type Quotation struct {
	*revel.Controller
}

func (q *Quotation) Index() revel.Result {
	email := q.Session["email"]
	nickName := q.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	quotations, err := manager.GetAllQuotation()

	return q.Render(email, nickName, quotations)
}

func (q *Quotation) TypeIndex(tag string) revel.Result {
	email := q.Session["email"]
	nickName := q.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	quotations, err := manager.GetQuotationByTag(tag)

	return q.Render(email, nickName, quotations)
}

func (q *Quotation) Add() revel.Result {
	email := q.Session["email"]
	nickName := q.Session["nickName"]
	return q.Render(email, nickName)
}

func (q *Quotation) Edit(id string) revel.Result {
	email := q.Session["email"]
	nickName := q.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		q.Response.Status = 500
		return q.RenderError(err)
	}
	defer manager.Close()
	originalQuotation, _ := manager.GetQuotationById(id)

	return q.Render(email, nickName, originalQuotation)
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
	email := q.Session["email"]
	nickName := q.Session["nickName"]

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
	return q.Render(email, nickName, quotation)
}

func (q *Quotation) Delete(id string) revel.Result {
	email := q.Session["email"]
	nickName := q.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		q.Response.Status = 500
		return q.RenderError(err)
	}
	defer manager.Close()
	err = manager.DeleteQuotationById(id)

	q.Render(email, nickName)
	return q.Redirect((*Quotation).Index)
}
