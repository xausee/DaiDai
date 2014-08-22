package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
	"time"
)

func (manager *DbManager) GetUserById(userid int) (userInfo User, err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	err = uc.Find(bson.M{"id": userid}).One(&userInfo)
	if err != nil {
		fmt.Println("查询用户信息失败")
	}

	return userInfo, err
}

func (manager *DbManager) GetUserByNickName(nickName string) (userInfo User, err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	err = uc.Find(bson.M{"nickname": nickName}).One(&userInfo)
	if err != nil {
		fmt.Println("查询用户信息失败")
	}

	return userInfo, err
}
func (manager *DbManager) GetAllUser() (userInfo []User, err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	err = uc.Find(nil).All(&userInfo)
	if err != nil {
		fmt.Println("查询用户信息失败")
	}

	return userInfo, err
}

func (manager *DbManager) SearchUser(keywords string) (users []User, err error) {
	type Items map[string]string

	uc := manager.session.DB(DbName).C(UserCollection)
	err = uc.Find(bson.M{"nickname": Items{"$regex": keywords}}).All(&users)

	return users, err
}

func (manager *DbManager) DeleteUserByNickName(nickName string) (err error) {
	type Items map[string]string

	uc := manager.session.DB(DbName).C(UserCollection)
	err = uc.Remove(bson.M{"nickname": nickName})

	return err
}

func (manager *DbManager) VerifyPasswordByNickName(nickName, password string) (user *User, err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	i, _ := uc.Find(bson.M{"nickname": nickName}).Count()
	if i == 0 {
		fmt.Println("该用户不存在")
		err = errors.New("该用户不存在")
		return
	}

	uc.Find(bson.M{"nickname": nickName}).One(&user)

	if user.Password == nil {
		err = errors.New("获取密码错误")
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		fmt.Println("密码不正确,请重新输入原始密码")
		err = errors.New("密码不正确,请重新输入原始密码")
	}
	return
}

func (manager *DbManager) UpdateUserPasswordById(userid int, newPassword string) (err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	var oldUserInfo User
	err = uc.Find(bson.M{"id": userid}).One(&oldUserInfo)
	if err != nil {
		fmt.Println("查询用户信息失败")
	}

	// 仅仅修改密码字段
	tempInfo := oldUserInfo
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	tempInfo.Password = encryptedPassword

	err = uc.Update(oldUserInfo, tempInfo)
	if err != nil {
		fmt.Println("修改密码失败")
	}

	return err
}

func (manager *DbManager) UpdateUserInfoById(userid int, newUserInfo User) (err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	var oldUserInfo User
	err = uc.Find(bson.M{"id": userid}).One(&oldUserInfo)
	if err != nil {
		fmt.Println("查询用户信息失败")
	}
	// 修改一些基本的信息，并不是全部，参看修改页面的内容
	tempInfo := oldUserInfo
	tempInfo.AvatarUrl = newUserInfo.AvatarUrl
	tempInfo.NickName = newUserInfo.NickName
	tempInfo.PenName = newUserInfo.PenName
	tempInfo.Gender = newUserInfo.Gender
	tempInfo.Email = newUserInfo.Email
	tempInfo.Birth = newUserInfo.Birth
	tempInfo.FavoriteAuthor = newUserInfo.FavoriteAuthor
	tempInfo.FavoriteBook = newUserInfo.FavoriteBook
	tempInfo.Intrest = newUserInfo.Intrest
	tempInfo.Introduction = newUserInfo.Introduction

	err = uc.Update(oldUserInfo, tempInfo)
	if err != nil {
		fmt.Println("修改用户信息失败")
	}

	return err
}

func (manager *DbManager) UpdateUserInfoByNickName(nickName string, newUserInfo User) (err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	var oldUserInfo User
	err = uc.Find(bson.M{"nickname": nickName}).One(&oldUserInfo)
	if err != nil {
		fmt.Println("查询用户信息失败")
	}
	// 修改一些基本的信息，并不是全部，参看修改页面的内容
	tempInfo := oldUserInfo
	tempInfo.AvatarUrl = newUserInfo.AvatarUrl
	tempInfo.NickName = newUserInfo.NickName
	tempInfo.PenName = newUserInfo.PenName
	tempInfo.Gender = newUserInfo.Gender
	tempInfo.Email = newUserInfo.Email
	tempInfo.Birth = newUserInfo.Birth
	tempInfo.FavoriteAuthor = newUserInfo.FavoriteAuthor
	tempInfo.FavoriteBook = newUserInfo.FavoriteBook
	tempInfo.Intrest = newUserInfo.Intrest
	tempInfo.Introduction = newUserInfo.Introduction

	err = uc.Update(oldUserInfo, tempInfo)
	if err != nil {
		fmt.Println("修改用户信息失败")
	}

	return err
}

func (manager *DbManager) UpdateMessage(userid int, message Comment) (err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	// 根据用户ID, 查找用户信息
	var oldUserInfo, tmpUser User
	err = uc.Find(bson.M{"id": userid}).One(&oldUserInfo)
	err = uc.Find(bson.M{"id": userid}).One(&tmpUser)

	// 给新留言创建ID，增加日期并格式化
	message.Id = bson.NewObjectId().Hex()
	message.Time = time.Now().Format("2006-01-02 15:04:05")

	// 追加留言到已有的集合
	ms := tmpUser.Message
	ms = append(ms, message)
	tmpUser.Message = ms

	// 更新整个用户信息，包括新加的文章
	err = uc.Update(oldUserInfo, tmpUser)
	if err != nil {
		fmt.Println("添加留言失败")
	}

	return err
}

func (manager *DbManager) AddWatch(watcherNickName string, watchedUserNickName string) (err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	// 根据用户ID, 查找用户信息
	var oldUserInfo, tmpUser User
	err = uc.Find(bson.M{"nickname": watcherNickName}).One(&oldUserInfo)
	err = uc.Find(bson.M{"nickname": watcherNickName}).One(&tmpUser)

	// 查找被关注的用户信息，获取avatarurl
	var watchedUser User
	err = uc.Find(bson.M{"nickname": watchedUserNickName}).One(&watchedUser)

	// 追加关注到已有的集合
	var watch Watch
	watch.NickName = watchedUserNickName
	watch.AvatarUrl = watchedUser.AvatarUrl

	wts := tmpUser.Watch
	for _, v := range wts {
		// 已经关注过此作者 ，返回
		if v.NickName == watchedUserNickName {
			fmt.Println("已经关注此作者")
			return nil
		}
	}
	// 没有关注过此作者，则添加
	wts = append(wts, watch)
	tmpUser.Watch = wts

	// 更新整个用户信息
	err = uc.Update(oldUserInfo, tmpUser)
	if err != nil {
		fmt.Println("添加关注失败")
	}

	return err
}

func (manager *DbManager) AddFans(fansNickName string, ownerNickName string) (err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	// 根据用户ID, 查找用户信息
	var oldUserInfo, tmpUser User
	err = uc.Find(bson.M{"nickname": ownerNickName}).One(&oldUserInfo)
	err = uc.Find(bson.M{"nickname": ownerNickName}).One(&tmpUser)

	// 查找粉丝信息，获取avatarurl
	var fansUser User
	err = uc.Find(bson.M{"nickname": fansNickName}).One(&fansUser)

	// 追加粉丝到已有的集合
	var fans Fans
	fans.NickName = fansNickName
	fans.AvatarUrl = fansUser.AvatarUrl

	fs := tmpUser.Fans
	for _, v := range fs {
		// 已经包含该粉丝 ，返回
		if v.NickName == fansNickName {
			fmt.Println("已经包含该粉丝")
			return nil
		}
	}
	// 没有关包含该粉丝 ，则添加
	fs = append(fs, fans)
	tmpUser.Fans = fs

	// 更新整个用户信息
	err = uc.Update(oldUserInfo, tmpUser)
	if err != nil {
		fmt.Println("添加粉丝失败")
	}

	return err
}

func (manager *DbManager) DeleteWatch(watcherNickName string, watchedUserNickName string) (err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	// 根据用户ID, 查找用户信息
	var oldUserInfo, tmpUser User
	err = uc.Find(bson.M{"nickname": watcherNickName}).One(&oldUserInfo)
	err = uc.Find(bson.M{"nickname": watcherNickName}).One(&tmpUser)

	// 查找，删除关注
	wts := tmpUser.Watch
	count := len(wts)
	var pos int
	for p, v := range wts {
		if v.NickName == watchedUserNickName {
			// 找到该关注，删除数据
			fmt.Println("找到该关注")
			wts = append(wts[:p], wts[p+1:]...)
			break
		}
		pos = p
	}
	// 没有关注则返回，不做任何更改
	if pos == count {
		return nil
	}

	tmpUser.Watch = wts

	// 更新整个用户信息
	err = uc.Update(oldUserInfo, tmpUser)
	if err != nil {
		fmt.Println("删除关注失败")
	}

	return err
}

func (manager *DbManager) DeleteFans(fansNickName string, ownerNickName string) (err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	// 根据用户ID, 查找用户信息
	var oldUserInfo, tmpUser User
	err = uc.Find(bson.M{"nickname": ownerNickName}).One(&oldUserInfo)
	err = uc.Find(bson.M{"nickname": ownerNickName}).One(&tmpUser)

	// 查找，删除粉丝
	var fans Fans
	fans.NickName = fansNickName

	fs := tmpUser.Fans
	count := len(fs)
	var pos int
	for p, v := range fs {
		if v.NickName == fansNickName {
			// 找到该粉丝，删除数据
			fmt.Println("找到该粉丝")
			fs = append(fs[:p], fs[p+1:]...)
			break
		}
		pos = p
	}
	// 没有该粉丝则返回，不做任何更改
	if pos == count {
		return nil
	}

	tmpUser.Fans = fs

	// 更新整个用户信息
	err = uc.Update(oldUserInfo, tmpUser)
	if err != nil {
		fmt.Println("删除粉丝失败")
	}

	return err
}

func (manager *DbManager) UpdateMessageByNickName(nickName string, message Comment) (err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	// 根据用户ID, 查找用户信息
	var oldUserInfo, tmpUser User
	err = uc.Find(bson.M{"nickname": nickName}).One(&oldUserInfo)
	err = uc.Find(bson.M{"nickname": nickName}).One(&tmpUser)

	// 给新留言创建ID，增加日期并格式化
	message.Id = bson.NewObjectId().Hex()
	message.Time = time.Now().Format("2006-01-02 15:04:05")

	// 追加留言到已有的集合
	ms := tmpUser.Message
	ms = append(ms, message)
	tmpUser.Message = ms

	// 更新整个用户信息，包括新加的文章
	err = uc.Update(oldUserInfo, tmpUser)
	if err != nil {
		fmt.Println("添加留言失败")
	}

	return err
}

func (manager *DbManager) AddUserArticle(article *UserArticle) error {
	uc := manager.session.DB(DbName).C(UserCollection)

	// 根据文章作者的ID, 查找作者的信息
	var oldUserInfo User
	err := uc.Find(bson.M{"id": article.AuthorId}).One(&oldUserInfo)
	tmpUser := oldUserInfo

	// 给新文章创建ID和创作日期并格式化
	article.Id = bson.NewObjectId().Hex()
	article.CreateTime = time.Now().Format("2006-01-02 15:04:05")

	// 追加文章到已有的集合
	as := oldUserInfo.Articles
	as = append(as, *article)
	tmpUser.Articles = as

	// 更新整个用户信息，包括新加的文章
	err = uc.Update(oldUserInfo, tmpUser)
	if err != nil {
		fmt.Println("新增文章失败")
	}

	return err
}

func (manager *DbManager) UpdateUserArticle(authorNickName string, newAritlce UserArticle) error {
	uc := manager.session.DB(DbName).C(UserCollection)

	// 根据文章作者的ID, 查找作者的信息
	var oldUserInfo, tmpUser User
	err := uc.Find(bson.M{"nickname": authorNickName}).One(&oldUserInfo)
	err = uc.Find(bson.M{"nickname": authorNickName}).One(&tmpUser)

	flag, index := false, 0
	for _, art := range oldUserInfo.Articles {
		if art.Id == newAritlce.Id {
			flag = true
			fmt.Println("找到指定的文章")
			break
		}
		index += 1
	}

	if flag {
		// 更新指定的文章的类型，标题和内容
		as := tmpUser.Articles
		as[index].Tag = newAritlce.Tag
		as[index].Title = newAritlce.Title
		as[index].Content = newAritlce.Content

		// 更新整个用户信息
		err = uc.Update(oldUserInfo, tmpUser)
		if err != nil {
			fmt.Println("更新失败")
		}
	}

	return err
}

func (manager *DbManager) AddArticleComment(article UserArticle, comment *Comment) error {
	uc := manager.session.DB(DbName).C(UserCollection)

	// 根据文章作者的ID, 查找作者的信息
	var oldUserInfo, tmpUser User
	err := uc.Find(bson.M{"id": article.AuthorId}).One(&oldUserInfo)
	err = uc.Find(bson.M{"id": article.AuthorId}).One(&tmpUser)

	// 给新评论创建ID，增加日期并格式化
	comment.Id = bson.NewObjectId().Hex()
	comment.Time = time.Now().Format("2006-01-02 15:04:05")

	// 追加评论到已有的集合
	cs := article.Comments
	cs = append(cs, *comment)
	article.Comments = cs

	flag, index := false, 0
	for _, art := range oldUserInfo.Articles {
		if art.Id == article.Id {
			flag = true
			fmt.Println("找到指定的文章")
			break
		}
		index += 1
	}

	if flag {
		// 更新该篇文章
		tmpUser.Articles[index] = article
	} else {
		fmt.Println("找不到指定的文章， 添加评论失败")
	}

	// 更新整个用户信息，包括新加的文章
	err = uc.Update(oldUserInfo, tmpUser)
	if err != nil {
		fmt.Println("新增文章失败")
	}

	return err
}

func (manager *DbManager) GetAllArticlesByUserId(userid int) (articles []UserArticle, err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	var userInfo User
	err = uc.Find(bson.M{"id": userid}).One(&userInfo)
	if err != nil {
		fmt.Println("查询用户信息失败")
	}
	articles = userInfo.Articles

	return articles, err
}

func (manager *DbManager) GetAllArticlesByUserNickName(nickName string) (articles []UserArticle, err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	var userInfo User
	err = uc.Find(bson.M{"nickname": nickName}).One(&userInfo)
	if err != nil {
		fmt.Println("查询用户信息失败")
	}
	articles = userInfo.Articles

	return articles, err
}

func (manager *DbManager) GetAllMessageByUserId(userid int) (message []Comment, err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	var userInfo User
	err = uc.Find(bson.M{"id": userid}).One(&userInfo)
	if err != nil {
		fmt.Println("查询用户信息失败")
	}
	message = userInfo.Message

	return message, err
}

func (manager *DbManager) GetArticleCollectionByNickName(nickName string) (articles []ArticleInCollection, err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	var userInfo User
	err = uc.Find(bson.M{"nickname": nickName}).One(&userInfo)
	if err != nil {
		fmt.Println("查询用户信息失败")
	}
	articles = userInfo.ArticleCollection

	return articles, err
}

func (manager *DbManager) AddAticleToArticleCollection(userNickName string, article ArticleInCollection) (err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	// 根据用户ID, 查找用户信息
	var oldUserInfo, tmpUser User
	err = uc.Find(bson.M{"nickname": userNickName}).One(&oldUserInfo)
	err = uc.Find(bson.M{"nickname": userNickName}).One(&tmpUser)

	ac := tmpUser.ArticleCollection
	for _, v := range ac {
		// 已经收藏 ，返回
		if v.Id == article.Id {
			fmt.Println("已经收藏此文章")
			return nil
		}
	}
	// 没有关收藏过，则添加
	ac = append(ac, article)
	tmpUser.ArticleCollection = ac

	// 更新整个用户信息
	err = uc.Update(oldUserInfo, tmpUser)
	if err != nil {
		fmt.Println("添加留言失败")
	}

	return err
}

func (manager *DbManager) DeleteAticleFromArticleCollection(userNickName string, article ArticleInCollection) (err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	// 根据用户ID, 查找用户信息
	var oldUserInfo, tmpUser User
	err = uc.Find(bson.M{"nickname": userNickName}).One(&oldUserInfo)
	err = uc.Find(bson.M{"nickname": userNickName}).One(&tmpUser)

	ac := tmpUser.ArticleCollection
	count := len(ac)
	var pos int
	for p, v := range ac {
		// 找到该文章 ，删除
		if v.Id == article.Id {
			ac = append(ac[:p], ac[p+1:]...)
			break
		}
		pos = p
	}
	// 没有收藏过该文章则返回，不做任何更改
	if pos == count {
		return nil
	}

	tmpUser.ArticleCollection = ac

	// 更新整个用户信息
	err = uc.Update(oldUserInfo, tmpUser)
	if err != nil {
		fmt.Println("添加留言失败")
	}

	return err
}

func (manager *DbManager) GetAllFansByUserNickName(nickName string) (fans []Fans, err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	var userInfo User
	err = uc.Find(bson.M{"nickname": nickName}).One(&userInfo)
	if err != nil {
		fmt.Println("查询用户信息失败")
	}
	fans = userInfo.Fans

	return fans, err
}

func (manager *DbManager) GetAllWatchByUserNickName(nickName string) (watch []Watch, err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	var userInfo User
	err = uc.Find(bson.M{"nickname": nickName}).One(&userInfo)
	if err != nil {
		fmt.Println("查询用户信息失败")
	}
	watch = userInfo.Watch

	return watch, err
}

func (manager *DbManager) GetAllMessageByUserNickName(nickName string) (message []Comment, err error) {
	uc := manager.session.DB(DbName).C(UserCollection)

	var userInfo User
	err = uc.Find(bson.M{"nickname": nickName}).One(&userInfo)
	if err != nil {
		fmt.Println("查询用户信息失败")
	}
	message = userInfo.Message

	return message, err
}

func (manager *DbManager) GetArticleByUserIdAndArticleId(userid int, articleid string) (article UserArticle, err error) {
	articles, _ := manager.GetAllArticlesByUserId(userid)

	flag := false
	for _, art := range articles {
		if art.Id == articleid {
			article = art
			flag = true
			fmt.Println("找到指定的文章")
			return article, err
		}
	}

	if !flag {
		fmt.Println("找到指定的文章")
	}

	return article, err
}

func (manager *DbManager) GetArticleByUserNickNameAndArticleId(nickName string, articleid string) (article UserArticle, err error) {
	articles, _ := manager.GetAllArticlesByUserNickName(nickName)

	flag := false
	for _, art := range articles {
		if art.Id == articleid {
			article = art
			flag = true
			fmt.Println("找到指定的文章")
			return article, err
		}
	}

	if !flag {
		fmt.Println("找到指定的文章")
	}

	return article, err
}

func (manager *DbManager) DeleteUserArticle(nickName string, articleid string) error {
	uc := manager.session.DB(DbName).C(UserCollection)

	// 根据文章作者的ID, 查找作者的信息
	var oldUserInfo User
	err := uc.Find(bson.M{"nickname": nickName}).One(&oldUserInfo)
	tmpUser := oldUserInfo

	flag, index := false, 0
	for _, art := range oldUserInfo.Articles {
		if art.Id == articleid {
			flag = true
			fmt.Println("找到指定的文章")
			break
		}
		index += 1
	}

	if flag {
		// 删除指定的文章
		as := oldUserInfo.Articles
		as = append(as[:index], as[index+1:]...)

		// 更新整个用户信息，包括新加的文章
		tmpUser.Articles = as
		err = uc.Update(oldUserInfo, tmpUser)
		if err != nil {
			fmt.Println("更新失败")
		}
	}

	return err
}
