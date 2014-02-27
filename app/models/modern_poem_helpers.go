package models

import (
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
)

func (manager *DbManager) AddModernPeom(pm *ModernPoem) error {
	uc := manager.session.DB(DbName).C(ModernPoemCollection)

	i, _ := uc.Find(bson.M{"Content": pm.Content}).Count()
	if i != 0 {
		return errors.New("此篇现代诗已经存在")
	}

	pm.Id = bson.NewObjectId().Hex()
	err := uc.Insert(pm)

	return err
}

func (manager *DbManager) GetAllModernPoem() ([]ModernPoem, error) {
	uc := manager.session.DB(DbName).C(ModernPoemCollection)

	count, err := uc.Count()
	fmt.Println("共有现代诗： ", count, "篇")

	poems := []ModernPoem{}
	err = uc.Find(nil).All(&poems)

	return poems, err
}

func (manager *DbManager) GetModernPoemByTag(tag string) ([]ModernPoem, error) {
	uc := manager.session.DB(DbName).C(ModernPoemCollection)
	poems := []ModernPoem{}
	err := uc.Find(bson.M{"tag": tag}).All(&poems)
	return poems, err
}

func (manager *DbManager) GetModernPoemById(id string) (mp *ModernPoem, err error) {
	uc := manager.session.DB(DbName).C(ModernPoemCollection)
	err = uc.Find(bson.M{"id": id}).One(&mp)
	return
}
