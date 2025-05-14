package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gooferOrm/goofer/dialect"
	"github.com/gooferOrm/goofer/engine"
	_ "github.com/mattn/go-sqlite3"
)

// 1) Define entities as before:
type User struct {
    ID        uint      `orm:"primaryKey;autoIncrement" validate:"required"`
    Name      string    `orm:"type:varchar(255);notnull" validate:"required"`
    Email     string    `orm:"unique;type:varchar(255);notnull" validate:"required,email"`
    CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
    Posts     []Post    `orm:"relation:OneToMany;foreignKey:UserID"`
}
func (User) TableName() string { return "users" }

type Post struct {
    ID        uint      `orm:"primaryKey;autoIncrement"`
    Title     string    `orm:"type:varchar(255);notnull"`
    Content   string    `orm:"type:text"`
    UserID    uint      `orm:"index;notnull"`
    CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
    User      *User     `orm:"relation:ManyToOne;foreignKey:UserID"`
}
func (Post) TableName() string { return "posts" }

func main() {
    db, err := sql.Open("sqlite3", "./db.db")
    if err != nil {
        log.Fatalf("open db: %v", err)
    }
    defer db.Close()

    // engine setup:
    gooferClient, err := engine.NewClient(
        db,
        dialect.NewSQLiteDialect(), //Todo: add the new dialect function on other dialects as well
        User{}, Post{},
    )
    if err != nil {
        log.Fatalf("engine init: %v", err)
    }

    // Now use your repos directly:
    u := &User{Name: "Bob", Email: "bob@example1.com"}
    if err := engine.Repo[User](gooferClient).Save(u); err != nil {
        log.Fatalf("save user: %v", err)
    }
    log.Printf("user ID = %d", u.ID)

    p := &Post{Title: "Hi", Content: "First post", UserID: u.ID}
    postRepo := engine.Repo[Post](gooferClient)
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
