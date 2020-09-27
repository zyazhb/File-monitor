package model

//User 表基本信息
type User struct {
	UID      uint   `gorm:"primary_key"`
	Email    string `gorm:"unique"`
	Password string
}

//DbInit 连接数据库,表迁移
func DbInit() {
	// db, err := gorm.Open(sqlite.Open("./foo.db"), &gorm.Config{})
	// if err != nil {
	// 	panic(err)
	// }
	//迁移 schema
	// db.AutoMigrate(&User{})
	// //创建测试数据
	// u1 := User{1, "123@123.com", "123456"}
	// db.Create(&u1)
}

//DbSel 数据查询
// func DbSel() {
// 	db, err := gorm.Open(sqlite.Open("./foo.db"), &gorm.Config{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	var user User
// 	db.Select("email").Find(&user)
// 	fmt.Println("", user)
// }

//DbInit 初始化Db
// func DbInit() {

// 	db, err := gorm.Open("sqlite3", "./foo.db")
// 	if err != nil {
// 		panic("failed to connect database")
// 	}
// 	// 迁移 schema
// 	db.AutoMigrate(&User{})
// }

// //DbInsert 增加数据
// func DbInsert() {
// 	db, err := sql.Open("sqlite3", "./foo.db")
// 	checkErr(err)

// 	//插入数据
// 	log.Print("插入数据, ID=")
// 	stmt, err := db.Prepare("INSERT INTO userinfo(username, departname)  values(?, ?)")
// 	checkErr(err)
// 	res, err := stmt.Exec("astaxie", "研发部门")
// 	checkErr(err)
// 	id, err := res.LastInsertId()
// 	checkErr(err)
// 	log.Println(id)

// 	//更新数据
// 	log.Print("更新数据 ")
// 	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
// 	checkErr(err)
// 	res, err = stmt.Exec("astaxieupdate", id)
// 	checkErr(err)
// 	affect, err := res.RowsAffected()
// 	checkErr(err)
// 	log.Println(affect)

// 	//查询数据
// 	log.Println("查询数据")
// 	rows, err := db.Query("SELECT * FROM userinfo")
// 	checkErr(err)
// 	for rows.Next() {
// 		var uid int
// 		var username string
// 		var department string
// 		var created string
// 		err = rows.Scan(&uid, &username, &department, &created)
// 		checkErr(err)
// 		log.Println(uid, username, department, created)
// 	}

// 	/*
// 	   //删除数据
// 	   log.Println("删除数据")
// 	   stmt, err = db.Prepare("delete from userinfo where uid=?")
// 	   checkErr(err)
// 	   res, err = stmt.Exec(id)
// 	   checkErr(err)
// 	   affect, err = res.RowsAffected()
// 	   checkErr(err)
// 	   log.Println(affect)
// 	*/
// 	db.Close()
// }

// func checkErr(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }
