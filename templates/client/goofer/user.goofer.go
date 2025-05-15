package goofer

import (
	"time"
)

type User struct {
    ID        uint      `orm:"primaryKey;autoIncrement" validate:"required"`
    Name      string    `orm:"type:varchar(255);notnull" validate:"required"`
    Email     string    `orm:"unique;type:varchar(255);notnull" validate:"required,email"`
    CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
    Posts     []Post    `orm:"relation:OneToMany;foreignKey:UserID"`
}

func (User) TableName() string { return "users" }

