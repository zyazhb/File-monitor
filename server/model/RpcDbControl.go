package model

import (
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
	Hash       string
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
	if len(RPCDbSel()) < 1 {
		db.Create(&RPCDb{})
	}
}

//RPCDbSel 数据查询
func RPCDbSel() []RPCDb {
	//连接数据库
	db, err := gorm.Open(sqlite.Open("./report.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	var rpcdb []RPCDb
	db.Find(&rpcdb)
	return rpcdb
}

//RPCDbInsert 注册插入数据
func RPCDbInsert(filename string, operation string, hash string) error {
	db, err := gorm.Open(sqlite.Open("./report.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	var rpcdb RPCDb
	db.Last(&rpcdb)
	lastid := rpcdb.RID
	newid := lastid + 1
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	newreport := RPCDb{newid, filename, operation, currentTime, hash}
	res := db.Create(&newreport)
	return res.Error
}
