package models

type User struct { 
  Email    string 
  Nickname string 
  Password []byte 
}

type MockUser struct { 
  Email           string 
  Nickname        string 
  Password        string 
  ConfirmPassword string 
}