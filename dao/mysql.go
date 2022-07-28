package dao

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

//type MsDB struct {
//	*sql.DB
//}
type MsDB struct {
	*sqlx.DB
}

var DB MsDB

func init() {
	//执行main之前 先执行init方法
	dataSourceName := fmt.Sprintf("root:12345678@tcp(localhost:3306)/go_db?charset=utf8")
	db, err := sqlx.Open("mysql", dataSourceName)
	if err != nil {
		log.Println("连接数据库异常")
		panic(err)
	}
	//最大空闲连接数，默认不配置，是2个最大空闲连接
	db.SetMaxIdleConns(5)
	//最大连接数，默认不配置，是不限制最大连接数
	db.SetMaxOpenConns(100)
	// 连接最大存活时间
	db.SetConnMaxLifetime(time.Minute * 3)
	//空闲连接最大存活时间
	db.SetConnMaxIdleTime(time.Minute * 1)
	err = db.Ping()
	if err != nil {
		log.Println("数据库无法连接")
		_ = db.Close()
		panic(err)
	}
	DB = MsDB{db}
}
