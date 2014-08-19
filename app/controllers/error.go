package controllers

import (
	"github.com/revel/revel"
)

type ErrorPages struct {
	*revel.Controller
}

func (this ErrorPages) Page404() revel.Result {
	return this.Render()
}

func (this ErrorPages) Page500() revel.Result {
	return this.Render()
}
