package controllers

import (
	"SanWenJia/app/models"
	"fmt"
	"github.com/robfig/revel"
	"strconv"
)

type HintFiction struct {
	*revel.Controller
}

func (hf *HintFiction) Index() revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	hintFictions, err := manager.GetAllHintFiction()
	count := len(hintFictions)

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

	hintFictionsOnOnePage := []models.HintFiction{}
	if count > models.ArticlesInSinglePage {
		hintFictionsOnOnePage = hintFictions[(count - models.ArticlesInSinglePage):]
	} else {
		hintFictionsOnOnePage = hintFictions
	}

	hf.RenderArgs["id"] = hf.Session["id"]
	hf.RenderArgs["email"] = hf.Session["email"]
	hf.RenderArgs["nickName"] = hf.Session["nickName"]
	hf.RenderArgs["allHintFinctions"] = hintFictions
	hf.RenderArgs["hintFictionsOnOnePage"] = hintFictionsOnOnePage
	hf.RenderArgs["pageCount"] = pageCount
	hf.RenderArgs["pages"] = pages

	return hf.Render()
}

func (hf *HintFiction) TypeIndex(tag string) revel.Result {
	email := hf.Session["email"]
	nickName := hf.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	hintFictions, err := manager.GetHintFictionByTag(tag)

	return hf.Render(email, nickName, hintFictions)
}

func (hf *HintFiction) Add() revel.Result {
	email := hf.Session["email"]
	nickName := hf.Session["nickName"]
	return hf.Render(email, nickName)
}

func (hf *HintFiction) Edit(id string) revel.Result {
	email := hf.Session["email"]
	nickName := hf.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		hf.Response.Status = 500
		return hf.RenderError(err)
	}
	defer manager.Close()
	oringinalHintFiction, _ := manager.GetHintFictionById(id)

	return hf.Render(email, nickName, oringinalHintFiction)
}

func (hf *HintFiction) PostAdd(hintFiction *models.HintFiction) revel.Result {
	hf.Validation.Required(hintFiction.Tag).Message("请选择一个标签")
	hf.Validation.Required(hintFiction.Content).Message("内容不能为空")
	hf.Validation.Required(hintFiction.Author).Message("作者不能为空")

	fmt.Println("微小说标签： ", hintFiction.Tag)
	fmt.Println("微小说标题： ", hintFiction.Title)
	fmt.Println("微小说内容： ", hintFiction.Content)
	fmt.Println("微小说作者： ", hintFiction.Author)

	if hf.Validation.HasErrors() {
		hf.Validation.Keep()
		hf.FlashParams()
		return hf.Redirect((*HintFiction).Add)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		hf.Response.Status = 500
		return hf.RenderError(err)
	}
	defer manager.Close()

	err = manager.AddHintFiction(hintFiction)
	if err != nil {
		hf.Flash.Error(err.Error())
		return hf.Redirect((*HintFiction).Add)
	}

	return hf.Redirect((*App).Add)
}

func (hf *HintFiction) PostEdit(originalHintFictionID string, newHintFiction *models.HintFiction) revel.Result {
	hf.Validation.Required(newHintFiction.Tag).Message("请选择一个标签")
	hf.Validation.Required(newHintFiction.Content).Message("内容不能为空")
	hf.Validation.Required(newHintFiction.Author).Message("作者不能为空")

	fmt.Println("微小说标签： ", newHintFiction.Tag)
	fmt.Println("微小说标题： ", newHintFiction.Title)
	fmt.Println("微小说内容： ", newHintFiction.Content)
	fmt.Println("微小说作者： ", newHintFiction.Author)

	if hf.Validation.HasErrors() {
		hf.Validation.Keep()
		hf.FlashParams()
		return hf.Redirect((*HintFiction).Edit)
	}

	manager, err := models.NewDbManager()
	if err != nil {
		hf.Response.Status = 500
		return hf.RenderError(err)
	}
	defer manager.Close()

	err = manager.UpdateHintFiction(originalHintFictionID, newHintFiction)
	if err != nil {
		hf.Flash.Error(err.Error())
		return hf.Redirect((*HintFiction).Edit)
	}

	return hf.Redirect((*HintFiction).Index)
}

func (hf *HintFiction) Show(id string) revel.Result {
	email := hf.Session["email"]
	nickName := hf.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		hf.Response.Status = 500
		return hf.RenderError(err)
	}
	defer manager.Close()
	hintFiction, _ := manager.GetHintFictionById(id)
	// if err != nil {
	// 	hf.Flash.Error(err.Error())
	// 	//return hf.Redirect((*Essay).Add)
	// }
	return hf.Render(email, nickName, hintFiction)
}

func (hf *HintFiction) PageList(pageNumber string) revel.Result {
	manager, err := models.NewDbManager()
	if err != nil {
		fmt.Println("链接数据库失败")
	}
	defer manager.Close()

	allHintfictions, err := manager.GetAllHintFiction()
	count := len(allHintfictions)

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

	hintFictionsOnOnePage := []models.HintFiction{}
	if count >= iPageNumber*models.ArticlesInSinglePage {
		hintFictionsOnOnePage = allHintfictions[(iPageNumber-1)*models.ArticlesInSinglePage : iPageNumber*models.ArticlesInSinglePage]
	} else if (iPageNumber-1)*models.ArticlesInSinglePage < count && count < iPageNumber*models.ArticlesInSinglePage {
		hintFictionsOnOnePage = allHintfictions[(iPageNumber-1)*models.ArticlesInSinglePage:]
	} else if (iPageNumber-1)*models.ArticlesInSinglePage > count {
		fmt.Println("已超过最大页码")
	}
	fmt.Println("pageCount:", pageCount)
	fmt.Println("pageNumber:", pageNumber)

	hf.RenderArgs["allHintfictions"] = allHintfictions
	hf.RenderArgs["hintFictionsOnOnePage"] = hintFictionsOnOnePage
	hf.RenderArgs["pageCount"] = pageCount
	hf.RenderArgs["pageNumber"] = pageNumber

	return hf.Render()
}

func (hf *HintFiction) Delete(id string) revel.Result {
	email := hf.Session["email"]
	nickName := hf.Session["nickName"]

	manager, err := models.NewDbManager()
	if err != nil {
		hf.Response.Status = 500
		return hf.RenderError(err)
	}
	defer manager.Close()
	err = manager.DeleteHintFictionById(id)

	hf.Render(email, nickName)
	return hf.Redirect((*HintFiction).Index)
}
