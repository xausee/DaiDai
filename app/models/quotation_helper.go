package models

import (
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
)

func (manager *DbManager) AddQuotation(q *Quotation) error {
	uc := manager.session.DB(DbName).C(QuotationCollection)

	i, _ := uc.Find(bson.M{"Content": q.Content}).Count()
	if i != 0 {
		return errors.New("此条摘录已经存在")
	}

	q.Id = bson.NewObjectId().Hex()
	err := uc.Insert(q)

	return err
}

func (manager *DbManager) GetAllQuotation() ([]Quotation, error) {
	uc := manager.session.DB(DbName).C(QuotationCollection)

	count, err := uc.Count()
	fmt.Println("共有摘录： ", count, "条")
	allquotation := []Quotation{}
	err = uc.Find(nil).All(&allquotation)

	return allquotation, err
}

func (manager *DbManager) GetQuotationByTag(tag string) ([]Quotation, error) {
	uc := manager.session.DB(DbName).C(QuotationCollection)

	allquotation := []Quotation{}
	err := uc.Find(bson.M{"tag": tag}).All(&allquotation)

	return allquotation, err
}

func (manager *DbManager) GetQuotationById(id string) (q *Quotation, err error) {
	uc := manager.session.DB(DbName).C(QuotationCollection)
	err = uc.Find(bson.M{"id": id}).One(&q)
	return
}

func (manager *DbManager) DeleteQuotationById(id string) (err error) {
	uc := manager.session.DB(DbName).C(QuotationCollection)
	err = uc.Remove(bson.M{"id": id})
	return err
}
