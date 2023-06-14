package models

import "time"

type User struct {
	Id          int64     `gorm:"primaryKey" json:"id"`
	NamaLengkap string    `gorm:"varchar(300)" json:"nama_lengkap"`
	Email       string    `gorm:"varchar(300)" json:"email"`
	Password    string    `gorm:"varchar(300)" json:"password"`
	Username    string    `gorm:"varchar(300)" json:"username"`
	Accounts    []Account `gorm:"foreignKey:UserId"`
}

type Account struct {
	AccountId int64    `gorm:"primaryKey" json:"account_id"`
	UserId    int64    `gorm:"integer" json:"user_id"`
	Number    int64    `gorm:"integer" json:"number"`
	Pin       string   `gorm:"varchar(300)" json:"pin"`
	Balance   int64    `gorm:"integer" json:"balance"`
	Status    Status   `gorm:"foreignKey:user_id"`
	Activity  Activity `gorm:"foreignKey:user_id"`
}

type Status struct {
	StatusId   int64 `gorm:"primaryKey" json:"status_id"`
	UserId     int64 `gorm:"integer" json:"user_id"`
	AccNumber  int64 `gorm:"integer" json:"acc_number"`
	PinAttempt int64 `gorm:"integer" json:"pin_attempt"`
	IsBlocked  int64 `gorm:"integer" json:"is_blocked"`
}

type Activity struct {
	ActivityId  int64     `gorm:"primaryKey" json:"activity_id"`
	UserId      int64     `gorm:"integer" json:"user_id"`
	AccNumber   int64     `gorm:"integer" json:"acc_number"`
	Description string    `gorm:"varchar(300)" json:"description"`
	Recipient   int64     `gorm:"integer" json:"recipient"`
	DateTime    time.Time `gorm:"autoCreateTime:false"`
}
