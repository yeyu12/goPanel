package core

import (
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"goPanel/src/panel/config"
	"xorm.io/core"
)

var Db *xorm.Engine

func init() {
	var err error

	if Db, err = xorm.NewEngine("sqlite3", config.Conf.Database.DbPath); err != nil {
		panic(err)
	}

	if err = Db.Ping(); err != nil {
		panic(err)
	}

	Db.ShowSQL(true)
	Db.Logger().SetLevel(core.LOG_DEBUG)
}

func CreateTables(beans ...interface{}) {
	for _, item := range beans {
		isTable, _ := Db.IsTableExist(item)
		if !isTable {
			_ = Db.CreateTables(item)
		}

		_ = Db.Sync2(item)
	}
}
