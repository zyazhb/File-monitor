package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//User 表基本信息
type User struct {
	UID        int    `gorm:"primary_key"`
	Email      string `gorm:"unique"`
	Password   string
	Status     int
	Createtime int64
}

//DbInit 连接数据库,表迁移
func DbInit() {
	//连接数据库
	db, err := gorm.Open(sqlite.Open("./foo.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//迁移 schema
	db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{})
	u1 := User{0, "admin@1.com", "123456", 1, 2020}
	db.Create(&u1)
}

//DbSel 数据查询
func DbSel(u *User) {
	//连接数据库
	db, err := gorm.Open(sqlite.Open("./foo.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Select("email", "password").Find(u)
}
