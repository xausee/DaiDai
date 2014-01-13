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
	fmt.Println("Total quotation count is ", count)
	allquotation := []Quotation{}
	err = uc.Find(nil).All(&allquotation)
	for _, quotation := range allquotation {
		fmt.Println(quotation)
		fmt.Println("==================")
	}

	return allquotation, err
}
