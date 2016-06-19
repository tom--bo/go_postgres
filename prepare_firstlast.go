package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "user=postgres password=pass dbname=mecab sslmode=disable")
	checkErr(err)
	//データの検索
	//rows, err := db.Query("SELECT * FROM words LIMIT 3")
	rows, err := db.Query("SELECT * FROM words WHERE first_char IS NULL")
	checkErr(err)

	for rows.Next() {
		var surface string
		var original string
		var reading string
		var first_c sql.NullString
		var last_c sql.NullString
		err = rows.Scan(&surface, &original, &reading, &first_c, &last_c)
		checkErr(err)
		/*
			fmt.Println(surface)
			fmt.Println(original)
			fmt.Println(reading)
			fmt.Println(reading[0:3])
			fmt.Println(reading[len(reading)-3 : len(reading)])
			fmt.Println(first_c)
			fmt.Println(last_c)
		*/

		first_c.String = reading[0:3]
		last_c.String = reading[len(reading)-3 : len(reading)]

		//データの更新
		stmt, err := db.Prepare("update words set first_char=$1, last_char=$2 where original=$3")
		checkErr(err)

		first_c.Valid = true
		last_c.Valid = true
		_, err = stmt.Exec(first_c.String, last_c.String, original)
		checkErr(err)

		/*
			affect, err := res.RowsAffected()
			checkErr(err)
			fmt.Println(affect)
		*/
	}

	//データの挿入
	/*
		stmt, err := db.Prepare("INSERT INTO userinfo(username,departname,created) VALUES($1,$2,$3) RETURNING uid")
		checkErr(err)

		res, err := stmt.Exec("astaxie", "研究開発部門", "2012-12-09")
		checkErr(err)

		//pgはこの関数をサポートしていません。MySQLのインクリメンタルなIDのようなものが無いためです。
		id, err := res.LastInsertId()
		checkErr(err)

		fmt.Println(id)
	*/

	//データの削除
	/*
		stmt, err = db.Prepare("delete from userinfo where uid=$1")
		checkErr(err)

		res, err = stmt.Exec(1)
		checkErr(err)

		affect, err = res.RowsAffected()
		checkErr(err)

		fmt.Println(affect)
	*/
	db.Close()
	fmt.Println("finish")

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
