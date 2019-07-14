package model

import (
	"net/http"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/go-macaron/session"
	"github.com/KellyLSB/demondin/graphql/utils"
)

func UpdateSession(tx *gorm.DB, s session.Store, r *http.Request) (*Session) {
	var session Session

	if sUUID := utils.EnsureUUID(s.Get("sessionUUID")); sUUID != uuid.Nil {
		tx.Where("id = ?", sUUID).Preload("Account", "Invoice").First(&session)
	} else {
		tx.Create(&session)
		s.Set("sessionUUID", session.ID)
	}

	session.RemoteAddr = &r.RemoteAddr
	session.UserAgent  = utils.StringPtr(r.UserAgent())
	session.Referer    = utils.StringPtr(r.Referer())
	session.Method     = &r.Method
	session.URL        = utils.StringPtr(r.URL.String())

	session.InitializeInvoice(tx)

	tx.Save(&session)

	return &session
}

func (s *Session) AttachAccount(a *Account) (*Session) {
	s.AccountID = &a.ID
	s.Account = a
	return s
}

func (s *Session) InitializeInvoice(tx *gorm.DB) (*Session) {
	if s.Invoice == nil {
		s.AttachInvoice(FetchOrCreateInvoice(tx, s.InvoiceID))
	}
	
	return s
}

func (s *Session) AttachInvoice(i *Invoice) (*Session) {
	s.InvoiceID = &i.ID
	s.Invoice = i
	return s
}
