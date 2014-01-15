package controllers

import (
	"ZhaiLuBaiKe/app/models"
	"fmt"
	"github.com/robfig/revel"
	//"strings"
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

	allQuotation, err := manager.GetAllQuotation()
	count := len(allQuotation)
	fmt.Println("总共存在", count, "条数据")
	if count >= 2 {
		//allQuotation[count-1].Content = strings.Replace(allQuotation[count-1].Content, "\r\n", "<br>", -1)
		c.RenderArgs["quotations1"] = allQuotation[count-1]
		c.RenderArgs["quotations2"] = allQuotation[count-2]
	}

	allWitticism, _ := manager.GetAllWitticismQuotation()
	count = len(allWitticism)
	fmt.Println("总共存在", count, "名人语录")
	if count >= 4 {
		c.RenderArgs["witticism1"] = allWitticism[count-1]
		c.RenderArgs["witticism2"] = allWitticism[count-2]
		c.RenderArgs["witticism3"] = allWitticism[count-3]
		c.RenderArgs["witticism4"] = allWitticism[count-4]
	}

	return c.Render(email, nickName)
}

func (c App) Add() revel.Result {
	email := c.Session["email"]
	nickName := c.Session["nickName"]

	return c.Render(email, nickName)
}
