package model

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

//DbInit 初始化Db
func DbInit() {
	log.Println("打开数据")
	db, err := sql.Open("sqlite3", "./foo.db")
	checkErr(err)

	log.Println("生成数据表")
	SQLTable := `
CREATE TABLE IF NOT EXISTS "users" (
   "uid" INTEGER PRIMARY KEY AUTOINCREMENT,
   "username" VARCHAR(64) NULL,
   "password" VARCHAR(64) NULL,
   "created" TIMESTAMP default (datetime('now', 'localtime'))  
);

CREATE TABLE IF NOT EXISTS "userinfo" (
   "uid" INTEGER PRIMARY KEY AUTOINCREMENT,
   "username" VARCHAR(64) NULL,
   "departname" VARCHAR(64) NULL,
   "created" TIMESTAMP default (datetime('now', 'localtime'))  
);

CREATE TABLE IF NOT EXISTS "userdeatail" (
   "uid" INT(10) NULL,
   "intro" TEXT NULL,
   "profile" TEXT NULL,
   PRIMARY KEY (uid)
);
   `
	db.Exec(SQLTable)
}

//DbInsert 增加数据
func DbInsert() {
	db, err := sql.Open("sqlite3", "./foo.db")
	checkErr(err)

	//插入数据
	log.Print("插入数据, ID=")
	stmt, err := db.Prepare("INSERT INTO userinfo(username, departname)  values(?, ?)")
	checkErr(err)
	res, err := stmt.Exec("astaxie", "研发部门")
	checkErr(err)
	id, err := res.LastInsertId()
	checkErr(err)
	log.Println(id)

	//更新数据
	log.Print("更新数据 ")
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)
	res, err = stmt.Exec("astaxieupdate", id)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	log.Println(affect)

	//查询数据
	log.Println("查询数据")
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)
	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		log.Println(uid, username, department, created)
	}

	/*
	   //删除数据
	   log.Println("删除数据")
	   stmt, err = db.Prepare("delete from userinfo where uid=?")
	   checkErr(err)
	   res, err = stmt.Exec(id)
	   checkErr(err)
	   affect, err = res.RowsAffected()
	   checkErr(err)
	   log.Println(affect)
	*/
	db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
