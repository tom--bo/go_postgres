package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"unicode/utf8"
)

func main() {
	//データの検索
	ti := 0
	for ti < 20000 {
		db, err := sql.Open("postgres", "user=postgres password=p dbname=mecab sslmode=disable")
		checkErr(err)
		var surface_arr [5]string
		var original_arr [5]string
		var reading_arr [5]string
		var length_arr [5]int

		tx, err := db.Begin()

		if err != nil {
			fmt.Println("トランザクションの取得に失敗しました。: %v", err)
		}

		rows, err := tx.Query("SELECT surface, original, reading FROM words_queue LIMIT 5 FOR UPDATE")
		checkErr(err)

		cnt := 0
		for rows.Next() {
			var surface string
			var original string
			var reading string
			err = rows.Scan(&surface, &original, &reading)
			checkErr(err)

			surface_arr[cnt] = surface
			original_arr[cnt] = original
			reading_arr[cnt] = reading
			length_arr[cnt] = utf8.RuneCountInString(reading)
			cnt++
		}

		for i := 0; i < 5; i++ {
			query := "DELETE FROM words_queue WHERE original = $1"
			_, err := tx.Exec(query, original_arr[i])
			checkErr(err)
		}

		// 本来ならerrの内容を確認してcommitまたはrollbackを決める必要がある
		err = tx.Commit()

		if err != nil {
			fmt.Println("トランザクションのコミットに失敗しました。: %v", err)
			return
		}

		for j := 0; j < 5; j++ {
			if length_arr[j] < 3 {
				stmt, err := db.Prepare("INSERT INTO metadata(original, minimum_length) VALUES($1,$2)")
				checkErr(err)

				_, err = stmt.Exec(original_arr[j], length_arr[j])
				checkErr(err)
			} else {
				tmp_len := length_arr[j]

				for tmp_len >= 2 {
					var wcnt int
					row, err := db.Query("SELECT count(*) FROM words WHERE reading LIKE $1", reading_arr[j][0:tmp_len*3]+"%")
					for row.Next() {
						err = row.Scan(&wcnt)
						checkErr(err)

						if wcnt >= 2 && tmp_len == length_arr[j] {
							stmt, err := db.Prepare("INSERT INTO metadata(original, minimum_length) VALUES($1,$2)")
							checkErr(err)

							_, err = stmt.Exec(original_arr[j], length_arr[j])
							checkErr(err)
							break
						} else if wcnt >= 2 {
							stmt, err := db.Prepare("INSERT INTO metadata(original, minimum_length) VALUES($1,$2)")
							checkErr(err)

							_, err = stmt.Exec(original_arr[j], tmp_len+1)
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
		ti += 1
		if ti%50 == 0 {
			fmt.Println(ti)
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
