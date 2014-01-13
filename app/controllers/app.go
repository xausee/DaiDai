package controllers

import (
	"ZhaiLuBaiKe/app/models"
	"fmt"
	"github.com/robfig/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	email := c.Session["email"]
	nickName := c.Session["nickName"]
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()
	qs, err := manager.GetAllQuotation()
	c.RenderArgs["tag"] = qs[3].Tag
	c.RenderArgs["content"] = qs[3].Content
	c.RenderArgs["original"] = qs[3].Original
	c.RenderArgs["author"] = qs[3].Author
	return c.Render(email, nickName)
}
