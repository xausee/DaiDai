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
	Tag           string
	Content       string
	Original      string
	OriginalTitle string
	Author        string
}

type AncientPoem struct {
	Title   string
	Style   string
	Tag     string
	Content string
	Author  string
}

type ModernPoem struct {
	Title   string
	Tag     string
	Content string
	Author  string
}

type Essay struct {
	Title   string
	Tag     string
	Content string
	Author  string
}
