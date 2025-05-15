package goofer

import "time"


type Post struct {
    ID        uint      `orm:"primaryKey;autoIncrement"`
    Title     string    `orm:"type:varchar(255);notnull"`
    Content   string    `orm:"type:text"`
    UserID    uint      `orm:"index;notnull"`
    CreatedAt time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
    User      *User     `orm:"relation:ManyToOne;foreignKey:UserID"`
}
func (Post) TableName() string { return "posts" }


