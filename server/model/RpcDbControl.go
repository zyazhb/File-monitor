package model

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//RPCDb 表基本信息
type RPCDb struct {
	RID        int `gorm:"primary_key"`
	FileName   string
	Operation  string
	Createtime string
	Hash       int
}

//RPCDbInit 连接数据库,表迁移
func RPCDbInit() {
	//连接数据库
	db, err := gorm.Open(sqlite.Open("./report.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//迁移 schema
	db.AutoMigrate(&RPCDb{})
	newdb := RPCDb{0, "nil", "nil", "nil", 0}
	db.Create(&newdb)
}

//RPCDbSel 数据查询
func RPCDbSel() *gorm.DB {
	//连接数据库
	db, err := gorm.Open(sqlite.Open("./report.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db.Select("*")
}

//RPCDbInsert 注册插入数据
func RPCDbInsert(filename string, operation string, hash []byte) error {
	db, err := gorm.Open(sqlite.Open("./report.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	var rpcdb RPCDb
	db.Last(&rpcdb)
	lastid := rpcdb.RID
	newid := lastid + 1
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	hashf, err := fmt.Printf("%x", hash)
	newreport := RPCDb{newid, filename, operation, currentTime, hashf}
	res := db.Create(&newreport)
	return res.Error
}
