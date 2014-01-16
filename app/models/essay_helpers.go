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

	err := uc.Insert(e)

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
