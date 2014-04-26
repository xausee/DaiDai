package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/revel/revel"
)

type AncientPoem struct {
	*revel.Controller
}

func (ap *AncientPoem) Add() revel.Result {
	email := ap.Session["email"]
	nickName := ap.Session["nickName"]
	avatarUrl := ap.Session["avatarUrl"]

	return ap.Render(email, nickName, avatarUrl)
}

func (ap *AncientPoem) PostAdd(ancientPoem *models.AncientPoem) revel.Result {
	ap.Validation.Required(ancientPoem.Tag).Message("请选择一个标签")
	ap.Validation.Required(ancientPoem.Content).Message("摘录内容不能为空")
	ap.Validation.Required(ancientPoem.Author).Message("作者不能为空")

	fmt.Println("诗歌标题： ", ancientPoem.Title)
	fmt.Println("诗歌标签： ", ancientPoem.Tag)
	fmt.Println("诗歌类型： ", ancientPoem.Style)
	fmt.Println("诗歌内容： ", ancientPoem.Content)
	fmt.Println("作者： ", ancientPoem.Author)

	if ap.Validation.HasErrors() {
		ap.Validation.Keep()
		ap.FlashParams()
		return ap.Redirect((*AncientPoem).Add)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		ap.Response.Status = 500
		return ap.RenderError(err)
	}
	defer manager.Close()

	err = manager.AddAncientPeom(ancientPoem)
	if err != nil {
		ap.Flash.Error(err.Error())
		return ap.Redirect((*AncientPoem).Add)
	}

	return ap.Redirect((*App).Add)
}
