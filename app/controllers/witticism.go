package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/revel/revel"
)

type Witticism struct {
	*revel.Controller
}

func (w *Witticism) Index() revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	witticisms, err := manager.GetAllWitticism()

	w.RenderArgs["userid"] = w.Session["userid"]
	w.RenderArgs["nickName"] = w.Session["nickName"]
	w.RenderArgs["witticisms"] = witticisms

	return w.Render()
}

func (w *Witticism) Add() revel.Result {
	w.RenderArgs["userid"] = w.Session["userid"]
	w.RenderArgs["nickName"] = w.Session["nickName"]

	return w.Render()
}

func (w *Witticism) Edit(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		w.Response.Status = 500
		return w.RenderError(err)
	}
	defer manager.Close()
	originalWitticism, _ := manager.GetWitticismById(id)

	w.RenderArgs["userid"] = w.Session["userid"]
	w.RenderArgs["nickName"] = w.Session["nickName"]
	w.RenderArgs["originalWitticism"] = originalWitticism

	return w.Render()
}

func (w *Witticism) PostAdd(witticism *models.Witticism) revel.Result {
	w.Validation.Required(witticism.Content).Message("慧语内容不能为空")
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

func (w *Witticism) PostEdit(originalWitticismID string, newWitticism *models.Witticism) revel.Result {
	w.Validation.Required(newWitticism.Content).Message("内容不能为空")
	w.Validation.Required(newWitticism.Author).Message("作者不能为空")

	fmt.Println("内容： ", newWitticism.Content)
	fmt.Println("作者： ", newWitticism.Author)

	if w.Validation.HasErrors() {
		w.Validation.Keep()
		w.FlashParams()
		fmt.Println("error in validation ")
		return w.Redirect((*Witticism).Edit)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		w.Response.Status = 500
		return w.RenderError(err)
	}
	defer manager.Close()

	err = manager.UpdateWitticism(originalWitticismID, newWitticism)
	if err != nil {
		w.Validation.Keep()
		w.FlashParams()
		w.Flash.Error(err.Error())
		fmt.Println("error in update Witticism ")
		return w.Redirect((*Witticism).Edit)
	}

	return w.Redirect((*Witticism).Index)
}

func (w *Witticism) Show(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		w.Response.Status = 500
		return w.RenderError(err)
	}
	defer manager.Close()
	witticism, _ := manager.GetWitticismById(id)
	fmt.Println("作者： ", witticism.Author)
	// if err != nil {
	// 	w.Flash.Error(err.Error())
	// 	//return w.Redirect((*Essay).Add)
	// }

	w.RenderArgs["userid"] = w.Session["userid"]
	w.RenderArgs["nickName"] = w.Session["nickName"]
	w.RenderArgs["witticism"] = witticism

	return w.Render()
}

func (w *Witticism) Delete(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		w.Response.Status = 500
		return w.RenderError(err)
	}
	defer manager.Close()
	err = manager.DeleteWitticismById(id)

	w.RenderArgs["userid"] = w.Session["userid"]
	w.RenderArgs["nickName"] = w.Session["nickName"]

	return w.Redirect((*Witticism).Index)
}
