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
	count := len(qs)
	fmt.Println("总共存在", count, "条数据")
	if count > 4 {
		c.RenderArgs["quotations1"] = qs[count-4]
		c.RenderArgs["quotations2"] = qs[count-3]
		c.RenderArgs["quotations3"] = qs[count-2]
		c.RenderArgs["quotations4"] = qs[count-1]
	}
	return c.Render(email, nickName)
}
