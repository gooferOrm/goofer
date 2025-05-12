package goofer

import (
	"database/sql"
	"log"
)

func Init()*sql.DB{
	db,err := sql.Open("sqlite3","./goofer.sqlite")
	//TODO:Use tripwire to handler errors later

	if err != nil{
		log.Fatalf("Failed to initialize the database: %v",err)
		return nil
	}
	return db

}