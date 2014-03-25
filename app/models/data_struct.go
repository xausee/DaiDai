package models

import (
	"time"
)

// type User struct {
// 	Id        int
// 	Email     string
// 	Nickname  string
// 	Gender    string
// 	Password  []byte
// 	AvatarUrl string
// 	PenName   string
// 	Birth     time.Time
// }

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

type Commenter struct {
	Email     string
	Nickname  string
	AvatarUrl string
}

type Comment struct {
	Id     string
	Author Commenter
	Time   time.Time
	Score  int
	Cotent string
}

type UserArticle struct {
	Id         string
	Title      string
	Tag        string
	Content    string
	CreateTime string
	AuthorId   int
	Comments   []Comment
}

type User struct {
	Id                int
	Email             string
	Nickname          string
	AvatarUrl         string
	Password          []byte
	PenName           string
	Birth             time.Time
	Gender            string
	FavoriteAuthor    string
	FavoriteBook      string
	Intrest           string
	Introduction      string
	Fans              string
	Watch             string
	Message           string
	ArticleCollection []UserArticle
	Articles          []UserArticle
}

type Quotation struct {
	Id            string
	Tag           string
	Content       string
	Original      string
	OriginalTitle string
	Author        string
}

type Witticism struct {
	Id      string
	Content string
	Author  string
}

type AncientPoem struct {
	Title   string
	Style   string
	Tag     string
	Content string
	Author  string
}

type ModernPoem struct {
	Id      string
	Title   string
	Tag     string
	Content string
	Author  string
}

type Essay struct {
	Id      string
	Title   string
	Tag     string
	Content string
	Author  string
}

type HintFiction struct {
	Id      string
	Title   string
	Tag     string
	Content string
	Author  string
}
