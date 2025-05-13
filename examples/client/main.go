package main

import (
    "database/sql"
    "log"
    "time"

    _ "github.com/mattn/go-sqlite3"
    "github.com/gooferOrm/goofer/dialect"
    "github.com/gooferOrm/goofer/engine"
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
        dialect.NewSQLiteDialect(),
        User{}, Post{},
    )
    if err != nil {
        log.Fatalf("engine init: %v", err)
    }

    // Now use your repos directly:
    u := &User{Name: "Bob", Email: "bob@example.com"}
    if err := gooferClient.Repo[User]().Save(u); err != nil {
        log.Fatalf("save user: %v", err)
    }
    log.Printf("user ID = %d", u.ID)

    p := &Post{Title: "Hi", Content: "First post", UserID: u.ID}
    if err := gooferClient.Repo[Post]().Save(p); err != nil {
        log.Fatalf("save post: %v", err)
    }
    log.Printf("post ID = %d", p.ID)
}
