package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type MysqlDB struct {
	*gorm.DB
	ServiceAddr string
}

type MysqlFactory struct {
}

type MysqlDBer interface {
	NewDBClient() *MysqlDB
}

func (MysqlFactory) Create(sAddr string) MysqlDBer {
	return &MysqlDB{
		ServiceAddr: sAddr,
	}
}

func (Mdb *MysqlDB) NewDBClient() *MysqlDB {
	url := fmt.Sprintf("root:example@tcp(%s)/dbtest?charset=utf8&parseTime=True&loc=Local", Mdb.ServiceAddr)
	db, err := gorm.Open("mysql", url)
	if err != nil {
		panic(err)
	}
	Mdb.DB = db
	return Mdb
}
