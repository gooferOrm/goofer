package goofer

import "time"

type Post struct {
	ID        uint      `orm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Title     string    `orm:"size:255;not null" json:"title" validate:"required"`
	Content   string    `orm:"type:text;not null" json:"content" validate:"required"`
	UserID    uint      `orm:"index;not null" json:"user_id"`
	CreatedAt time.Time `orm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `orm:"autoUpdateTime" json:"updated_at"`
	User      *User     `orm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the database table name for the Post model
func (Post) TableName() string {
	return "posts"
}

// BeforeCreate is a hook that runs before creating a post
func (p *Post) BeforeCreate() error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate is a hook that runs before updating a post
func (p *Post) BeforeUpdate() error {
	p.UpdatedAt = time.Now()
	return nil
}
