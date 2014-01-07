package controllers

import "github.com/robfig/revel"

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	email := c.Session["email"]
	nickName := c.Session["nickName"]
	return c.Render(email, nickName)
}
