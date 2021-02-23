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
	Role       int
	Createtime string
}

//DbInit 连接数据库,表迁移
func DbInit() {
	//连接数据库
	db, err := gorm.Open(sqlite.Open("./user.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//迁移 schema
	db.AutoMigrate(&User{})
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	email := "admin@null.com"
	password := GetRandomString(time.Now().UnixNano())
	PrintLog("Your email is: " + email)
	PrintLog("Your password is: " + password)
	u1 := User{1, email, GenMD5(password), 0, currentTime}
	db.Create(&u1)
}

//DbSel 数据查询
func DbSel(u *User, email, passmd5 string) (int, int) {
	//连接数据库
	db, err := gorm.Open(sqlite.Open("./user.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Where("email=? AND password=?", email, passmd5).Find(u)
	return u.UID, u.Role
}

//AllUserInfo 全部用户信息
func AllUserInfo() []User {
	//连接数据库
	db, err := gorm.Open(sqlite.Open("./user.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	var userdb []User
	db.Order("UID").Find(&userdb)
	return userdb
}

//DbInsert 注册插入数据
func DbInsert(email string, passmd5 string) error {
	db, err := gorm.Open(sqlite.Open("./user.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	var user User
	db.Last(&user)
	lastid := user.UID
	newid := lastid + 1
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	u := User{newid, email, passmd5, 999, currentTime} //默认role999 意味无权限
	res := db.Create(&u)
	return res.Error
}
