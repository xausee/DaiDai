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
		return errors.New("此条摘录已经存在")
	}

	err := uc.Insert(e)

	return err
}

func (manager *DbManager) GetAllEssay() ([]Essay, error) {
	uc := manager.session.DB(DbName).C(EssayCollection)

	allEssay := []Essay{}
	err := uc.Find(bson.M{"tag": "散文"}).All(&allEssay)

	count, err := uc.Count()
	fmt.Println("所有的条目数： ", count)

	return allEssay, err
}
