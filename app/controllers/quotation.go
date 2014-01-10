package controllers

import (
	"ZhaiLuBaiKe/app/models"
	"fmt"
	"github.com/robfig/revel"
)

type Quotation struct {
	*revel.Controller
}

func (q *Quotation) Add() revel.Result {
	email := q.Session["email"]
	nickName := q.Session["nickName"]
	return q.Render(email, nickName)
}

func (q *Quotation) PostAdd(quotation *models.Quotation) revel.Result {
	q.Validation.Required(quotation.Content).Message("摘录内容不能为空")
	//q.Validation.Required(quotation.Original).Message("原文不能为空")
	q.Validation.Required(quotation.Author).Message("作者不能为空")

	fmt.Println(quotation.Content)
	fmt.Println(quotation.Author)

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

	return q.Redirect((*Account).RegisterSuccessful)
	//return q.Redirect((*Quotation).AddSuccessful)
}
