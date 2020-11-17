package config

type Database struct {
	SqlPath string
	DbPath  string
}

func (c *Database) initialization() {
	c.SqlPath = "./script/sqlite.sql"
	c.DbPath = "./go-panel.db"

	Conf.Database = c
}
