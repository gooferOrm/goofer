package goofer

import (
	"github.com/gooferOrm/goofer/engine"
)

func Init(){
	_,err := engine.Connect("sqlite3", "./db.db")
	if err != nil {
		panic(err)
	}
	// postRepo := engine.Repo[Post](c)
}