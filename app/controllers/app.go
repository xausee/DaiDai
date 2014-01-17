package controllers

import (
	"ZhaiLuBaiKe/app/models"
	"fmt"
	"github.com/robfig/revel"
	"strings"
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
	aqCount := len(allQuotation)
	fmt.Println("总共存在", aqCount, "条数据")
	if aqCount >= 2 {
		allQuotation[aqCount-1].Content = strings.Replace(allQuotation[aqCount-1].Content, "\r\n", "<br>", -1)
		c.RenderArgs["quotations1"] = allQuotation[aqCount-1]
		c.RenderArgs["quotations2"] = allQuotation[aqCount-2]
	}

	// 获取名人语录最新4条数据
	allWitticism, _ := manager.GetAllWitticismQuotation()
	awCount := len(allWitticism)
	fmt.Println("总共存在", awCount, "条名人语录")
	if awCount >= 4 {
		c.RenderArgs["witticism1"] = allWitticism[awCount-1]
		c.RenderArgs["witticism2"] = allWitticism[awCount-2]
		c.RenderArgs["witticism3"] = allWitticism[awCount-3]
		c.RenderArgs["witticism4"] = allWitticism[awCount-4]
	}

	// 获取古诗词最新15条数据
	c.RenderAncientPoems(manager)

	// 获取现代诗最新15条数据
	// mps, _ := manager.GetAllModernPoem()
	// mpCount := len(mps)
	// if mpCount >= 1 {
	// 	c.RenderArgs["modernPoem1"] = mps[mpCount-1]
	// }

	c.RenderModernPoems(manager)

	// 获取散文最新1条数据
	es, _ := manager.GetAllEssay()
	eCount := len(es)
	if eCount >= 1 {
		c.RenderArgs["essay1"] = es[eCount-1]
	}

	return c.Render(email, nickName)
}

func (c App) Add() revel.Result {
	email := c.Session["email"]
	nickName := c.Session["nickName"]

	return c.Render(email, nickName)
}

func (c *App) RenderAncientPoems(manager *models.DbManager) error {
	poems, err := manager.GetAllAncientPoem()
	count := len(poems)
	if count >= 4 {
		c.RenderArgs["ancientPoem1"] = poems[count-1]
		c.RenderArgs["ancientPoem2"] = poems[count-2]
		c.RenderArgs["ancientPoem3"] = poems[count-3]
		c.RenderArgs["ancientPoem4"] = poems[count-4]
		c.RenderArgs["ancientPoem5"] = poems[count-5]
	}
	return err
}

func (c *App) RenderModernPoems(manager *models.DbManager) error {
	poems, err := manager.GetAllModernPoem()
	count := len(poems)
	if count >= 1 {
		c.RenderArgs["modernPoem1"] = poems[count-1]
		// c.RenderArgs["modernPoem2"] = poems[count-2]
		// c.RenderArgs["modernPoem3"] = poems[count-3]
		// c.RenderArgs["modernPoem4"] = poems[count-4]
		// c.RenderArgs["modernPoem5"] = poems[count-5]
	}
	return err
}
