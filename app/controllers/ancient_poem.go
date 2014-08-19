package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/revel/revel"
)

type AncientPoem struct {
	*revel.Controller
}

func (this *AncientPoem) Add() revel.Result {
	email := this.Session["email"]
	nickName := this.Session["nickName"]
	avatarUrl := this.Session["avatarUrl"]

	return this.Render(email, nickName, avatarUrl)
}

func (this *AncientPoem) PostAdd(ancientPoem *models.AncientPoem) revel.Result {
	this.Validation.Required(ancientPoem.Tag).Message("请选择一个标签")
	this.Validation.Required(ancientPoem.Content).Message("摘录内容不能为空")
	this.Validation.Required(ancientPoem.Author).Message("作者不能为空")

	fmt.Println("诗歌标题： ", ancientPoem.Title)
	fmt.Println("诗歌标签： ", ancientPoem.Tag)
	fmt.Println("诗歌类型： ", ancientPoem.Style)
	fmt.Println("诗歌内容： ", ancientPoem.Content)
	fmt.Println("作者： ", ancientPoem.Author)

	if this.Validation.HasErrors() {
		this.Validation.Keep()
		this.FlashParams()
		return this.Redirect((*AncientPoem).Add)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()

	err = manager.AddAncientPeom(ancientPoem)
	if err != nil {
		this.Flash.Error(err.Error())
		return this.Redirect((*AncientPoem).Add)
	}

	return this.Redirect((*App).Add)
}
