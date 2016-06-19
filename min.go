package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"unicode/utf8"
)

func main() {
	//データの検索
	cnt := 2
	for cnt < 51 {
		db, err := sql.Open("postgres", "user=postgres password=pass dbname=mecab sslmode=disable")
		checkErr(err)

		rows, err := db.Query("SELECT * FROM words WHERE original NOT IN(select original from metadata) LIMIT 100")
		checkErr(err)

		for rows.Next() {
			var surface string
			var original string
			var reading string
			var first_c sql.NullString
			var last_c sql.NullString
			var length int
			var wcnt int
			err = rows.Scan(&surface, &original, &reading, &first_c, &last_c)
			checkErr(err)
			/*
				fmt.Println(original)
				fmt.Println(surface)
				fmt.Println(reading)
					fmt.Println(reading[0:3])
					fmt.Println(reading[len(reading)-3 : len(reading)])
					fmt.Println(first_c)
					fmt.Println(last_c)
			*/

			length = utf8.RuneCountInString(reading)

			if length < 3 {
				// insert
				//データの挿入
				stmt, err := db.Prepare("INSERT INTO metadata(original, same_reading_count, minimum_length) VALUES($1,$2,$3)")
				checkErr(err)

				_, err = stmt.Exec(original, length, length)
				checkErr(err)
			} else {
				tmp_len := length

				for tmp_len >= 2 {
					row, err := db.Query("SELECT count(*) FROM words WHERE reading LIKE $1", reading[0:tmp_len*3]+"%")
					for row.Next() {
						err = row.Scan(&wcnt)
						checkErr(err)

						if wcnt >= 2 && tmp_len == length {
							stmt, err := db.Prepare("INSERT INTO metadata(original, same_reading_count, minimum_length) VALUES($1,$2,$3)")
							checkErr(err)

							_, err = stmt.Exec(original, wcnt, length)
							checkErr(err)
							break
						} else if wcnt >= 2 {
							stmt, err := db.Prepare("INSERT INTO metadata(original, same_reading_count, minimum_length) VALUES($1,$2,$3)")
							checkErr(err)

							_, err = stmt.Exec(original, 1, tmp_len+1)
							checkErr(err)
							break
						}
					}
					if wcnt >= 2 {
						break
					}
					tmp_len -= 1
				}
			}
		}
		cnt += 1
		if cnt%50 == 0 {
			fmt.Println(cnt)
			break
		}
		db.Close()
	}
	fmt.Println("finish")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
