// package data

// import (
// 	"database/sql"
// 	"log"

// 	_ "github.com/lib/pq"
// )

// //CREATE is a sql query
// var CREATE = "CREATE TABLE sites (id integer, domain varchar(40), userid integer, datetime timestamp without time zone)"

// // CreateSites db query to create site table
// func CreateSites() string {
// 	return CREATE
// }

// //Init initializedata
// func Init() {

// 	db, err := sql.Open("postgres", "user=Garrett dbname=cms_scratch sslmode=disable")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer db.Close()
// }
