package goofer

import "time"

type User struct {
	ID        uint      `orm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Name      string    `orm:"size:255;not null" json:"name"`
	Email     string    `orm:"size:255;not null;uniqueIndex" json:"email"`
	CreatedAt time.Time `orm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `orm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type Post struct {
	ID        uint      `orm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Title     string    `orm:"size:255;not null" json:"title"`
	Content   string    `orm:"type:text;not null" json:"content"`
	UserID    uint      `orm:"index;not null" json:"user_id"`
	CreatedAt time.Time `orm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `orm:"autoUpdateTime" json:"updated_at"`
}

func (Post) TableName() string {
	return "posts"
}
