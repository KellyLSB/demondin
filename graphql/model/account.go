package model

import (
	"github.com/jinzhu/gorm"
)

func FindAccountByEmail(tx *gorm.DB, email string) (*Account) {
	var account = new(Account)
	tx.FirstOrCreate(account, Account{
		Email: email,
	})
	
	return account
}
