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

	// 获取摘抄最新15条数据
	c.RenderQuotations(manager)

	// 获取名人语录最新4条数据
	c.RenderWitticismQuotation(manager)

	// 获取古诗词最新15条数据
	c.RenderAncientPoems(manager)

	// 获取现代诗最新15条数据
	c.RenderModernPoems(manager)

	// 获取散文最新15条数据
	c.RenderEssays(manager)

	return c.Render(email, nickName)
}

func (c App) Add() revel.Result {
	email := c.Session["email"]
	nickName := c.Session["nickName"]

	return c.Render(email, nickName)
}

func (c *App) RenderQuotations(manager *models.DbManager) error {
	quotations, err := manager.GetAllQuotation()
	count := len(quotations)
	if count >= 4 {
		c.RenderArgs["quotation1"] = quotations[count-1]
		c.RenderArgs["quotation2"] = quotations[count-2]
		c.RenderArgs["quotation3"] = quotations[count-3]
		c.RenderArgs["quotation4"] = quotations[count-4]
	}
	return err
}

func (c *App) RenderWitticismQuotation(manager *models.DbManager) error {
	witticisms, err := manager.GetAllQuotation()
	count := len(witticisms)
	if count >= 4 {
		c.RenderArgs["witticism1"] = witticisms[count-1]
		c.RenderArgs["witticism2"] = witticisms[count-2]
		c.RenderArgs["witticism3"] = witticisms[count-3]
		c.RenderArgs["witticism4"] = witticisms[count-4]
	}
	return err
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

func (c *App) RenderEssays(manager *models.DbManager) error {
	essays, err := manager.GetAllEssay()
	count := len(essays)
	if count >= 1 {
		c.RenderArgs["essay1"] = essays[count-1]
		// c.RenderArgs["essay2"] = essays[count-2]
		// c.RenderArgs["essay3"] = essays[count-3]
		// c.RenderArgs["essay4"] = essays[count-4]
		// c.RenderArgs["essay5"] = essays[count-5]
	}
	return err
}
