package models

type User struct {
	Email    string
	Nickname string
	Gender   string
	Password []byte
}

type MockUser struct {
	Email           string
	Nickname        string
	Gender          string
	Password        string
	ConfirmPassword string
}

type LoginUser struct {
	Email    string
	Nickname string
	Password string
}

type Quotation struct {
	Tag      string
	Content  string
	Original string
	Author   string
}
