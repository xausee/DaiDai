package models

import (
	//"fmt"
	"labix.org/v2/mgo/bson"
)

func (manager *DbManager) Search(keywords string) (results []ArticleBasicInfo, err error) {
	type Items map[string]string

	re0 := []Essay{}
	re1 := []Essay{}
	uc := manager.session.DB(DbName).C(EssayCollection)
	err = uc.Find(bson.M{"title": Items{"$regex": keywords}}).All(&re0)
	err = uc.Find(bson.M{"author": Items{"$regex": keywords}}).All(&re1)
	//err = uc.Find(bson.M{"title": Items{"$regex": bson.RegEx{1}}}).All(&re0)

	rm0 := []ModernPoem{}
	rm1 := []ModernPoem{}
	uc = manager.session.DB(DbName).C(ModernPoemCollection)
	err = uc.Find(bson.M{"title": Items{"$regex": keywords}}).All(&rm0)
	err = uc.Find(bson.M{"author": Items{"$regex": keywords}}).All(&rm1)

	rq0 := []Quotation{}
	rq1 := []Quotation{}
	uc = manager.session.DB(DbName).C(QuotationCollection)
	err = uc.Find(bson.M{"title": Items{"$regex": keywords}}).All(&rq0)
	err = uc.Find(bson.M{"author": Items{"$regex": keywords}}).All(&rq1)

	rh0 := []HintFiction{}
	rh1 := []HintFiction{}
	uc = manager.session.DB(DbName).C(HintFictionCollection)
	err = uc.Find(bson.M{"title": Items{"$regex": keywords}}).All(&rh0)
	err = uc.Find(bson.M{"author": Items{"$regex": keywords}}).All(&rh1)

	for _, f := range re0 {
		var info ArticleBasicInfo
		info.Id = f.Id
		info.Tag = "Essay"
		info.Title = f.Title
		info.Author = f.Author
		results = append(results, info)
	}

	for _, f := range re1 {
		var info ArticleBasicInfo
		info.Id = f.Id
		info.Tag = "Essay"
		info.Title = f.Title
		info.Author = f.Author
		results = append(results, info)
	}

	for _, f := range rm0 {
		var info ArticleBasicInfo
		info.Id = f.Id
		info.Tag = "ModernPoem"
		info.Title = f.Title
		info.Author = f.Author
		results = append(results, info)
	}

	for _, f := range rm1 {
		var info ArticleBasicInfo
		info.Id = f.Id
		info.Tag = "ModernPoem"
		info.Title = f.Title
		info.Author = f.Author
		results = append(results, info)
	}

	for _, f := range rq0 {
		var info ArticleBasicInfo
		info.Id = f.Id
		info.Tag = "Quotation"
		info.Title = f.OriginalTitle
		info.Author = f.Author
		results = append(results, info)
	}

	for _, f := range rq1 {
		var info ArticleBasicInfo
		info.Id = f.Id
		info.Tag = "Quotation"
		info.Title = f.OriginalTitle
		info.Author = f.Author
		results = append(results, info)
	}

	for _, f := range rh0 {
		var info ArticleBasicInfo
		info.Id = f.Id
		info.Tag = "HintFiction"
		info.Title = f.Title
		info.Author = f.Author
		results = append(results, info)
	}

	for _, f := range rh1 {
		var info ArticleBasicInfo
		info.Id = f.Id
		info.Tag = "HintFiction"
		info.Title = f.Title
		info.Author = f.Author
		results = append(results, info)
	}

	resultss := RemoveDuplicates(results)

	return resultss, err
}
