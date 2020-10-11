package model

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//User 表基本信息
type User struct {
	UID        int    `gorm:"primary_key"`
	Email      string `gorm:"unique"`
	Password   string
	Status     int
	Createtime string
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
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	u1 := User{1, "admin@1.com", "123456", 1, currentTime}
	db.Create(&u1)
}

//DbSel 数据查询
func DbSel(u *User, email, pass string) int {
	//连接数据库
	db, err := gorm.Open(sqlite.Open("./foo.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Select("uid").Where("email=? AND password=?", email, pass).Find(u)
	return u.UID
}

//DbInsert 注册插入数据
func DbInsert(email string, pass string) error {
	db, err := gorm.Open(sqlite.Open("./foo.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	var user User
	db.Last(&user)
	lastid := user.UID
	newid := lastid + 1
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	u := User{newid, email, pass, 1, currentTime}
	res := db.Create(&u)
	return res.Error
}
