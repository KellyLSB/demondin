package models

import "github.com/google/uuid"

type M2M struct {
	XID, YID uuid.UUID `gorm:"type:uuid; index;"`
	XNm, YNm string
}

func ManyOf(m interface{}) {

}
