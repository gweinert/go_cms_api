package models

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

//InitDB init db
func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}

// func BulkInsert(unsavedRows []*ExampleRowStruct) error {
//     valueStrings := make([]string, 0, len(unsavedRows))
//     valueArgs := make([]interface{}, 0, len(unsavedRows) * 3)
//     for _, post := range unsavedRows {
//         valueStrings = append(valueStrings, "(?, ?, ?)")
//         valueArgs = append(valueArgs, post.Column1)
//         valueArgs = append(valueArgs, post.Column2)
//         valueArgs = append(valueArgs, post.Column3)
//     }
//     stmt := fmt.Sprintf("INSERT INTO my_sample_table (column1, column2, column3) VALUES %s", strings.Join(valueStrings, ","))
//     _, err := db.Exec(stmt, valueArgs...)
//     return err
// }
