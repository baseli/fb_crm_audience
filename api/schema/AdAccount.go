package schema

type AdAccount struct {
	Id				uint	`gorm:"primaryKey,autoincrement"`
	AdAccountId		string	`gorm:"index,size:31"`
	AccountId		string	`gorm:"index,size:31"`
	Name			string	`gorm:"size:63"`
}
