package main

import (
	"fmt"
	"log"

	"github.com/gooferOrm/goofer/dialect"
	"github.com/gooferOrm/goofer/examples/todo/goofer"
	"github.com/gooferOrm/goofer/schema"
)


func main(){
	sqliteDialect := dialect.NewSQLiteDialect()
	fmt.Println(sqliteDialect.Name())

	if err:= schema.Registry.RegisterEntity(goofer.Todo{}); err != nil{
		log.Fatalf("Error while registering the todo entity: %v",err)
		return
	}

}