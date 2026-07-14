package recorder

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Row struct {
	User    string
	Ip      string
	Time    string
	Message string
}

func Recorder(messages chan Row, database *sql.DB) {

	defer database.Close()

	for {
		currentRow := <-messages

		res, err := database.Exec(
			"insert into first_table(`user`, ip, `time`, message) values (?, ?, ?, ?)",
			currentRow.User,
			currentRow.Ip,
			currentRow.Time,
			currentRow.Message,
		)

		if err != nil {
			log.Fatal(err)
		}
		id, _ := res.LastInsertId()
		fmt.Println("[+] Inserted ID:", id)
	}
}
