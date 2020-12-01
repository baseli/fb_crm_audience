package schema

import "time"

type Task struct {
	Id			uint	`gorm:"primaryKey,autoincrement"`
	AdAccountId	string	`gorm:"index,size:31,not null"`
	File		string	`gorm:"size:255,not null"`
	Msg			string	`gorm:"size:511,not null,default:''"`
	Status		uint	`gorm:"not null,default:0"`
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}
