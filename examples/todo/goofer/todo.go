package goofer


type Todo struct{
	ID uint 		`orm:"primaryKey;autoIncrement"`
	Title string	`orm:"type:varchar;notnull"`
	desc string		`orm:type:varchar;default:null`
}

