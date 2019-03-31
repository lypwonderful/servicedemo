package db

import (
	"github.com/lypwonderful/servicedemo/pkgs/core/db/mongodb"
	"github.com/lypwonderful/servicedemo/pkgs/core/db/mysql"
)

type DBClients struct {
	mongodb.MongoDB
	mysql.MysqlDB
}

//type DBer interface {
//	mongodb.MongoDBer
//	mysql.MysqlDBer
//}

//func t() {
//	dbObj := mongodb.MongoFactory{}
//	dbCli := dbObj.Create("").NewDBClient()
//}
