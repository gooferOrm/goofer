package goofer

import (
	"time"
)

type User struct {
	ID        uint      `orm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Name      string    `orm:"size:255;not null" json:"name" validate:"required"`
	Email     string    `orm:"size:255;not null;uniqueIndex" json:"email" validate:"required,email"`
	CreatedAt time.Time `orm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `orm:"autoUpdateTime" json:"updated_at"`
	Posts     []Post    `orm:"foreignKey:UserID" json:"posts,omitempty"`
}

// TableName specifies the database table name for the User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate is a hook that runs before creating a user
func (u *User) BeforeCreate() error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate is a hook that runs before updating a user
func (u *User) BeforeUpdate() error {
	u.UpdatedAt = time.Now()
	return nil
}
