package models

import (
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
)

func (manager *DbManager) AddHintFiction(hf *HintFiction) error {
	uc := manager.session.DB(DbName).C(HintFictionCollection)

	i, _ := uc.Find(bson.M{"Content": hf.Content}).Count()
	if i != 0 {
		return errors.New("此篇微小说已经存在")
	}

	hf.Id = bson.NewObjectId().Hex()
	err := uc.Insert(hf)

	return err
}

func (manager *DbManager) GetAllHintFiction() ([]HintFiction, error) {
	uc := manager.session.DB(DbName).C(HintFictionCollection)

	count, err := uc.Count()
	fmt.Println("共有微小说： ", count, "篇")

	hfs := []HintFiction{}
	err = uc.Find(nil).All(&hfs)

	return hfs, err
}

func (manager *DbManager) GetHintFictionByTag(tag string) ([]HintFiction, error) {
	uc := manager.session.DB(DbName).C(HintFictionCollection)

	hfs := []HintFiction{}
	err := uc.Find(bson.M{"tag": tag}).All(&hfs)

	return hfs, err
}

func (manager *DbManager) GetHintFictionById(id string) (hf *HintFiction, err error) {
	uc := manager.session.DB(DbName).C(HintFictionCollection)
	err = uc.Find(bson.M{"id": id}).One(&hf)
	return
}

func (manager *DbManager) DeleteHintFictionById(id string) (err error) {
	uc := manager.session.DB(DbName).C(HintFictionCollection)
	err = uc.Remove(bson.M{"id": id})
	return err
}
