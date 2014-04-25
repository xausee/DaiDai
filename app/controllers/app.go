package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (this App) Index() revel.Result {
	userid := this.Session["userid"]
	//不再使用邮件地址登陆，用昵称替代掉
	nickName := this.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	// 获取摘抄最新15条数据
	this.RenderQuotations(manager)

	// 获取名人语录最新5条数据
	this.RenderWitticism(manager)

	// 获取古诗词最新15条数据
	this.RenderAncientPoems(manager)

	// 获取现代诗最新15条数据
	this.RenderModernPoems(manager)

	// 获取散文最新15条数据
	this.RenderEssays(manager)

	// 获取微小说最新15条数据
	this.RenderHintFictions(manager)

	// 获取所有的推荐文章（15）篇
	this.RenderRecommendArticles(manager)

	return this.Render(userid, nickName)
}

func (this App) Add() revel.Result {
	userid := this.Session["userid"]
	nickName := this.Session["nickName"]

	return this.Render(userid, nickName)
}

func (this App) PostSearch(keywords string) revel.Result {
	userid := this.Session["userid"]
	nickName := this.Session["nickName"]

	fmt.Println(keywords)
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	results, err := manager.Search(keywords)

	noResults := false
	if len(results) == 0 {
		noResults = true
	}

	this.RenderArgs["userid"] = userid
	this.RenderArgs["nickName"] = nickName
	this.RenderArgs["results"] = results
	this.RenderArgs["noResults"] = noResults
	this.RenderArgs["keywords"] = keywords

	return this.Render()
}

func (this *App) RenderQuotations(manager *models.DbManager) error {
	qs, err := manager.GetAllQuotation()

	// 倒序处理
	var count = len(qs)
	for i := 0; i < count/2; i++ {
		qs[i], qs[count-i-1] = qs[count-i-1], qs[i]
	}

	// 只取前15条摘录用于首页的显示
	var quotations []models.Quotation
	var more bool
	if count > models.ArticlesInHomePanel {
		more = true
		quotations = qs[:models.ArticlesInHomePanel]
	} else {
		more = false
		quotations = qs
	}

	this.RenderArgs["quotations"] = quotations
	this.RenderArgs["moreQuotation"] = more
	return err
}

func (this *App) RenderWitticism(manager *models.DbManager) error {
	ws, err := manager.GetAllWitticism()

	// 倒序处理
	var count = len(ws)
	for i := 0; i < count/2; i++ {
		ws[i], ws[count-i-1] = ws[count-i-1], ws[i]
	}

	// 只取前5条名人语录用于首页的显示
	var witticisms []models.Witticism
	var more bool
	if count > 5 {
		more = true
		witticisms = ws[:5]
	} else {
		more = false
		witticisms = ws
	}

	// 轮播4条慧语
	show := false
	if count >= 4 {
		show = true
		this.RenderArgs["witticism1"] = witticisms[count-1]
		this.RenderArgs["witticism2"] = witticisms[count-2]
		this.RenderArgs["witticism3"] = witticisms[count-3]
		this.RenderArgs["witticism4"] = witticisms[count-4]
	}

	this.RenderArgs["witticisms"] = witticisms
	this.RenderArgs["showWitticism"] = show
	this.RenderArgs["moreMitticism"] = more
	return err
}

func (this *App) RenderAncientPoems(manager *models.DbManager) error {
	poems, err := manager.GetAllAncientPoem()
	this.RenderArgs["ancientPoems"] = poems

	return err
}

func (this *App) RenderModernPoems(manager *models.DbManager) error {
	ps, err := manager.GetAllModernPoem()

	// 倒序处理
	var count = len(ps)
	for i := 0; i < count/2; i++ {
		ps[i], ps[count-i-1] = ps[count-i-1], ps[i]
	}

	// 只取前15篇现代诗用于首页的显示
	var poems []models.ModernPoem
	var more bool
	if count > models.ArticlesInHomePanel {
		more = true
		poems = ps[:models.ArticlesInHomePanel]
	} else {
		more = false
		poems = ps
	}

	this.RenderArgs["modernPoems"] = poems
	this.RenderArgs["moreModernPoem"] = more
	return err
}

func (this *App) RenderEssays(manager *models.DbManager) error {
	es, err := manager.GetAllEssay()

	// 倒序处理
	var count = len(es)
	for i := 0; i < count/2; i++ {
		es[i], es[count-i-1] = es[count-i-1], es[i]
	}

	// 只取前15篇散文用于首页的显示
	var essays []models.Essay
	var more bool
	if count > models.ArticlesInHomePanel {
		more = true
		essays = es[:models.ArticlesInHomePanel]
	} else {
		more = false
		essays = es
	}

	this.RenderArgs["essays"] = essays
	this.RenderArgs["moreEssay"] = more
	return err
}

func (this *App) RenderHintFictions(manager *models.DbManager) error {
	hs, err := manager.GetAllHintFiction()

	// 倒序处理
	var count = len(hs)
	for i := 0; i < count/2; i++ {
		hs[i], hs[count-i-1] = hs[count-i-1], hs[i]
	}

	// 只取前15篇现代诗用于首页的显示
	var hintFictions []models.HintFiction
	var more bool
	if count > models.ArticlesInHomePanel {
		more = true
		hintFictions = hs[:models.ArticlesInHomePanel]
	} else {
		more = false
		hintFictions = hs
	}

	this.RenderArgs["hintFictions"] = hintFictions
	this.RenderArgs["moreHintFiction"] = more
	return err
}

func (this *App) RenderRecommendArticles(manager *models.DbManager) error {
	as, err := manager.GetAllRecommendArticle()

	// 倒序处理
	var count = len(as)
	for i := 0; i < count/2; i++ {
		as[i], as[count-i-1] = as[count-i-1], as[i]
	}

	this.RenderArgs["allRecommendArticles"] = as
	return err
}

func (this App) AboutUs() revel.Result {
	userid := this.Session["userid"]
	nickName := this.Session["nickName"]

	return this.Render(userid, nickName)
}

func (this App) Donate() revel.Result {
	userid := this.Session["userid"]
	nickName := this.Session["nickName"]

	return this.Render(userid, nickName)
}
