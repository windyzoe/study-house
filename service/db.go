package service

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Start() {
	initDB()
	createDB()
}

func initDB() {
	log.Print("initDB--SUCCESS")
	var err error
	DB, err = sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func createDB() {
	log.Print("createDB--SUCCES")
	sql, err := ioutil.ReadFile("./create.sql")
	if err != nil {
		log.Fatal(err)
	}
	sqlStmt := string(sql)
	_, err = DB.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

// 数据库更新操作,有id走更新,没id创建
// rows:默认会有id,id=0就是创建
func UpdateWithCommit(tableName string, rows []map[string]interface{}) {
	transaction, err := DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	if len(rows) == 0 {
		return
	}
	// 字符串拼接,会把id排除
	oneRow := rows[0]
	keys := make([]string, 0, len(oneRow))
	spaces := make([]string, 0, len(oneRow))
	updateKeys := make([]string, 0, len(oneRow))
	for k := range oneRow {
		if k != "Id" {
			keys = append(keys, k)
			spaces = append(spaces, "?")
			updateKeys = append(updateKeys, k+"=?")
		}
	}
	strKeys := strings.Join(keys, ",")
	strSpaces := strings.Join(spaces, ",")
	strUpdateKeys := strings.Join(updateKeys, ",")
	insertSql := "insert into " + tableName + "(" + strKeys + ") values(" + strSpaces + ")"
	updateSql := "update " + tableName + " set " + strUpdateKeys + " where Id=?"
	// 准备
	prepareInsert, err := transaction.Prepare(insertSql)
	if err != nil {
		log.Fatal(err)
	}
	prepareUpdate, err := transaction.Prepare(updateSql)
	if err != nil {
		log.Fatal(err)
	}
	defer prepareInsert.Close()
	defer prepareUpdate.Close()
	for _, rowMap := range rows {
		values := make([]interface{}, 0, len(oneRow))
		id, idExsit := rowMap["Id"]
		idInt, idOk := id.(int64)
		for _, key := range keys {
			v, ok := rowMap[key]
			if ok {
				values = append(values, v)
			} else {
				values = append(values, "")
			}
		}
		// 基于id拆分
		if idExsit && idOk && idInt != 0 {
			values = append(values, id)
			_, err = prepareUpdate.Exec(values...)
			log.Println(values)
		} else {
			log.Println(values)
			_, err = prepareInsert.Exec(values...)
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	err = transaction.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// go-sqlite 官网示例
func example() {
	var err error
	DB, err = sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	sqlStmt := `
	create table foo (id integer not null primary key, name text);
	delete from foo;
	`
	_, err = DB.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	tx, err := DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()

	rows, err := DB.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err = DB.Prepare("select name from foo where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow("3").Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)

	_, err = DB.Exec("delete from foo")
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	if err != nil {
		log.Fatal(err)
	}

	rows, err = DB.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
