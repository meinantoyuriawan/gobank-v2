package models

type User struct {
	Id          int64     `gorm:"primaryKey" json:"id"`
	NamaLengkap string    `gorm:"varchar(300)" json:"nama_lengkap"`
	Email       string    `gorm:"varchar(300)" json:"email"`
	Password    string    `gorm:"varchar(300)" json:"password"`
	Username    string    `gorm:"varchar(300)" json:"username"`
	Accounts    []Account `gorm:"foreignKey:UserId"`
}

type Account struct {
	AccountId int64  `gorm:"primaryKey" json:"account_id"`
	UserId    int64  `gorm:"integer" json:"user_id"`
	Number    int64  `gorm:"integer" json:"number"`
	Pin       string `gorm:"varchar(300)" json:"pin"`
	Balance   int64  `gorm:"integer" json:"balance"`
}

// type Status struct {
// 	AccountId
// 	PinAttempt
// 	IsBlocked
// }
