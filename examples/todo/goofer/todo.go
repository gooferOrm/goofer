package goofer

// import "github.com/gooferOrm/goofer/schema"


type Todo struct{
	ID uint 		`orm:"primaryKey;autoIncrement"`
	Title string	`orm:"type:varchar;notnull"`
	Desc string		`orm:type:varchar;default:null`
}

func (Todo) TableName() string{
	return "todos"
}

// var err error = schema.Registry.RegisterEntity(Todo{})
