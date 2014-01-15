package models

import (
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
)

func (manager *DbManager) AddAncientPeom(pm *AncientPoem) error {
	uc := manager.session.DB(DbName).C(AncientPoemCollection)

	i, _ := uc.Find(bson.M{"Content": pm.Content}).Count()
	if i != 0 {
		return errors.New("此条摘录已经存在")
	}

	err := uc.Insert(pm)

	return err
}

func (manager *DbManager) GetAllAncientPoem() ([]AncientPoem, error) {
	uc := manager.session.DB(DbName).C(AncientPoemCollection)

	count, err := uc.Count()
	fmt.Println("所有的条目数： ", count)

	poems := []AncientPoem{}
	err = uc.Find(bson.M{"tag": "古体诗歌"}).All(&poems)

	return poems, err
}
