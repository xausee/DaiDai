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

	err := uc.Insert(q)

	return err
}

func (manager *DbManager) GetAllQuotation() ([]Quotation, error) {
	uc := manager.session.DB(DbName).C(QuotationCollection)

	count, err := uc.Count()
	fmt.Println("所有的条目数： ", count)
	allquotation := []Quotation{}
	err = uc.Find(nil).All(&allquotation)
	for _, quotation := range allquotation {
		fmt.Println(quotation)
		fmt.Println("==================")
	}

	return allquotation, err
}

func (manager *DbManager) GetAllWitticismQuotation() ([]Quotation, error) {
	uc := manager.session.DB(DbName).C(QuotationCollection)

	allquotation := []Quotation{}
	err := uc.Find(bson.M{"tag": "名人语录"}).All(&allquotation)

	return allquotation, err
}
