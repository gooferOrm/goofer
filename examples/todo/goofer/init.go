package goofer

import (
	"database/sql"
	"log"

)

type Goofer struct{
	db *sql.DB
	// dialect any
}

func Init()Goofer{
	db,err := sql.Open("sqlite3","./goofer.sqlite")
	//TODO:Use tripwire to handler errors later
	if err != nil{
		log.Fatalf("Failed to initialize the database: %v",err)
		return Goofer{db:nil} //Will fix the logic here later
	}
	defer db.Close()
	client := Goofer{db}
	return client

}