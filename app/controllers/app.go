package controllers

import (
	"SanWenJia/app/models"
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

	// 获取微小说最新15条数据
	c.RenderHintFictions(manager)

	return c.Render(email, nickName)
}

func (c App) Add() revel.Result {
	email := c.Session["email"]
	nickName := c.Session["nickName"]

	return c.Render(email, nickName)
}

func (c *App) RenderQuotations(manager *models.DbManager) error {
	quotations, err := manager.GetAllQuotation()
	c.RenderArgs["quotations"] = quotations

	return err
}

func (c *App) RenderWitticismQuotation(manager *models.DbManager) error {
	witticisms, err := manager.GetAllQuotation()
	c.RenderArgs["witticisms"] = witticisms

	return err
}

func (c *App) RenderAncientPoems(manager *models.DbManager) error {
	poems, err := manager.GetAllAncientPoem()
	c.RenderArgs["ancientPoems"] = poems

	return err
}

func (c *App) RenderModernPoems(manager *models.DbManager) error {
	poems, err := manager.GetAllModernPoem()
	c.RenderArgs["modernPoems"] = poems

	return err
}

func (c *App) RenderEssays(manager *models.DbManager) error {
	essays, err := manager.GetAllEssay()
	c.RenderArgs["essays"] = essays

	return err
}

func (c *App) RenderHintFictions(manager *models.DbManager) error {
	hintFictions, err := manager.GetAllHintFiction()
	c.RenderArgs["hintFictions"] = hintFictions

	return err
}
