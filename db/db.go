package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/windyzoe/study-house/util"
)

var DB *sql.DB

func Start() {
	initDB()
	createDB()
}

func initDB() {
	log.Print("initDB--SUCCESS")
	var err error
	DB, err = sql.Open("sqlite3", util.Configs.Db.Path)
	if err != nil {
		log.Println(err)
	}
	err = DB.Ping()
	if err != nil {
		log.Println(err)
	}
}

func createDB() {
	log.Print("createDB--SUCCES")
	sql, err := ioutil.ReadFile("./create.sql")
	if err != nil {
		log.Println(err)
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
		log.Println(err)
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
		log.Println(err)
	}
	prepareUpdate, err := transaction.Prepare(updateSql)
	if err != nil {
		log.Println(err)
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
			log.Println(err)
		}
	}
	err = transaction.Commit()
	if err != nil {
		log.Println(err)
	}
}

// 通用查询的封装
// @selectStr select的key值,"a.name,b.name"
// @mappingKeys key值映射到map里的实际的key aName bName
// @fromStr from后面那一大堆
func Query(mapper map[string]string, fromStr string) (list []map[string]interface{}) {
	selectStr, mappingKeys := getSelectMapperKeys(mapper)
	querySql := "select " + selectStr + " from " + fromStr
	log.Println(querySql)
	rows, err := DB.Query(querySql)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	keyLength := len(mappingKeys)
	cache := make([]interface{}, keyLength)
	for i, _ := range cache {
		var a interface{}
		cache[i] = &a
	}
	for rows.Next() {
		err = rows.Scan(cache...)
		if err != nil {
			log.Println(err)
		}
		item := make(map[string]interface{})
		for i, v := range cache {
			item[mappingKeys[i]] = *v.(*interface{})
		}
		list = append(list, item)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}
	return
}

func QueryWithCount(mapper map[string]string, fromStr string) (list []map[string]interface{}, count int64) {
	list = Query(mapper, fromStr)

	selectStr, _ := getSelectMapperKeys(mapper)
	querySql := "select " + selectStr + " from " + fromStr
	countMapper := map[string]string{
		"Count": "count(*)",
	}
	countMap := Query(countMapper, "("+querySql+")")
	countIf := countMap[0]["Count"]
	count = countIf.(int64)
	return
}

type PageData struct {
	PageNumber int64                    `json:"pageNumber"`
	PageSize   int64                    `json:"pageSize"`
	Count      int64                    `json:"count"`
	List       []map[string]interface{} `json:"list"`
}

func QueryPage(mapper map[string]string, fromStr string, pageNumber int64, pageSize int64) (pageData PageData) {
	offset := (pageNumber - 1) * pageSize
	limit := pageSize
	pageData.List = Query(mapper, fromStr+" limit "+strconv.FormatInt(limit, 10)+" offset "+strconv.FormatInt(offset, 10))
	// count
	selectStr, _ := getSelectMapperKeys(mapper)
	countSql := "select " + selectStr + " from " + fromStr
	countMapper := map[string]string{
		"Count": "count(*)",
	}
	countMap := Query(countMapper, "("+countSql+")")
	countInterface := countMap[0]["Count"]
	pageData.Count = countInterface.(int64)
	pageData.PageNumber = pageNumber
	pageData.PageSize = pageSize
	return
}

// tableName 表名 keys 要查的字段 id 要查的id,小于1时为查全部
func QueryAll(tableName string, keys []string, id int64) []map[string]interface{} {
	var querySql string
	if id > 0 {
		querySql = "select " + strings.Join(keys, ",") + " from " + tableName + " where id=" + strconv.FormatInt(id, 10)
	} else {
		querySql = "select " + strings.Join(keys, ",") + " from " + tableName
	}
	rows, err := DB.Query(querySql)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	keyLength := len(keys)
	cache := make([]interface{}, keyLength)
	for i, _ := range cache {
		var a interface{}
		cache[i] = &a
	}
	var list []map[string]interface{}
	for rows.Next() {
		err = rows.Scan(cache...)
		if err != nil {
			log.Println(err)
		}
		item := make(map[string]interface{})
		for i, v := range cache {
			item[keys[i]] = *v.(*interface{})
		}
		list = append(list, item)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}
	return list
}

func getSelectMapperKeys(mapper map[string]string) (selectStr string, mappingKeys []string) {
	selectStrs := []string{}
	for k, v := range mapper {
		selectStrs = append(selectStrs, v)
		mappingKeys = append(mappingKeys, k)
	}
	selectStr = strings.Join(selectStrs, ",")
	return
}

// go-sqlite 官网示例
func example() {
	var err error
	DB, err = sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Println(err)
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
		log.Println(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
		if err != nil {
			log.Println(err)
		}
	}
	tx.Commit()

	rows, err := DB.Query("select id, name from foo")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}

	stmt, err = DB.Prepare("select name from foo where id = ?")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow("3").Scan(&name)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(name)

	_, err = DB.Exec("delete from foo")
	if err != nil {
		log.Println(err)
	}

	_, err = DB.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	if err != nil {
		log.Println(err)
	}

	rows, err = DB.Query("select id, name from foo")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}
}
