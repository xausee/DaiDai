package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/revel/revel"
	"strconv"
)

type ModernPoem struct {
	*revel.Controller
}

func (this *ModernPoem) Index() revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	allPoems, err := manager.GetAllModernPoem()

	// 倒序处理
	count := len(allPoems)
	for i := 0; i < count/2; i++ {
		allPoems[i], allPoems[count-i-1] = allPoems[count-i-1], allPoems[i]
	}

	var pageCount int
	if (count % models.ArticlesInSinglePage) == 0 {
		pageCount = count / models.ArticlesInSinglePage
	} else {
		pageCount = count/models.ArticlesInSinglePage + 1
	}

	pages := make([]int, 0)
	for i := 1; i <= pageCount; i++ {
		pages = append(pages, i)
	}

	poemsOnOnePage := []models.ModernPoem{}
	if count > models.ArticlesInSinglePage {
		poemsOnOnePage = allPoems[:models.ArticlesInSinglePage]
	} else {
		poemsOnOnePage = allPoems
	}

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["allPoems"] = allPoems
	this.RenderArgs["poemsOnOnePage"] = poemsOnOnePage
	this.RenderArgs["pageCount"] = pageCount
	this.RenderArgs["pages"] = pages

	return this.Render()
}

func (this *ModernPoem) TypeIndex(tag string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	allPoems, err := manager.GetModernPoemByTag(tag)

	// 倒序处理
	count := len(allPoems)
	for i := 0; i < count/2; i++ {
		allPoems[i], allPoems[count-i-1] = allPoems[count-i-1], allPoems[i]
	}

	var pageCount int
	if (count % models.ArticlesInSinglePage) == 0 {
		pageCount = count / models.ArticlesInSinglePage
	} else {
		pageCount = count/models.ArticlesInSinglePage + 1
	}

	pageSlice := make([]int, 0)
	for i := 1; i <= pageCount; i++ {
		pageSlice = append(pageSlice, i)
	}

	poemsOnOnePage := []models.ModernPoem{}
	if count > models.ArticlesInSinglePage {
		poemsOnOnePage = allPoems[:models.ArticlesInSinglePage]
	} else {
		poemsOnOnePage = allPoems
	}

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["allPoems"] = allPoems
	this.RenderArgs["poemsOnOnePage"] = poemsOnOnePage
	this.RenderArgs["pageCount"] = pageCount
	this.RenderArgs["type"] = tag
	this.RenderArgs["pageSlice"] = pageSlice

	return this.Render()
}

func (this *ModernPoem) Add() revel.Result {
	userid := this.Session["userid"]
	nickName := this.Session["nickName"]
	return this.Render(userid, nickName)
}

func (this *ModernPoem) Edit(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()
	oringinalModernPoem, _ := manager.GetModernPoemById(id)

	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["oringinalModernPoem"] = oringinalModernPoem

	return this.Render()
}

func (this *ModernPoem) PostAdd(modernPoem *models.ModernPoem) revel.Result {
	this.Validation.Required(modernPoem.Tag).Message("请选择一个标签")
	this.Validation.Required(modernPoem.Content).Message("摘录内容不能为空")
	this.Validation.Required(modernPoem.Author).Message("作者不能为空")

	fmt.Println("诗歌标签： ", modernPoem.Tag)
	fmt.Println("诗歌标题： ", modernPoem.Title)
	fmt.Println("诗歌内容： ", modernPoem.Content)
	fmt.Println("作者： ", modernPoem.Author)

	if this.Validation.HasErrors() {
		this.Validation.Keep()
		this.FlashParams()
		return this.Redirect((*ModernPoem).Add)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()

	err = manager.AddModernPeom(modernPoem)
	if err != nil {
		this.Flash.Error(err.Error())
		return this.Redirect((*ModernPoem).Add)
	}

	// 返回到管理员首页
	return this.Redirect((*Admin).Index)
}

func (this *ModernPoem) PostEdit(originalModernPoemID string, newModernPoem *models.ModernPoem) revel.Result {
	this.Validation.Required(newModernPoem.Tag).Message("请选择一个标签")
	this.Validation.Required(newModernPoem.Content).Message("摘录内容不能为空")
	this.Validation.Required(newModernPoem.Author).Message("作者不能为空")

	fmt.Println("诗歌标签： ", newModernPoem.Tag)
	fmt.Println("诗歌标题： ", newModernPoem.Title)
	fmt.Println("诗歌内容： ", newModernPoem.Content)
	fmt.Println("作者： ", newModernPoem.Author)

	if this.Validation.HasErrors() {
		this.Validation.Keep()
		this.FlashParams()
		return this.Redirect((*ModernPoem).Edit)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()

	err = manager.UpdateModernPeom(originalModernPoemID, newModernPoem)
	if err != nil {
		this.Flash.Error(err.Error())
		return this.Redirect((*ModernPoem).Edit)
	}

	return this.Redirect((*ModernPoem).Index)
}

func (this *ModernPoem) Show(id string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()
	modernPoem, _ := manager.GetModernPoemById(id)
	// if err != nil {
	// 	this.Flash.Error(err.Error())
	// 	//return this.Redirect((*Essay).Add)
	// }

	// session 里保存是字符串，所以需要类型转换
	var admin = false
	var role models.Role
	roleStr, isExists := this.Session[models.CSessionRole]

	if isExists {
		value, _ := strconv.Atoi(roleStr)
		role = models.Role(value)

		if role == models.AdminRole {
			admin = true
		}
	}

	this.RenderArgs["admin"] = admin
	this.RenderArgs["userid"] = this.Session["userid"]
	this.RenderArgs["nickName"] = this.Session["nickName"]
	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]
	this.RenderArgs["modernPoem"] = modernPoem

	return this.Render()
}

func (this *ModernPoem) PageList(pageNumber string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	allPoems, err := manager.GetAllModernPoem()

	// 倒序处理
	count := len(allPoems)
	for i := 0; i < count/2; i++ {
		allPoems[i], allPoems[count-i-1] = allPoems[count-i-1], allPoems[i]
	}

	var pageCount int
	if (count % models.ArticlesInSinglePage) == 0 {
		pageCount = count / models.ArticlesInSinglePage
	} else {
		pageCount = count/models.ArticlesInSinglePage + 1
	}

	var iPageNumber int
	iPageNumber, err = strconv.Atoi(pageNumber)
	if err != nil {
		fmt.Println(err)
	}

	poemsOnOnePage := []models.ModernPoem{}
	if count >= iPageNumber*models.ArticlesInSinglePage {
		poemsOnOnePage = allPoems[(iPageNumber-1)*models.ArticlesInSinglePage : iPageNumber*models.ArticlesInSinglePage]
	} else if (iPageNumber-1)*models.ArticlesInSinglePage < count && count < iPageNumber*models.ArticlesInSinglePage {
		poemsOnOnePage = allPoems[(iPageNumber-1)*models.ArticlesInSinglePage:]
	} else if (iPageNumber-1)*models.ArticlesInSinglePage > count {
		fmt.Println("已超过最大页码")
	}
	fmt.Println("pageCount:", pageCount)
	fmt.Println("pageNumber:", pageNumber)

	this.RenderArgs["allPoems"] = allPoems
	this.RenderArgs["poemsOnOnePage"] = poemsOnOnePage
	this.RenderArgs["pageCount"] = pageCount
	this.RenderArgs["pageNumber"] = pageNumber

	return this.Render()
}

func (this *ModernPoem) PageListWithTag(uPageNumber string, tag string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	allPoems, err := manager.GetModernPoemByTag(tag)

	// 倒序处理
	count := len(allPoems)
	for i := 0; i < count/2; i++ {
		allPoems[i], allPoems[count-i-1] = allPoems[count-i-1], allPoems[i]
	}

	var pageCount int
	if (count % models.ArticlesInSinglePage) == 0 {
		pageCount = count / models.ArticlesInSinglePage
	} else {
		pageCount = count/models.ArticlesInSinglePage + 1
	}

	var iPageNumber int
	iPageNumber, err = strconv.Atoi(uPageNumber)
	if err != nil {
		fmt.Println(err)
	}

	poemsOnOnePage := []models.ModernPoem{}
	if count >= iPageNumber*models.ArticlesInSinglePage {
		poemsOnOnePage = allPoems[(iPageNumber-1)*models.ArticlesInSinglePage : iPageNumber*models.ArticlesInSinglePage]
	} else if (iPageNumber-1)*models.ArticlesInSinglePage < count && count < iPageNumber*models.ArticlesInSinglePage {
		poemsOnOnePage = allPoems[(iPageNumber-1)*models.ArticlesInSinglePage:]
	} else if (iPageNumber-1)*models.ArticlesInSinglePage > count {
		fmt.Println("已超过最大页码")
	}
	fmt.Println("pageCount:", pageCount)
	fmt.Println("uPageNumber:", uPageNumber)

	this.RenderArgs["allPoems"] = allPoems
	this.RenderArgs["poemsOnOnePage"] = poemsOnOnePage
	this.RenderArgs["pageCount"] = pageCount
	this.RenderArgs["uPageNumber"] = uPageNumber

	return this.Render()
}

func (this *ModernPoem) Delete(id string) revel.Result {
	userid := this.Session["userid"]
	nickName := this.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		this.Response.Status = 500
		return this.RenderError(err)
	}
	defer manager.Close()
	err = manager.DeleteModernPoemById(id)

	this.RenderArgs["avatarUrl"] = this.Session["avatarUrl"]

	this.Render(userid, nickName)
	return this.Redirect((*ModernPoem).Index)
}
