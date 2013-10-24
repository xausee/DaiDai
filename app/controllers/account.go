package controllers

import ( 
  "github.com/robfig/revel" 
  "DaiDai/app/models" 
)

type Account struct { 
  *revel.Controller 
}

func (c *Account) Register() revel.Result { 
  return c.Render() 
}

func (c *Account) PostRegister(user *models.MockUser) revel.Result { 
  return c.Render() 
}