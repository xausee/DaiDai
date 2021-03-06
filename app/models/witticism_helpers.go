package models

import (
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
)

func (manager *DbManager) AddWitticism(w *Witticism) error {
	uc := manager.session.DB(DbName).C(WitticismCollection)

	i, _ := uc.Find(bson.M{"Content": w.Content}).Count()
	if i != 0 {
		return errors.New("此条慧语已经存在")
	}

	w.Id = bson.NewObjectId().Hex()
	err := uc.Insert(w)

	return err
}

func (manager *DbManager) UpdateWitticism(originalWitticismID string, newWitticism *Witticism) error {
	uc := manager.session.DB(DbName).C(WitticismCollection)

	var originalWitticism *Witticism
	newWitticism.Id = originalWitticismID

	err := uc.Find(bson.M{"id": originalWitticismID}).One(&originalWitticism)	
	err = uc.Update(originalWitticism, newWitticism)

	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (manager *DbManager) GetAllWitticism() ([]Witticism, error) {
	uc := manager.session.DB(DbName).C(WitticismCollection)

	count, err := uc.Count()
	fmt.Println("共有慧语： ", count, "条")
	allWitticism := []Witticism{}
	err = uc.Find(nil).All(&allWitticism)

	return allWitticism, err
}

func (manager *DbManager) GetWitticismById(id string) (w *Witticism, err error) {
	uc := manager.session.DB(DbName).C(WitticismCollection)
	err = uc.Find(bson.M{"id": id}).One(&w)
	return
}

func (manager *DbManager) DeleteWitticismById(id string) (err error) {
	uc := manager.session.DB(DbName).C(WitticismCollection)
	err = uc.Remove(bson.M{"id": id})
	return err
}