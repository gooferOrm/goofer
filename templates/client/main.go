package main

import (
	"fmt"
	"log"

	"github.com/gooferOrm/goofer/engine"
	"github.com/gooferOrm/goofer/examples/custom_queries/goofer"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
    // db, err := sql.Open("sqlite3", "./db.db")
    // if err != nil {
    //     log.Fatalf("open db: %v", err)
    // }
    // defer db.Close()
    GooferClient,err := engine.Connect("sqlite3", "./db.db")
	if err != nil {
		panic(err)
	}
    

    // engine setup:
    // GooferClient, err := engine.NewClient(
    //     db,
    //     dialect.NewSQLiteDialect(),
    //     goofer.User{}, goofer.Post{},
    // )
    if err != nil {
        log.Fatalf("engine init: %v", err)
    }

    // Now use your repos directly:
    u := &goofer.User{Name: "Bob", Email: "bob@example1.com"}
    if err := engine.Repo[goofer.User](GooferClient).Save(u); err != nil {
        log.Fatalf("save user: %v", err)
    }
    log.Printf("user ID = %d", u.ID)

    p := &goofer.Post{Title: "Hi", Content: "First post", UserID: u.ID}
    postRepo := engine.Repo[goofer.Post](GooferClient)
    if err := postRepo.Save(p); err != nil {
        log.Fatalf("save post: %v", err)
    }
    allPost,err := postRepo.Find().All()
    if err != nil {
        log.Fatalf("Error: %v",err)
        return
    }
    fmt.Println(allPost)
    for _, post := range allPost{
        fmt.Println(post)
    }
    log.Printf("post ID = %d", p.ID)
}
