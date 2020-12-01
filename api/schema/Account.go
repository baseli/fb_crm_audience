package schema

import "time"

type Account struct {
	AccountId		string	`gorm:"primaryKey,size:31,not null"`
	Email			string	`gorm:"size:63"`
	Token			string	`gorm:"size:255"`
	Name			string	`gorm:"size:63"`
	AuthTime		time.Time
	NextAuthTime	time.Time
}
