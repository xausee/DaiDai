package models

import (
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
)

const (
	DbName                = "SanWenJia"
	UserCollection        = "user"
	QuotationCollection   = "quotation"
	WitticismCollection   = "witticism"
	AncientPoemCollection = "ancientPoem"
	ModernPoemCollection  = "modernPoem"
	EssayCollection       = "essay"
	HintFictionCollection = "hintFiction"
)

type DbManager struct {
	session *mgo.Session
}

func NewDbManager() (*DbManager, error) {
	revel.Config.SetSection("db")
	ip, found := revel.Config.String("ip")
	if !found {
		revel.ERROR.Fatal("Cannot load database ip from app.conf")
	}

	session, err := mgo.Dial(ip)
	if err != nil {
		return nil, err
	}

	return &DbManager{session}, nil
}

func (manager *DbManager) Close() {
	manager.session.Close()
}
