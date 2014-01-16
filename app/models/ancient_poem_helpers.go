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
		return errors.New("此篇古体诗（词）已经存在")
	}

	err := uc.Insert(pm)

	return err
}

func (manager *DbManager) GetAllAncientPoem() ([]AncientPoem, error) {
	uc := manager.session.DB(DbName).C(AncientPoemCollection)

	count, err := uc.Count()
	fmt.Println("共有古诗词： ", count, "篇")

	poems := []AncientPoem{}
	err = uc.Find(nil).All(&poems)

	return poems, err
}
