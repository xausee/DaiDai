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
		return errors.New("此条摘录已经存在")
	}

	err := uc.Insert(pm)

	return err
}

func (manager *DbManager) GetAllModernPoem() ([]ModernPoem, error) {
	uc := manager.session.DB(DbName).C(ModernPoemCollection)

	poems := []ModernPoem{}
	err := uc.Find(bson.M{"tag": "现代诗"}).All(&poems)

	count, err := uc.Count()
	fmt.Println("所有的条目数： ", count)

	return poems, err
}
