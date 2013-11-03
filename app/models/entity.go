package models

type User struct { 
  Email    string 
  Nickname string 
  Gender      string
  Password []byte 
}

type MockUser struct { 
  Email           string 
  Nickname        string 
  Gender             string
  Password        string 
  ConfirmPassword string   
}