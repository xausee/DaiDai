package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/robfig/revel"
)

type Witticism struct {
	*revel.Controller
}

func (w *Witticism) Index() revel.Result {
	email := w.Session["email"]
	nickName := w.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	Witticisms, err := manager.GetAllWitticism()

	return w.Render(email, nickName, Witticisms)
}

func (w *Witticism) Add() revel.Result {
	email := w.Session["email"]
	nickName := w.Session["nickName"]
	return w.Render(email, nickName)
}

func (w *Witticism) PostAdd(witticism *models.Witticism) revel.Result {
	w.Validation.Required(witticism.Content).Message("摘录内容不能为空")
	w.Validation.Required(witticism.Author).Message("作者不能为空")

	fmt.Println("慧语被容： ", witticism.Content)
	fmt.Println("作者： ", witticism.Author)

	if w.Validation.HasErrors() {
		w.Validation.Keep()
		w.FlashParams()
		return w.Redirect((*Witticism).Add)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		w.Response.Status = 500
		return w.RenderError(err)
	}
	defer manager.Close()

	err = manager.AddWitticism(witticism)
	if err != nil {
		// q.Validation.Keep()
		// q.FlashParams()
		w.Flash.Error(err.Error())
		return w.Redirect((*Witticism).Add)
	}

	return w.Redirect((*App).Add)
}

func (w *Witticism) Show(id string) revel.Result {
	email := w.Session["email"]
	nickName := w.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		w.Response.Status = 500
		return w.RenderError(err)
	}
	defer manager.Close()
	Witticism, _ := manager.GetWitticismById(id)
	fmt.Println("作者： ", Witticism)
	// if err != nil {
	// 	w.Flash.Error(err.Error())
	// 	//return w.Redirect((*Essay).Add)
	// }
	return w.Render(email, nickName, Witticism)
}
