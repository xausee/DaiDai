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

func (c *Account) RegisterSuccessful() revel.Result { 
  return c.Render() 
}

func (c *Account) PostRegister(user *models.MockUser) revel.Result { 
  c.Validation.Email(user.Email).Message("电子邮件格式无效") 
  c.Validation.Required(user.Nickname).Message("用户昵称不能为空") 
  c.Validation.Required(user.Password).Message("密码不能为空") 
  c.Validation.Required(user.ConfirmPassword == user.Password).Message("两次输入的密码不一致")

  if c.Validation.HasErrors() { 
    c.Validation.Keep()
    c.FlashParams() 
    return c.Redirect((*Account).Register) 
  }

  dal, err := models.NewDal() 
  if err != nil { 
    c.Response.Status = 500 
    return c.RenderError(err) 
  } 
  defer dal.Close()

  err = dal.RegisterUser(user) 
  if err != nil { 
    c.Flash.Error(err.Error())         
    return c.Redirect((*Account).Register) 
  }

  return c.Redirect((*Account).RegisterSuccessful) 
}