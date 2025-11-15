package models

import "time"


type Comment struct {
	ID				uint   `gorm:"primaryKey" json:"id"`
	Name			string `json:"name" gorm:"size:100;not null"`
	Message			string `json:"message" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
}