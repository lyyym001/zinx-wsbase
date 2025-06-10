package znet

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type SqliteHandle struct {
	DB *sql.DB
}

func NewSqliteHandle(dbPath string) *SqliteHandle {

	sqlite := &SqliteHandle{}
	sqlite.Connect(dbPath)
	return sqlite

}

func (sqlite *SqliteHandle) Connect(dbPath string) {

	db, err := sql.Open("sqlite3", dbPath) //还要装gcc 算了
	if err != nil {
		fmt.Println("Sqlite3 Open Err")
		return
	}

	fmt.Println("3.dbOpened")

	sqlite.DB = db
	sqlite.DB.SetMaxOpenConns(2000) //用于设置最大打开的连接数，默认值为0表示不限制。
	sqlite.DB.SetMaxIdleConns(1000) //用于设置闲置的连接数。
	//fmt.Println("Sqlite Ping = ",sqlite.DB.Ping())                //保持链接  貌似没啥用

	//sqlite.Query()

}

func (sqlite *SqliteHandle) GetDB() *sql.DB {

	return sqlite.DB

}

// 查询数据
func (r *SqliteHandle) Query() {
	sql := r.DB
	rows, err := sql.Query("SELECT * FROM tb_user")
	fmt.Println("err - >", err)
	//checkErr(err)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var username string
			var department string
			var created string
			err = rows.Scan(&username, &department, &created)
			//checkErr(err)

			fmt.Print(username, "  ")
			fmt.Print(department, "  ")
			fmt.Print(created, "  \n")
		}
	}

}

// 删除数据

func (r *SqliteHandle) Delete() {
	sql := r.DB
	stmt, err := sql.Prepare("DELETE  FROM userinfo WHERE username = ?")
	if err != nil {
		log.Fatal(err)
	}
	result, err := stmt.Exec("astaxie")
	affectNum, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("delete affect rows is ", affectNum)

}

// 数据更新

func (r *SqliteHandle) Update() {
	sql := r.DB
	stmt, err := sql.Prepare("UPDATE userinfo SET created = ? WHERE username = ?")
	if err != nil {
		log.Fatal(err)
	}
	result, err := stmt.Exec("2016-09-7", "我的名字")
	affectNum, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("update affect rows is ", affectNum)
}

// 创建表
func (r *SqliteHandle) Create() {
	sql := r.DB
	sqlStmt := `
    create table userinfo (username text , departname text , created text);
    `
	_, err := sql.Exec(sqlStmt)
	checkErr(err)
	//fmt.Println(f)
}

// 插入数据
func (r *SqliteHandle) Insert() {
	sql := r.DB
	stmt, err := sql.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
	checkErr(err)
	res, err := stmt.Exec("我的名字", "IEC118", time.Now())
	checkErr(err)
	_, err = res.LastInsertId()
	checkErr(err)
	//fmt.Println("ID  ...  ",id)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Prepare Error  ", err)
	}
}
