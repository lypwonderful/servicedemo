package db

import (
	"github.com/lypwonderful/servicedemo/pkgs/core/db/mongodb"
)

type DBClients struct {
	mongodb.MongoDB
}

type DBer interface {
	mongodb.MongoDBer
}

//func t() {
//	dbObj := mongodb.MongoFactory{}
//	dbCli := dbObj.Create("").NewDBClient()
//}
