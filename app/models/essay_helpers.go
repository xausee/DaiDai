package models

import (
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
)

func (manager *DbManager) AddEssay(e *Essay) error {
	uc := manager.session.DB(DbName).C(EssayCollection)

	i, _ := uc.Find(bson.M{"Content": e.Content}).Count()
	if i != 0 {
		return errors.New("此篇散文已经存在")
	}

	e.Id = bson.NewObjectId().Hex()
	err := uc.Insert(e)

	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (manager *DbManager) GetAllEssay() ([]Essay, error) {
	uc := manager.session.DB(DbName).C(EssayCollection)

	count, err := uc.Count()
	fmt.Println("共有散文 ", count, "篇")

	allEssay := []Essay{}
	err = uc.Find(nil).All(&allEssay)

	return allEssay, err
}

func (manager *DbManager) GetEssayByTitle(title string) ([]Essay, error) {
	uc := manager.session.DB(DbName).C(EssayCollection)

	count, err := uc.Count()
	fmt.Println("共有散文 ", count, "篇")

	allEssay := []Essay{}
	err = uc.Find(bson.M{"title": title}).All(&allEssay)

	return allEssay, err
}

func (manager *DbManager) GetEssayById(id string) (e *Essay, err error) {
	uc := manager.session.DB(DbName).C(EssayCollection)
	err = uc.Find(bson.M{"id": id}).One(&e)
	return
}
