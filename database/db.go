package database

import "github.com/jinzhu/gorm"

// あとで直す
const (
	dbms     = "mysql"
	user     = "jb5"
	pass     = "h19life"
	protocol = "tcp(localhost:3306)"
	dbname   = "et"
)

// SQLConnect establishes a mysql connection.
func SQLConnect() (database *gorm.DB, err error) {
	DBMS := dbms
	USER := user
	PASS := pass
	PROTOCOL := protocol
	DBNAME := dbname

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"
	return gorm.Open(DBMS, CONNECT)
}
