package model

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

func FetchAccount(tx *gorm.DB, id uuid.UUID) (account *Account) {
	tx.First(account, id)
	return
}

func FindAccountByEmail(tx *gorm.DB, email string) (*Account) {
	var account = new(Account)
	tx.FirstOrCreate(account, Account{
		Email: email,
	})
	
	return account
}
