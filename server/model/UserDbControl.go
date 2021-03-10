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

//DbGetByuid 展示所有信息用于修改
func DbGetByuid(u *User, uid int) (user User) {
	db, err := gorm.Open(sqlite.Open("./user.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Where("uid=? ", uid).Find(u)
	return *u
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

//DbAddUser 添加用户信息
func DbAddUser(email string, passmd5 string, role int) error {
	db, err := gorm.Open(sqlite.Open("./user.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	var user User
	db.Last(&user)
	lastid := user.UID
	newid := lastid + 1
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	u := User{newid, email, passmd5, role, currentTime} //默认role999 意味无权限
	res := db.Create(&u)
	return res.Error
}

//EditUser 修改用户信息
func EditUser(uid, role int, email, pass string) {
	db, err := gorm.Open(sqlite.Open("./user.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	var user User
	db.Model(&user).Where("UID=?", uid).Updates(map[string]interface{}{
		"email":    email,
		"password": pass,
		"role":     role})
}

//DbDelUser 删除用户信息
func DbDelUser(uid string) {
	db, err := gorm.Open(sqlite.Open("./user.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Where("UID=?", uid).Delete(User{})
}
