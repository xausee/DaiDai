package models

// 用户基本信息其它数据分开定义？
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
	NickName        string
	PenName         string
	Gender          string
	Password        string
	ConfirmPassword string
}

type LoginUser struct {
	Email    string
	NickName string
	Password string
}

type Commenter struct {
	Email     string
	NickName  string
	AvatarUrl string
}

type Fans struct {
	Email     string
	NickName  string
	AvatarUrl string
}

type Watch struct {
	Email     string
	NickName  string
	AvatarUrl string
}

type Comment struct {
	Id      string
	Author  Commenter
	Time    string
	Score   int
	Content string
}

type UserArticle struct {
	Id             string
	Title          string
	Tag            string
	Content        string
	CreateTime     string
	AuthorId       int
	AuthorNickName string
	Comments       []Comment
}

type ArticleInCollection struct {
	Id             string
	Title          string
	AuthorId       int
	AuthorNickName string
}

type RecommendArticle struct {
	Article ArticleInCollection
}

type User struct {
	Id                int
	Email             string
	NickName          string
	RealName          string
	PenName           string
	AvatarUrl         string
	Password          []byte
	Birth             string
	Gender            string
	FavoriteAuthor    string
	FavoriteBook      string
	Intrest           string
	Introduction      string
	Fans              []Fans
	Watch             []Watch
	Message           []Comment
	ArticleCollection []ArticleInCollection
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
