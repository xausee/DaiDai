package models

type User struct { 
  Email    string 
  Nickname string 
  Sex      string
  Password []byte 
}

type MockUser struct { 
  Email           string 
  Nickname        string 
  Sex             string
  Password        string 
  ConfirmPassword string   
}