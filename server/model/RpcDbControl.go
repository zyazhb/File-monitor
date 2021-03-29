package model

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//RPCDb 表基本信息
type RPCDb struct {
	RID        int `gorm:"primary_key"`
	AgentIP    string
	FileName   string
	Operation  string
	Createtime string
	Hash       string
}

var db *gorm.DB
var err error

//RPCDbInit 连接数据库,表迁移
func RPCDbInit() {
	//连接数据库
	db, err = gorm.Open(sqlite.Open("./report.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//迁移 schema
	db.AutoMigrate(&RPCDb{})
	if len(RPCDbSel(1)) < 1 {
		db.Create(&RPCDb{})
	}
}

//RPCDbSel 数据查询
func RPCDbSel(page int) []RPCDb {
	//连接数据库
	if err != nil {
		panic(err)
	}
	var rpcdb []RPCDb
	//暂时pagesize是硬编码的
	pageSize := 15
	DB := db.Limit(pageSize).Offset((page - 1) * pageSize)
	DB.Order("r_id desc").Find(&rpcdb)
	return rpcdb
}

//RPCDbPageCount 计算总数
func RPCDbPageCount() int64 {
	//连接数据库
	if err != nil {
		panic(err)
	}
	var count int64
	db.Model(&RPCDb{}).Count(&count)
	return count / 15
}

//RPCDbDel 数据删除
func RPCDbDel(rid string) {
	//连接数据库
	if err != nil {
		panic(err)
	}
	db.Where("r_id=?", rid).Delete(RPCDb{})
}

//RPCDbInsert RPC数据库插入数据
func RPCDbInsert(agentip string, filename string, operation string, hash string) error {
	if err != nil {
		panic(err)
	}

	var rpcdb RPCDb
	db.Last(&rpcdb)
	newid := rpcdb.RID + 1
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	newreport := RPCDb{newid, agentip, filename, operation, currentTime, hash}
	res := db.Create(&newreport)
	return res.Error
}
